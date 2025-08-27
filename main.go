package main

import (
	"database/sql"
	"fmt"
	"motors-backup/internal"
	"motors-backup/internal/config"
	"motors-backup/internal/log"
	"os"
	"strings"
)

func main() {
	// 检查参数
	if len(os.Args) < 2 {
		fmt.Println("Usage: motors-backup [options] table")
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
		fmt.Println("Usage: motors-backup [options] table")
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

	err := internal.StartExport(cfg, func(database *sql.DB, info *internal.MySQLInfo) error {
		internal.PrintEnvironmentSettings(cfg, info)

		err := internal.DumpCreateDatabase(cfg, database)
		if err != nil {
			log.Logger.Errorf("Error creating database: %v\n", err)
			os.Exit(1)
			return err
		}

		// 执行导出操作
		for _, tableName := range tableNames {
			tableName = strings.TrimSpace(tableName)
			if tableName != "" {
				err := internal.DumpTableStructure(cfg, database, tableName)
				if err != nil {
					log.Logger.Errorf("Error dumping table structure %s: %v\n", tableName, err)
					os.Exit(1)
					return err
				}

				err = internal.DumpTable(cfg, database, tableName)
				if err != nil {
					log.Logger.Errorf("Error dumping table %s: %v\n", tableName, err)
					os.Exit(1)
					return err
				}
			}
		}

		internal.PrintRestoreConnectionSettings()

		return nil
	})

	if err != nil {
		log.Logger.Errorf("Error: %v\n", err)
		os.Exit(1)
	}
}
