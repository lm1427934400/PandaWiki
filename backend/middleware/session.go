package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/boj/redistore"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"

	"github.com/chaitin/panda-wiki/config"
	"github.com/chaitin/panda-wiki/log"
	"github.com/chaitin/panda-wiki/store/cache"
)

const (
	SessionKey = "SessionKey"
)

type SessionMiddleware struct {
	logger      *log.Logger
	store       *redistore.RediStore
	memoryStore sessions.Store
}

func NewSessionMiddleware(logger *log.Logger, config *config.Config, cache cache.Cache) (*SessionMiddleware, error) {

	secretKey, err := cache.GetOrSet(context.Background(), SessionKey, uuid.New().String(), time.Duration(0))
	if err != nil {
		logger.Error("session store create secret key failed: %v", log.Error(err))
		return nil, err
	}

	// 尝试创建Redis存储
	store, err := redistore.NewRediStore(
		10,
		"tcp",
		config.Redis.Addr,
		"",
		config.Redis.Password,
		[]byte(secretKey.(string)),
	)

	// 如果Redis连接失败，使用内存存储作为回退
	if err != nil {
		logger.Warn("Redis session store initialization failed, using in-memory store: %v", log.Error(err))
		// 使用gorilla/sessions的内存存储
		memoryStore := sessions.NewCookieStore([]byte(secretKey.(string)))
		memoryStore.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   30 * 86400,
			SameSite: http.SameSiteLaxMode,
			HttpOnly: true,
		}
		return &SessionMiddleware{
			logger:      logger.WithModule("middleware.session"),
			store:       nil,
			memoryStore: memoryStore,
		},
		nil
	}

	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   30 * 86400,
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
	}

	return &SessionMiddleware{
		logger: logger.WithModule("middleware.session"),
		store:  store,
	}, nil
}

func (s *SessionMiddleware) Session() echo.MiddlewareFunc {
	// 根据可用的存储类型选择合适的存储
	var store sessions.Store
	if s.store != nil {
		store = s.store
	} else {
		store = s.memoryStore
	}
	return session.MiddlewareWithConfig(session.Config{
		Store: store,
	})
}
