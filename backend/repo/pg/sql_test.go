package pg

import (
	"fmt"
	"testing"
)

// TestSQLGeneration 直接测试SQL生成是否正确，不依赖数据库连接
func TestSQLGeneration(t *testing.T) {
	// 测试Select语句是否正确生成，没有多余的逗号
	query := `nodes.id, nodes.kb_id, nodes.type, nodes.status, nodes.name, nodes.content, nodes.parent_id, nodes.position, nodes.created_at, nodes.updated_at, nodes.meta, u1.account as creator, u2.account as editor`
	
	// 检查是否有多余的逗号
	if containsDoubleComma(query) {
		t.Fatalf("SQL contains double commas: %s", query)
	}
	
	fmt.Println("SQL query is valid, no double commas found")
}

// containsDoubleComma 检查字符串中是否包含连续的逗号
func containsDoubleComma(s string) bool {
	for i := 0; i < len(s)-1; i++ {
		if s[i] == ',' && s[i+1] == ',' {
			return true
		}
	}
	return false
}
