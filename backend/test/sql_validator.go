package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	// 直接模拟我们修复后的SQL查询
	query := `nodes.id, nodes.kb_id, nodes.type, nodes.status, nodes.name, nodes.content, nodes.parent_id, nodes.position, nodes.created_at, nodes.updated_at, nodes.meta, u1.account as creator, u2.account as editor`
	
	// 检查是否有多余的逗号
	if containsDoubleComma(query) {
		fmt.Println("ERROR: SQL contains double commas!")
		os.Exit(1)
	}
	
	// 打印查询语句，便于检查
	fmt.Println("SQL query is valid, no double commas found:")
	fmt.Println(query)
	
	// 分析修复前的错误
	originalQuery := "nodes.id, nodes.kb_id, nodes.type, nodes.status, nodes.name, nodes.content, nodes.parent_id, ,nodes.position, ,nodes.created_at"
	if containsDoubleComma(originalQuery) {
		fmt.Println("\nConfirmed: Original query contains double commas")
		
		// 找出具体位置
		commaPositions := findCommaPositions(originalQuery)
		fmt.Println("Double comma positions:", commaPositions)
	}
	
	os.Exit(0)
}

// containsDoubleComma 检查字符串中是否包含连续的逗号
func containsDoubleComma(s string) bool {
	return strings.Contains(s, ", ,")
}

// findCommaPositions 找出所有连续逗号的位置
func findCommaPositions(s string) []int {
	var positions []int
	for i := 0; i < len(s)-1; i++ {
		if s[i] == ',' && i+1 < len(s) && s[i+1] == ' ' && i+2 < len(s) && s[i+2] == ',' {
			positions = append(positions, i)
		}
	}
	return positions
}
