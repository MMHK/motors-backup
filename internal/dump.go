package internal

import (
	"database/sql"
	"fmt"
	"motors-backup/internal/config"
	dbConn "motors-backup/internal/db"
	"motors-backup/internal/exporter"
	"motors-backup/internal/schema"
	"regexp"
	"runtime"
	"strings"
)

func DumpCreateDatabase(cfg *config.Config, database *sql.DB, withCreateDB bool) error {
	databaseDDL, err := schema.GetDatabaseDDL(database, cfg.DBName)
	if err != nil {
		return fmt.Errorf("failed to get database DDL: %w", err)
	}
	fmt.Println("--")
	fmt.Printf("-- Current Database: `%s`\n", cfg.DBName)
	fmt.Println("--")
	if withCreateDB {
		fmt.Printf("\n%s;\n", strings.Replace(databaseDDL, "CREATE DATABASE", "CREATE DATABASE /*!32312 IF NOT EXISTS*/", 1))
	}
	fmt.Printf("\nUSE `%s`;\n\n", cfg.DBName)

	return nil
}

func DumpTableStructure(cfg *config.Config, database *sql.DB, tableName string) error {
	tableDDL, err := schema.GetTableDDL(database, cfg.DBName, tableName)
	if err != nil {
		return fmt.Errorf("failed to get table DDL: %w", err)
	}
	fmt.Println("--")
	fmt.Printf("-- Table structure for table `%s`\n", tableName)
	fmt.Println("--")
	fmt.Println("")
	fmt.Println("")
	fmt.Printf("DROP TABLE IF EXISTS `%s`;\n", tableName)
	fmt.Println("/*!40101 SET @saved_cs_client     = @@character_set_client */;")
	fmt.Println("/*!40101 SET character_set_client = utf8mb4 */;")
	fmt.Printf("%s;\n", tableDDL)
	fmt.Printf("/*!40101 SET character_set_client = @saved_cs_client */;\n\n")
	return nil
}

// DumpTable dumps the specified table data as SQL INSERT statements
func DumpTable(cfg *config.Config, database *sql.DB, tableName string, whereClause string) error {

	// 获取MySQL服务器信息
	mysqlInfo, err := getMySQLInfo(database)
	if err != nil {
		return fmt.Errorf("failed to get MySQL info: %w", err)
	}

	// 检查MySQL版本兼容性
	if err := checkMySQLVersionCompatibility(mysqlInfo.Version); err != nil {
		return err
	}

	// 分析列结构，识别虚拟列
	columns, err := schema.AnalyzeColumns(database, cfg.DBName, tableName)
	if err != nil {
		return fmt.Errorf("failed to analyze columns: %w", err)
	}

	// 获取非虚拟列列表
	nonGeneratedColumns := schema.GetNonGeneratedColumns(columns)
	if len(nonGeneratedColumns) == 0 {
		return fmt.Errorf("no non-generated columns found in %s", tableName)
	}

	// 导出数据
	err = exporter.ExportData(database, tableName, nonGeneratedColumns, whereClause)
	if err != nil {
		return fmt.Errorf("failed to export data: %w", err)
	}

	return nil
}

func ReplaceDDLDefinerWithCurrentUser(ddl string) string {
	// 替换 DEFINER=...
	re := regexp.MustCompile(`DEFINER=[^ ]+`)
	return re.ReplaceAllString(ddl, "DEFINER=CURRENT_USER()")
}

func ReplaceViewDDLASReplace(ddl string) string {
	// 替换 AS REPLACE
	re := regexp.MustCompile(`CREATE ALGORITHM`)
	return re.ReplaceAllString(ddl, "CREATE OR REPLACE ALGORITHM")
}

func DumpViews(cfg *config.Config, database *sql.DB) error {
	viewDDLs, err := schema.AllViewDDL(database)
	if err != nil {
		return fmt.Errorf("failed to get view DDL: %w", err)
	}
	for _, viewDDL := range viewDDLs {
		fmt.Printf("\n--\n-- Temporary table structure for view `%s`\n--\n\n", viewDDL.Name)
		fmt.Printf("/*!50001 DROP VIEW IF EXISTS `%s`*/;\n", viewDDL.Name)
		fmt.Println("SET @saved_cs_client     = @@character_set_client;")
		fmt.Println("SET character_set_client = utf8mb4;")
		fmt.Printf("%s;\n", ReplaceDDLDefinerWithCurrentUser(ReplaceViewDDLASReplace(viewDDL.DDL)))
		fmt.Printf("SET character_set_client = @saved_cs_client;\n")
	}

	return nil
}

// MySQLInfo 存储MySQL服务器信息
type MySQLInfo struct {
	Version  string
	Charset  string
	Timezone string
}

func StartExport(cfg *config.Config, worker func(database *sql.DB, info *MySQLInfo) error) error {
	// 检查必需的DB_NAME配置
	if cfg.DBName == "" {
		return fmt.Errorf("DB_NAME environment variable is required")
	}

	// 建立数据库连接
	database, err := dbConn.Connect(cfg)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer dbConn.Close(database)

	// 获取MySQL服务器信息
	mysqlInfo, err := getMySQLInfo(database)
	if err != nil {
		return fmt.Errorf("failed to get MySQL info: %w", err)
	}

	// 检查MySQL版本兼容性
	if err := checkMySQLVersionCompatibility(mysqlInfo.Version); err != nil {
		return err
	}

	return worker(database, mysqlInfo)
}

// getMySQLInfo 获取MySQL服务器的版本、字符集和时区信息
func getMySQLInfo(database *sql.DB) (*MySQLInfo, error) {
	var version, charset, timezone string

	// 获取版本信息
	row := database.QueryRow("SELECT VERSION()")
	if err := row.Scan(&version); err != nil {
		return nil, fmt.Errorf("failed to get MySQL version: %w", err)
	}

	// 获取字符集信息
	row = database.QueryRow("SELECT @@character_set_server")
	if err := row.Scan(&charset); err != nil {
		return nil, fmt.Errorf("failed to get MySQL charset: %w", err)
	}

	// 获取时区信息
	row = database.QueryRow("SELECT TIME_FORMAT(TIMEDIFF(NOW(), UTC_TIMESTAMP()), '%H:%i') AS timezone_offset;")
	if err := row.Scan(&timezone); err != nil {
		return nil, fmt.Errorf("failed to get MySQL timezone: %w", err)
	}

	if !strings.Contains(timezone, "-") {
		timezone = fmt.Sprintf("+%s", timezone)
	}

	return &MySQLInfo{
		Version:  version,
		Charset:  charset,
		Timezone: timezone,
	}, nil
}

// checkMySQLVersionCompatibility 检查MySQL版本是否兼容（>= 8.0）
func checkMySQLVersionCompatibility(version string) error {
	// 简单检查版本号是否以8.0或更高版本开头
	if strings.HasPrefix(version, "8.") || strings.HasPrefix(version, "9.") {
		return nil
	}

	// 如果版本号小于8.0，则返回错误
	if strings.HasPrefix(version, "5.") || strings.HasPrefix(version, "6.") || strings.HasPrefix(version, "7.") {
		return fmt.Errorf("incompatible MySQL version: %s. Required version >= 8.0", version)
	}

	// 对于其他格式的版本号，假定为兼容
	return nil
}

// PrintEnvironmentSettings outputs basic MySQL environment settings
func PrintEnvironmentSettings(cfg *config.Config, mysqlInfo *MySQLInfo) {
	// 獲取 golang runtime 執行環境arc
	arch := runtime.GOARCH
	os := runtime.GOOS

	fmt.Printf("-- MOTORS_BACKUP 0.1  Distrib 8.0.x, for %s (%s)\n", os, arch)
	fmt.Println("--")
	fmt.Printf("-- Host: %s    Database: %s\n", cfg.DBHost, cfg.DBName)
	fmt.Println("-- ------------------------------------------------------")
	fmt.Printf("-- Server version	%s\n", mysqlInfo.Version)
	fmt.Println()
	fmt.Println("/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;")
	fmt.Println("/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;")
	fmt.Println("/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;")
	fmt.Println("/*!50503 SET NAMES utf8mb4 */;")
	fmt.Println("/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;")
	fmt.Printf("/*!40103 SET TIME_ZONE='%s' */;\n", mysqlInfo.Timezone)
	fmt.Println("/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;")
	fmt.Println("/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;")
	fmt.Println("/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;")
	fmt.Println("/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;")
	fmt.Println()
}

// PrintRestoreConnectionSettings outputs statements to restore connection settings
func PrintRestoreConnectionSettings() {
	fmt.Println("/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;")
	fmt.Println("/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;")
	fmt.Println("/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;")
	fmt.Println("/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;")
	fmt.Println("/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;")
	fmt.Println("/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;")
	fmt.Println("/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;")
	fmt.Println("/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;")
	fmt.Println()
}
