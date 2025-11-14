package pg

import (
	"context"

	v1 "github.com/chaitin/panda-wiki/api/stat/v1"
	"github.com/chaitin/panda-wiki/domain"
	"github.com/chaitin/panda-wiki/store/cache"
	"github.com/chaitin/panda-wiki/store/pg"
	"github.com/chaitin/panda-wiki/utils"
)

type StatRepository struct {
	db    *pg.DB
	cache cache.Cache
}

func NewStatRepository(db *pg.DB, cahe cache.Cache) *StatRepository {
	return &StatRepository{
		db:    db,
		cache: cahe,
	}
}

func (r *StatRepository) CreateStatPage(ctx context.Context, stat *domain.StatPage) error {
	return r.db.WithContext(ctx).Model(&domain.StatPage{}).Create(stat).Error
}

func (r *StatRepository) GetHotPages(ctx context.Context, kbID string) ([]*domain.HotPage, error) {
	var hotPages []*domain.HotPage
	if err := r.db.WithContext(ctx).Model(&domain.StatPage{}).
		Where("kb_id = ?", kbID).
		Where("scene = ?", domain.StatPageSceneNodeDetail).
		Group("node_id").
		Select("node_id, COUNT(*) as count").
		Order("count DESC").
		Limit(10).
		Find(&hotPages).Error; err != nil {
		return nil, err
	}
	return hotPages, nil
}

func (r *StatRepository) GetHotPagesNoLimit(ctx context.Context, kbID string) ([]*domain.HotPage, error) {
	var hotPages []*domain.HotPage
	if err := r.db.WithContext(ctx).Model(&domain.StatPage{}).
		Where("kb_id = ?", kbID).
		Where("scene = ?", domain.StatPageSceneNodeDetail).
		Group("node_id").
		Select("node_id, COUNT(*) as count").
		Find(&hotPages).Error; err != nil {
		return nil, err
	}
	return hotPages, nil
}

func (r *StatRepository) GetHotScene(ctx context.Context, kbID string) (map[domain.StatPageScene]int64, error) {
	var scenes map[domain.StatPageScene]int64
	if err := r.db.WithContext(ctx).Model(&domain.StatPage{}).
		Where("kb_id = ?", kbID).
		Group("scene").
		Select("scene, COUNT(*) as count").
		Order("count DESC").
		Limit(10).
		Find(&scenes).Error; err != nil {
		return nil, err
	}
	return scenes, nil
}

func (r *StatRepository) GetHotRefererHosts(ctx context.Context, kbID string) ([]*domain.HotRefererHost, error) {
	var hotRefererHosts []*domain.HotRefererHost
	if err := r.db.WithContext(ctx).Model(&domain.StatPage{}).
		Where("kb_id = ?", kbID).
		Group("referer_host").
		Select("referer_host, COUNT(*) as count").
		Order("count DESC").
		Limit(10).
		Find(&hotRefererHosts).Error; err != nil {
		return nil, err
	}
	return hotRefererHosts, nil
}

func (r *StatRepository) GetHotBrowsers(ctx context.Context, kbID string) (*domain.HotBrowser, error) {
	var hotBrowsers *domain.HotBrowser
	var osCount []domain.BrowserCount
	var browserCount []domain.BrowserCount

	query := r.db.WithContext(ctx).Model(&domain.StatPage{}).
		Where("kb_id = ?", kbID).
		Group("browser_name").
		Select("browser_name as name, COUNT(*) as count")
	if err := query.Order("count DESC").Limit(10).Find(&browserCount).Error; err != nil {
		return nil, err
	}

	query = r.db.WithContext(ctx).Model(&domain.StatPage{}).
		Where("kb_id = ?", kbID).
		Group("browser_os").
		Select("browser_os as name, COUNT(*) as count")
	if err := query.Order("count DESC").Limit(10).Find(&osCount).Error; err != nil {
		return nil, err
	}

	hotBrowsers = &domain.HotBrowser{
		OS:      osCount,
		Browser: browserCount,
	}

	return hotBrowsers, nil
}

func (r *StatRepository) GetStatPageCount(ctx context.Context, kbID string) (*v1.StatCountResp, error) {
	var count v1.StatCountResp
	if err := r.db.WithContext(ctx).Model(&domain.StatPage{}).
		Where("kb_id = ?", kbID).
		Select("COUNT(DISTINCT ip) as ip_count, COUNT(DISTINCT session_id) as session_count, COUNT(*) as page_visit_count").
		Scan(&count).Error; err != nil {
		return nil, err
	}
	return &count, nil
}

func (r *StatRepository) GetInstantCount(ctx context.Context, kbID string) ([]*domain.InstantCountResp, error) {
	var instantCount []*domain.InstantCountResp
	if err := r.db.WithContext(ctx).Model(&domain.StatPage{}).
		Where("kb_id = ? AND created_at >= NOW() - INTERVAL '1h'", kbID).
		Select("date_trunc('minute', created_at) as time, COUNT(*) as count").
		Group("time").
		Order("time ASC").
		Find(&instantCount).Error; err != nil {
		return nil, err
	}
	return instantCount, nil
}

func (r *StatRepository) GetInstantPages(ctx context.Context, kbID string) ([]*domain.InstantPageResp, error) {
	var instantPages []*domain.InstantPageResp
	if err := r.db.WithContext(ctx).Model(&domain.StatPage{}).
		Where("kb_id = ?", kbID).
		Select("node_id, ip, scene, created_at,user_id").
		Order("created_at DESC").
		Limit(10).
		Find(&instantPages).Error; err != nil {
		return nil, err
	}
	return instantPages, nil
}

func (r *StatRepository) RemoveOldData(ctx context.Context) error {
	if err := r.db.WithContext(ctx).Model(&domain.StatPage{}).
		Where("created_at < ?", utils.GetTimeHourOffset(-24)).
		Delete(&domain.StatPage{}).Error; err != nil {
		return err
	}
	return nil
}
