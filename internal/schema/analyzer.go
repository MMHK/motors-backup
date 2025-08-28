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
	var tmp interface{}
	params := make([]interface{}, 0)
	columns, err := rows.Columns()
	for _, col := range columns {
		if strings.Contains(col, "Create") {
			params = append(params, &ddl)
			continue
		}
		params = append(params, &tmp)
	}

	for rows.Next() {
		err := rows.Scan(params...)
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
	query := "SHOW FULL TABLES WHERE Table_Type = 'BASE TABLE';"
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query information_schema: %w", err)
	}
	defer rows.Close()

	var tables []string
	var tmp string
	for rows.Next() {
		var tableName string
		err := rows.Scan(&tableName, &tmp)
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

type TriggerInfo struct {
	Name     string
	SQL_MODE string
	DDL      string
}

func AllTriggersDDL(db *sql.DB, dbName string) ([]*TriggerInfo, error) {
	query := fmt.Sprintf("SELECT `trigger_name` FROM `information_schema`.`triggers` WHERE `trigger_schema` = ?")
	rows, err := db.Query(query, dbName)
	if err != nil {
		return nil, fmt.Errorf("failed to query triggers: %w", err)
	}
	defer rows.Close()

	var triggers []string
	for rows.Next() {
		var triggerName string
		err := rows.Scan(&triggerName)
		if err != nil {
			return nil, fmt.Errorf("failed to scan trigger name: %w", err)
		}
		triggers = append(triggers, triggerName)
	}

	triggersDDL := make([]*TriggerInfo, 0)

	for _, triggerName := range triggers {
		query := fmt.Sprintf("SHOW CREATE TRIGGER `%s`.`%s`", dbName, triggerName)
		rows, err := db.Query(query)
		if err != nil {
			return nil, fmt.Errorf("failed to query trigger DDL: %w", err)
		}
		defer rows.Close()

		columns, err := rows.Columns()
		if err != nil {
			return nil, fmt.Errorf("failed to get column names: %w", err)
		}
		params := make([]interface{}, len(columns))
		info := new(TriggerInfo)
		var tmp string
		for i, col := range columns {
			if strings.Contains(col, "sql_mode") {
				params[i] = &(info.SQL_MODE)
				continue
			}
			if strings.Contains(col, "Trigger") {
				params[i] = &(info.Name)
				continue
			}
			if strings.Contains(col, "Statement") {
				params[i] = &(info.DDL)
				continue
			}
			params[i] = &tmp
		}

		for rows.Next() {
			err := rows.Scan(params...)
			if err != nil {
				return nil, fmt.Errorf("failed to scan trigger DDL: %w", err)
			}
			triggersDDL = append(triggersDDL, info)
		}
	}

	return triggersDDL, nil
}

type ViewInfo struct {
	DDL  string
	Name string
}

func AllViewDDL(db *sql.DB) ([]*ViewInfo, error) {
	query := "SHOW FULL TABLES WHERE Table_Type = 'VIEW';"
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query view DDL: %w", err)
	}
	defer rows.Close()

	views := make([]string, 0)
	viewDDL := make([]*ViewInfo, 0)

	for rows.Next() {
		var viewName, tableType string
		err := rows.Scan(&viewName, &tableType)
		if err != nil {
			return nil, fmt.Errorf("failed to scan view DDL: %w", err)
		}
		views = append(views, viewName)
	}

	for _, view := range views {
		query := fmt.Sprintf("SHOW CREATE VIEW %s", view)
		rows, err := db.Query(query)
		if err != nil {
			return nil, fmt.Errorf("failed to query view DDL: %w", err)
		}
		defer rows.Close()

		columns, err := rows.Columns()
		if err != nil {
			return nil, fmt.Errorf("failed to get view DDL columns: %w", err)
		}
		var tmp string
		var viewName string
		var DDL string
		params := make([]interface{}, len(columns))
		for i, col := range columns {
			if strings.Contains(col, "Create") {
				params[i] = &DDL
				continue
			}
			if strings.Contains(col, "View") {
				params[i] = &viewName
				continue
			}
			params[i] = &tmp
		}

		for rows.Next() {
			err := rows.Scan(params...)
			if err != nil {
				return nil, fmt.Errorf("failed to scan view DDL: %w", err)
			}
			viewDDL = append(viewDDL, &ViewInfo{
				DDL:  DDL,
				Name: viewName,
			})
		}
	}

	return viewDDL, nil
}
