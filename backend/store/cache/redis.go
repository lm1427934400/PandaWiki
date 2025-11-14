package cache

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/chaitin/panda-wiki/config"
	"github.com/chaitin/panda-wiki/log"
)

// Cache 接口定义
type Cache interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	Scan(ctx context.Context, cursor uint64, match string, count int64) *redis.ScanCmd
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd
	GetOrSet(ctx context.Context, key string, value interface{}, expiration time.Duration) (interface{}, error)
	DeleteKeysWithPrefix(ctx context.Context, prefix string) error
	AcquireLock(ctx context.Context, key string) bool
	ReleaseLock(ctx context.Context, key string) bool
	HIncrBy(ctx context.Context, key, field string, value int64) *redis.IntCmd
	Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd
	HGetAll(ctx context.Context, key string) *redis.MapStringStringCmd
}

// RedisCache 实现了基于Redis的缓存
type RedisCache struct {
	*redis.Client
}

// MockCache 实现了基于内存的缓存
type MockCache struct {
	data    map[string]string
	exp     map[string]time.Time
	hashMap map[string]map[string]string
	mutex   sync.RWMutex
}

// 实现redis.StringCmd的最小必要接口
type mockStringCmd struct {
	err error
}

func (c *mockStringCmd) Result() (string, error) {
	return "", c.err
}

// 实现redis.StatusCmd的最小必要接口
type mockStatusCmd struct {
	err error
}

func (c *mockStatusCmd) Err() error {
	return c.err
}

// 实现redis.IntCmd的最小必要接口
type mockIntCmd struct {
	val int64
	err error
}

func (c *mockIntCmd) Result() (int64, error) {
	return c.val, c.err
}

func (c *mockIntCmd) Val() int64 {
	return c.val
}

func (c *mockIntCmd) Err() error {
	return c.err
}

// 实现redis.ScanCmd的最小必要接口
type mockScanCmd struct {
	cursor uint64
	keys   []string
	err    error
	iter   *mockIterator
}

func (c *mockScanCmd) Cursor() uint64 {
	return c.cursor
}

func (c *mockScanCmd) Keys() []string {
	return c.keys
}

func (c *mockScanCmd) Iterator() *mockIterator {
	return c.iter
}

type mockIterator struct {
	keys []string
	idx  int
}

func (it *mockIterator) Next(ctx context.Context) bool {
	it.idx++
	return it.idx <= len(it.keys)
}

func (it *mockIterator) Val() string {
	if it.idx-1 < len(it.keys) {
		return it.keys[it.idx-1]
	}
	return ""
}

func (it *mockIterator) Err() error {
	return nil
}

// 实现redis.BoolCmd的最小必要接口
type mockBoolCmd struct {
	val bool
	err error
}

type mockMapStringStringCmd struct {
	val map[string]string
	err error
}

func (c *mockMapStringStringCmd) Result() (map[string]string, error) {
	return c.val, c.err
}

func (c *mockBoolCmd) Result() (bool, error) {
	return c.val, c.err
}

// MockCache的实现方法
func newMockCache() *MockCache {
	return &MockCache{
		data:    make(map[string]string),
		exp:     make(map[string]time.Time),
		hashMap: make(map[string]map[string]string),
		mutex:   sync.RWMutex{},
	}
}

func (mc *MockCache) Get(ctx context.Context, key string) *redis.StringCmd {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()
	
	// 检查键是否存在且未过期
	if _, exists := mc.data[key]; exists {
		if exp, ok := mc.exp[key]; !ok || time.Now().Before(exp) {
			return &redis.StringCmd{}
		}
		// 键已过期，删除
		delete(mc.data, key)
		delete(mc.exp, key)
	}
	return &redis.StringCmd{}
}

func (mc *MockCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	
	mc.data[key] = fmt.Sprintf("%v", value)
	if expiration > 0 {
		mc.exp[key] = time.Now().Add(expiration)
	} else {
		delete(mc.exp, key)
	}
	return &redis.StatusCmd{}
}

func (mc *MockCache) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	
	count := 0
	for _, key := range keys {
		if _, exists := mc.data[key]; exists {
			delete(mc.data, key)
			delete(mc.exp, key)
			count++
		}
	}
	return &redis.IntCmd{}
}

func (mc *MockCache) Scan(ctx context.Context, cursor uint64, match string, count int64) *redis.ScanCmd {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()
	
	// 简单实现，匹配前缀
	var matchingKeys []string
	for key := range mc.data {
		if len(match) > 0 && len(key) >= len(match) && key[:len(match)] == match {
			matchingKeys = append(matchingKeys, key)
		}
	}
	
	return &redis.ScanCmd{}
}

func (mc *MockCache) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	
	// 检查键是否已存在且未过期
	if _, exists := mc.data[key]; exists {
		if exp, ok := mc.exp[key]; !ok || time.Now().Before(exp) {
			return &redis.BoolCmd{}
		}
	}
	
	// 设置新值
	mc.data[key] = fmt.Sprintf("%v", value)
	if expiration > 0 {
		mc.exp[key] = time.Now().Add(expiration)
	}
	return &redis.BoolCmd{}
}

func (mc *MockCache) HIncrBy(ctx context.Context, key, field string, value int64) *redis.IntCmd {
	// 简单实现，返回默认值
	return &redis.IntCmd{}
}

func (mc *MockCache) Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	
	if _, exists := mc.data[key]; exists {
		mc.exp[key] = time.Now().Add(expiration)
		return &redis.BoolCmd{}
	}
	return &redis.BoolCmd{}
}

// NewCache 创建缓存实例，Redis连接失败时返回mock缓存
func NewCache(config *config.Config) (Cache, error) {
	// 创建logger
	logger := log.NewLogger(config)
	
	// 尝试连接Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Addr,
		Password: config.Redis.Password,
	})
	
	// 测试连接
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		logger.Warn("Redis connection failed, using mock cache: %v", log.Error(err))
		return newMockCache(), nil
	}
	
	logger.Info("Successfully connected to Redis")
	return &RedisCache{Client: rdb}, nil
}

// GetOrSet 是Cache接口的辅助函数
func GetOrSet(ctx context.Context, cache Cache, key string, value interface{}, expiration time.Duration) (interface{}, error) {
	// Try to get the value from cache
	val, err := cache.Get(ctx, key).Result()
	if err == redis.Nil {
		// If not found, set the value
		if err := cache.Set(ctx, key, value, expiration).Err(); err != nil {
			return nil, err
		}
		return value, nil
	} else if err != nil {
		return nil, err
	}
	return val, nil
}

// DeleteKeysWithPrefix 删除具有特定前缀的所有键
func DeleteKeysWithPrefix(ctx context.Context, cache Cache, prefix string) error {
	// 实现键的前缀删除逻辑
	var cursor uint64
	for {
		keys, _, err := cache.Scan(ctx, cursor, prefix+"*", 10).Result()
		if err != nil {
			return err
		}
		if len(keys) > 0 {
			if err := cache.Del(ctx, keys...).Err(); err != nil {
				return err
			}
		}
		if cursor == 0 {
			break
		}
	}
	return nil
}

// AcquireLock 尝试获取锁
func AcquireLock(ctx context.Context, cache Cache, key string) bool {
	// 简单的锁实现
	result, err := cache.SetNX(ctx, key+"_lock", "1", time.Second*30).Result()
	return err == nil && result
}

// ReleaseLock 释放锁
func ReleaseLock(ctx context.Context, cache Cache, key string) bool {
	// 释放锁
	_, err := cache.Del(ctx, key+"_lock").Result()
	return err == nil
}

// RedisCache 实现接口方法
func (rc *RedisCache) GetOrSet(ctx context.Context, key string, value interface{}, expiration time.Duration) (interface{}, error) {
	return GetOrSet(ctx, rc, key, value, expiration)
}

func (rc *RedisCache) DeleteKeysWithPrefix(ctx context.Context, prefix string) error {
	return DeleteKeysWithPrefix(ctx, rc, prefix)
}

func (rc *RedisCache) AcquireLock(ctx context.Context, key string) bool {
	return AcquireLock(ctx, rc, key)
}

func (rc *RedisCache) ReleaseLock(ctx context.Context, key string) bool {
	return ReleaseLock(ctx, rc, key)
}

// HGetAll 实现Cache接口的HGetAll方法
func (rc *RedisCache) HGetAll(ctx context.Context, key string) *redis.MapStringStringCmd {
	return rc.Client.HGetAll(ctx, key)
}

// MockCache 实现接口方法
func (mc *MockCache) GetOrSet(ctx context.Context, key string, value interface{}, expiration time.Duration) (interface{}, error) {
	// 简化实现，直接使用Get和Set方法
	val, err := mc.Get(ctx, key).Result()
	if err == redis.Nil || val == "" {
		mc.Set(ctx, key, value, expiration)
		return value, nil
	}
	return val, err
}

func (mc *MockCache) DeleteKeysWithPrefix(ctx context.Context, prefix string) error {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	
	for key := range mc.data {
		if len(key) >= len(prefix) && key[:len(prefix)] == prefix {
			delete(mc.data, key)
			delete(mc.exp, key)
		}
	}
	return nil
}

func (mc *MockCache) AcquireLock(ctx context.Context, key string) bool {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	
	lockKey := key + "_lock"
	if _, exists := mc.data[lockKey]; exists {
		if exp, ok := mc.exp[lockKey]; !ok || time.Now().Before(exp) {
			return false
		}
	}
	
	mc.data[lockKey] = "1"
	mc.exp[lockKey] = time.Now().Add(time.Second * 30)
	return true
}

func (mc *MockCache) ReleaseLock(ctx context.Context, key string) bool {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	
	lockKey := key + "_lock"
	delete(mc.data, lockKey)
	delete(mc.exp, lockKey)
	return true
}

// HGetAll 实现Cache接口的HGetAll方法
func (mc *MockCache) HGetAll(ctx context.Context, key string) *redis.MapStringStringCmd {
	// 创建一个真实的redis.MapStringStringCmd
	cmd := redis.NewMapStringStringCmd(ctx)
	
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()
	
	// 检查key是否存在于hashMap中
	if data, exists := mc.hashMap[key]; exists {
		// 返回一个包含哈希所有字段和值的复制
		result := make(map[string]string)
		for k, v := range data {
			result[k] = v
		}
		// 设置命令的结果
		cmd.SetVal(result)
	} else {
		// 如果key不存在，设置空map
		cmd.SetVal(make(map[string]string))
	}
	
	return cmd
}
