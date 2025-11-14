package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	migratePG "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	// 获取数据库连接字符串
	dsn := os.Getenv("PG_DSN")
	if dsn == "" {
		fmt.Println("环境变量PG_DSN未设置")
		os.Exit(1)
	}

	// 连接数据库
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		fmt.Printf("连接数据库失败: %v\n", err)
		os.Exit(1)
	}
	
	driver, err := migratePG.WithInstance(db, &migratePG.Config{})
	if err != nil {
		fmt.Printf("创建迁移驱动失败: %v\n", err)
		os.Exit(1)
	}

	// 初始化迁移
	m, err := migrate.NewWithDatabaseInstance(
		"file://store/pg/migration",
		"postgres", driver)
	if err != nil {
		fmt.Printf("初始化迁移失败: %v\n", err)
		os.Exit(1)
	}

	// 运行迁移
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			fmt.Println("没有需要执行的迁移")
			err = nil
		} else {
			fmt.Printf("执行迁移失败: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Println("迁移执行成功")
	}
}
