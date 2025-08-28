package schema

import (
	"database/sql"
	"fmt"
	"strings"
)

// Column represents a table column
type Column struct {
	Name        string
	IsGenerated bool
}

// AnalyzeColumns analyzes the table schema to identify virtual columns
func AnalyzeColumns(db *sql.DB, dbName, tableName string) ([]Column, error) {
	// 使用 information_schema.COLUMNS 表查询列信息和生成表达式
	query := fmt.Sprintf("SHOW COLUMNS FROM `%s`.`%s`;", dbName, tableName)

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query information_schema: %w", err)
	}
	defer rows.Close()

	var columns []Column
	for rows.Next() {
		var tmp interface{}
		var col Column
		var Extra string
		err := rows.Scan(&col.Name, &tmp, &tmp, &tmp, &tmp, &Extra)
		if err != nil {
			return nil, fmt.Errorf("failed to scan column info: %w", err)
		}
		col.IsGenerated = strings.Contains(Extra, "GENERATED")
		columns = append(columns, col)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return columns, nil
}

func GetTableDDL(db *sql.DB, dbName string, tableName string) (string, error) {
	query := fmt.Sprintf("SHOW CREATE TABLE `%s`.`%s`;", dbName, tableName)
	rows, err := db.Query(query)
	if err != nil {
		return "", fmt.Errorf("failed to query information_schema: %w", err)
	}
	defer rows.Close()
	var ddl string
	for rows.Next() {
		var tmp interface{}
		err := rows.Scan(&tmp, &ddl)
		if err != nil {
			return "", fmt.Errorf("failed to scan column info: %w", err)
		}
	}
	return ddl, nil
}

// GetNonGeneratedColumns returns a list of non-virtual column names
func GetNonGeneratedColumns(columns []Column) []string {
	var nonVirtualColumns []string
	for _, col := range columns {
		if !col.IsGenerated {
			nonVirtualColumns = append(nonVirtualColumns, col.Name)
		}
	}
	return nonVirtualColumns
}

func GetDatabaseDDL(db *sql.DB, dbName string) (string, error) {
	query := fmt.Sprintf("SHOW CREATE DATABASE `%s`;", dbName)
	rows, err := db.Query(query)
	if err != nil {
		return "", fmt.Errorf("failed to query information_schema: %w", err)
	}
	defer rows.Close()
	var ddl string
	for rows.Next() {
		var tmp interface{}
		err := rows.Scan(&tmp, &ddl)
		if err != nil {
			return "", fmt.Errorf("failed to scan column info: %w", err)
		}
	}
	return ddl, nil
}

func ListAllTables(db *sql.DB) ([]string, error) {
	query := "SHOW TABLES;"
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query information_schema: %w", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		err := rows.Scan(&tableName)
		if err != nil {
			return nil, fmt.Errorf("failed to scan table name: %w", err)
		}
		tables = append(tables, tableName)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return tables, nil
}
