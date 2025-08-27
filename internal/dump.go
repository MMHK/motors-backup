package internal

import (
	"fmt"
	"motors-backup/internal/config"
	"motors-backup/internal/db"
	"motors-backup/internal/exporter"
	"motors-backup/internal/schema"
)

// DumpTable dumps the specified table data as SQL INSERT statements
func DumpTable(cfg *config.Config, tableName string) error {

	// 检查必需的DB_NAME配置
	if cfg.DBName == "" {
		return fmt.Errorf("DB_NAME environment variable is required")
	}

	// 建立数据库连接
	database, err := db.Connect(cfg)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close(database)

	// 输出基本环境设置
	printEnvironmentSettings(cfg)

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
	err = exporter.ExportData(database, tableName, nonGeneratedColumns)
	if err != nil {
		return fmt.Errorf("failed to export data: %w", err)
	}

	return nil
}

// printEnvironmentSettings outputs basic MySQL environment settings
func printEnvironmentSettings(cfg *config.Config) {
	fmt.Println("-- MySQL dump 10.13  Distrib 8.0.x, for Linux (x86_64)")
	fmt.Println("--")
	fmt.Printf("-- Host: %s    Database: %s\n", cfg.DBHost, cfg.DBName)
	fmt.Println("-- ------------------------------------------------------")
	fmt.Println("-- Server version	8.0.x")
	fmt.Println()
	fmt.Println("/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;")
	fmt.Println("/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;")
	fmt.Println("/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;")
	fmt.Println("/*!50503 SET NAMES utf8mb4 */;")
	fmt.Println("/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;")
	fmt.Println("/*!40103 SET TIME_ZONE='+00:00' */;")
	fmt.Println("/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;")
	fmt.Println("/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;")
	fmt.Println("/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;")
	fmt.Println("/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;")
	fmt.Println()
}
