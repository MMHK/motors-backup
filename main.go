package main

import (
	"fmt"
	"motors-backup/internal"
	"motors-backup/internal/config"
	"os"
	"strings"
)

func main() {
	// 检查参数
	if len(os.Args) < 2 {
		fmt.Println("Usage: dump [options] table")
		fmt.Println()
		fmt.Println("Options:")
		fmt.Println("  -h, --help     Show help information")
		fmt.Println()
		fmt.Println("Arguments:")
		fmt.Println("  table          Table name(s) to export data from, multiple names separated by commas")
		os.Exit(1)
	}

	// 处理帮助参数
	if os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Println("Usage: dump [options] table")
		fmt.Println()
		fmt.Println("Options:")
		fmt.Println("  -h, --help     Show help information")
		fmt.Println()
		fmt.Println("Arguments:")
		fmt.Println("  table          Table name(s) to export data from, multiple names separated by commas")
		os.Exit(0)
	}

	// 获取表名
	tableNames := strings.Split(os.Args[1], ",")
	cfg := config.LoadConfig()

	// 执行导出操作
	for _, tableName := range tableNames {
		tableName = strings.TrimSpace(tableName)
		if tableName != "" {
			err := internal.DumpTable(cfg, tableName)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error dumping table %s: %v\n", tableName, err)
				os.Exit(1)
			}
		}
	}
}
