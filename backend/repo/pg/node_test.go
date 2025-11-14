package pg

import (
	"context"
	"testing"
	"time"

	"gorm.io/gorm"

	"github.com/chaitin/panda-wiki/domain"
	"github.com/chaitin/panda-wiki/log"
	"github.com/chaitin/panda-wiki/store/pg"
)

// MockDB 模拟数据库连接
type MockDB struct {
	*gorm.DB
}

// TestNodeRepository_GetList 测试GetList方法
func TestNodeRepository_GetList(t *testing.T) {
	// 准备测试数据
	ctx := context.Background()
	logger := log.NewLogger()
	
	// 创建测试数据库连接（实际项目中可能需要使用测试数据库或mock）
	// 这里为了简化，我们假设db已经初始化
	db, err := pg.NewDB(pg.Config{
		DSN: "host=localhost user=postgres password=postgres dbname=panda_wiki_test port=5432 sslmode=disable",
	})
	if err != nil {
		t.Skip("Skipping test due to database connection failure: " + err.Error())
	}
	
	// 创建仓库实例
	repo := NewNodeRepository(db, logger)
	
	// 清理测试数据
	defer func() {
		db.Exec("DELETE FROM nodes WHERE kb_id = ?", "test-kb-id")
		db.Exec("DELETE FROM users WHERE id IN (?, ?)", "test-creator-id", "test-editor-id")
	}()
	
	// 插入测试用户
	db.Exec(
		"INSERT INTO users (id, account) VALUES (?, ?), (?, ?)",
		"test-creator-id", "creator",
		"test-editor-id", "editor",
	)
	
	// 插入测试节点
	now := time.Now()
	db.Exec(
		"INSERT INTO nodes (id, kb_id, type, status, name, content, parent_id, position, creator_id, editor_id, created_at, updated_at, meta) " +
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		"test-node-id", "test-kb-id", domain.NodeTypeDoc, domain.NodeStatusActive, "Test Node", "Test Content", "", 0.0,
		"test-creator-id", "test-editor-id", now, now, 
		`{\"summary\": \"Test Summary\", \"emoji\": \"📝\", \"content_type\": \"md\"}`,
	)
	
	// 测试用例1: 基本查询
	t.Run("BasicQuery", func(t *testing.T) {
		req := &domain.GetNodeListReq{
			KBID: "test-kb-id",
		}
		
		nodes, err := repo.GetList(ctx, req)
		if err != nil {
			t.Fatalf("GetList failed: %v", err)
		}
		
		if len(nodes) != 1 {
			t.Fatalf("Expected 1 node, got %d", len(nodes))
		}
		
		node := nodes[0]
		// 验证节点字段
		if node.ID != "test-node-id" {
			t.Errorf("Expected ID 'test-node-id', got '%s'", node.ID)
		}
		if node.Name != "Test Node" {
			t.Errorf("Expected Name 'Test Node', got '%s'", node.Name)
		}
		if node.Creator != "creator" {
			t.Errorf("Expected Creator 'creator', got '%s'", node.Creator)
		}
		if node.Editor != "editor" {
			t.Errorf("Expected Editor 'editor', got '%s'", node.Editor)
		}
		if node.Summary != "Test Summary" {
			t.Errorf("Expected Summary 'Test Summary', got '%s'", node.Summary)
		}
		if node.Emoji != "📝" {
			t.Errorf("Expected Emoji '📝', got '%s'", node.Emoji)
		}
		if node.ContentType != "md" {
			t.Errorf("Expected ContentType 'md', got '%s'", node.ContentType)
		}
		// 验证时间字段映射正确（修复的重点）
		if node.UpdatedAt.IsZero() {
			t.Error("UpdatedAt should not be zero")
		}
		if node.CreatedAt.IsZero() {
			t.Error("CreatedAt should not be zero")
		}
	})
	
	// 测试用例2: 带搜索条件的查询
	t.Run("SearchQuery", func(t *testing.T) {
		req := &domain.GetNodeListReq{
			KBID:   "test-kb-id",
			Search: "Test",
		}
		
		nodes, err := repo.GetList(ctx, req)
		if err != nil {
			t.Fatalf("GetList with search failed: %v", err)
		}
		
		if len(nodes) == 0 {
			t.Error("Expected to find nodes with search term 'Test'")
		}
	})
	
	// 测试用例3: 不存在的知识库
	t.Run("NonExistentKB", func(t *testing.T) {
		req := &domain.GetNodeListReq{
			KBID: "non-existent-kb",
		}
		
		nodes, err := repo.GetList(ctx, req)
		if err != nil {
			t.Fatalf("GetList for non-existent KB failed: %v", err)
		}
		
		if len(nodes) != 0 {
			t.Errorf("Expected 0 nodes for non-existent KB, got %d", len(nodes))
		}
	})
}

// TestNodeRepository_GetList_FieldMapping 专门测试字段映射修复
func TestNodeRepository_GetList_FieldMapping(t *testing.T) {
	// 这个测试用例专注于验证字段映射，特别是修复的updated_at字段
	ctx := context.Background()
	logger := log.NewLogger()
	
	// 创建测试数据库连接
	db, err := pg.NewDB(pg.Config{
		DSN: "host=localhost user=postgres password=postgres dbname=panda_wiki_test port=5432 sslmode=disable",
	})
	if err != nil {
		t.Skip("Skipping test due to database connection failure: " + err.Error())
	}
	
	repo := NewNodeRepository(db, logger)
	
	// 清理测试数据
	defer func() {
		db.Exec("DELETE FROM nodes WHERE kb_id = ?", "field-mapping-test-kb")
		db.Exec("DELETE FROM users WHERE id = ?", "field-mapping-test-user")
	}()
	
	// 插入测试用户
	db.Exec(
		"INSERT INTO users (id, account) VALUES (?, ?)",
		"field-mapping-test-user", "test-user",
	)
	
	// 插入测试节点，设置明确的时间值
	createdTime := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	updatedTime := time.Date(2024, 1, 1, 11, 0, 0, 0, time.UTC)
	
	db.Exec(
		"INSERT INTO nodes (id, kb_id, type, status, name, content, parent_id, position, creator_id, editor_id, created_at, updated_at, meta) " +
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		"field-mapping-test-node", "field-mapping-test-kb", domain.NodeTypeDoc, domain.NodeStatusActive,
		"Field Mapping Test", "Content for field mapping test", "", 0.0,
		"field-mapping-test-user", "field-mapping-test-user",
		createdTime, updatedTime, 
		`{\"summary\": \"Field Mapping Test Summary\", \"emoji\": \"✅\", \"content_type\": \"md\"}`,
	)
	
	// 执行查询
	req := &domain.GetNodeListReq{
		KBID: "field-mapping-test-kb",
	}
	
	nodes, err := repo.GetList(ctx, req)
	if err != nil {
		t.Fatalf("GetList failed: %v", err)
	}
	
	if len(nodes) != 1 {
		t.Fatalf("Expected 1 node, got %d", len(nodes))
	}
	
	node := nodes[0]
	
	// 验证时间字段映射正确
	// 注意：由于数据库可能有时区转换，我们比较时间是否接近而不是完全相等
	if !timesAreClose(node.CreatedAt, createdTime) {
		t.Errorf("CreatedAt mapping incorrect. Expected: %v, Got: %v", createdTime, node.CreatedAt)
	}
	
	if !timesAreClose(node.UpdatedAt, updatedTime) {
		t.Errorf("UpdatedAt mapping incorrect. Expected: %v, Got: %v", updatedTime, node.UpdatedAt)
	}
}

// timesAreClose 检查两个时间是否接近（考虑数据库时区转换等因素）
func timesAreClose(t1, t2 time.Time) bool {
	diff := t1.Sub(t2)
	if diff < 0 {
		diff = -diff
	}
	// 允许最多1秒的差异
	return diff <= time.Second
}
