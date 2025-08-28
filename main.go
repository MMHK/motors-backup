package main

import (
	"database/sql"
	"flag"
	"fmt"
	"motors-backup/internal"
	"motors-backup/internal/config"
	"motors-backup/internal/log"
	"motors-backup/internal/schema"
	"os"
	"strings"
)

// parseFlags 处理命令行参数解析
func parseFlags() (tableNames []string, createDatabase bool, ignoreTables ignoreList, ignoreTableDataList ignoreList, whereCondition string, err error) {
	// 定义命令行参数
	help := flag.Bool("help", false, "Show help information")
	h := flag.Bool("h", false, "Show help information")
	createDatabaseFlag := flag.Bool("create-database", true, "Enable create database statement")

	// 定义可重复使用的ignore-table和ignore-table-data参数
	flag.Var(&ignoreTables, "ignore-table", "Table name(s) to ignore structure and data, can be specified multiple times")
	flag.Var(&ignoreTableDataList, "ignore-table-data", "Table name(s) to ignore data only, can be specified multiple times")

	// 定义where条件参数
	whereFlag := flag.String("where", "", "WHERE condition for querying table data")

	flag.Usage = func() {
		fmt.Println("Usage: motors-backup [options] table")
		fmt.Println()
		fmt.Println("Options:")
		flag.PrintDefaults()
		fmt.Println()
		fmt.Println("Arguments:")
		fmt.Println("  table          Table name(s) to export data from, multiple names separated by commas")
		fmt.Println()
		fmt.Println("Environment Variables:")
		fmt.Println("  DB_HOST=")
		fmt.Println("  DB_PORT=")
		fmt.Println("  DB_USER=root")
		fmt.Println("  DB_PASSWORD=")
		fmt.Println("  DB_NAME=")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  motors-backup                            Export all tables in the database")
		fmt.Println("  motors-backup users                      Export users table")
		fmt.Println("  motors-backup users,orders               Export users and orders tables")
		fmt.Println("  motors-backup --ignore-table=logs --ignore-table=users")
		fmt.Println("                                           Export all tables except logs and users")
		fmt.Println("  motors-backup --ignore-table-data=logs --ignore-table-data=users users,orders")
		fmt.Println("                                           Export users table structure and data,")
		fmt.Println("                                           but only structure for logs table")
		fmt.Println("  motors-backup --create-database=false users")
		fmt.Println("                                           Export users table without create database statement")
		fmt.Println("  motors-backup --where='id>100' users")
		fmt.Println("                                           Export users table with condition id>100")
	}

	flag.Parse()

	// 处理帮助参数
	if *help || *h {
		return tableNames, createDatabase, ignoreTables, ignoreTableDataList, whereCondition, nil
	}

	// 检查参数
	if flag.NArg() > 0 {
		// 获取表名
		tableNames = strings.Split(flag.Arg(0), ",")
	}

	createDatabase = *createDatabaseFlag
	whereCondition = *whereFlag

	return tableNames, createDatabase, ignoreTables, ignoreTableDataList, whereCondition, nil
}

func main() {
	// 解析命令行参数
	tableNames, createDatabase, ignoreTables, ignoreTableDataList, whereCondition, err := parseFlags()
	if err != nil {
		flag.Usage()
		os.Exit(1)
	}

	cfg := config.LoadConfig()

	err = internal.StartExport(cfg, func(database *sql.DB, info *internal.MySQLInfo) error {
		internal.PrintEnvironmentSettings(cfg, info)

		// 如果启用了create-database参数，则执行创建数据库操作
		err := internal.DumpCreateDatabase(cfg, database, createDatabase)
		if err != nil {
			log.Logger.Errorf("Error creating database: %v\n", err)
			os.Exit(1)
			return err
		}

		allTables, err := schema.ListAllTables(database)
		if err != nil {
			log.Logger.Errorf("Error listing all tables: %v\n", err)
			os.Exit(1)
			return err
		}

		if len(tableNames) == 0 {
			// 如果没有指定表名，则导出所有表
			tableNames = allTables
		}
		// 如果 tableNames 不wei空 filter 掉不在 allTables 的 table
		filteredTables := make([]string, 0)
		for _, tableName := range tableNames {
			for _, allTable := range allTables {
				if tableName == allTable {
					filteredTables = append(filteredTables, tableName)
					break
				}
			}
		}
		if len(filteredTables) == 0 {
			log.Logger.Errorf("No tables found in database:%s", cfg.DBName)
			os.Exit(1)
		}
		tableNames = filteredTables

		// 执行导出操作
		for _, tableName := range tableNames {
			tableName = strings.TrimSpace(tableName)
			if tableName != "" {
				// 检查是否在忽略表列表中
				if ignoreTables.Contains(tableName) {
					continue
				}

				// 如果不在忽略结构列表中，则导出表结构
				err := internal.DumpTableStructure(cfg, database, tableName)
				if err != nil {
					log.Logger.Errorf("Error dumping table structure %s: %v\n", tableName, err)
					os.Exit(1)
					return err
				}

				// 如果不在忽略数据列表中，则导出表数据
				if !ignoreTableDataList.Contains(tableName) {
					err = internal.DumpTable(cfg, database, tableName, whereCondition)
					if err != nil {
						log.Logger.Errorf("Error dumping table %s: %v\n", tableName, err)
						os.Exit(1)
						return err
					}
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

// ignoreList 实现了 flag.Value 接口，用于处理可重复的参数
type ignoreList []string

func (i *ignoreList) String() string {
	return strings.Join(*i, ",")
}

func (i *ignoreList) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func (i *ignoreList) Contains(value string) bool {
	for _, item := range *i {
		if item == value {
			return true
		}
	}
	return false
}
