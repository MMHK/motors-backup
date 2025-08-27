package schema

import (
	"database/sql"
	"fmt"
)

// Column represents a table column
type Column struct {
	Name        string
	IsGenerated bool
}

// AnalyzeColumns analyzes the table schema to identify virtual columns
func AnalyzeColumns(db *sql.DB, dbName, tableName string) ([]Column, error) {
	// 使用 information_schema.COLUMNS 表查询列信息和生成表达式
	query := `
		SELECT COLUMN_NAME, 
		       CASE 
		           WHEN EXTRA LIKE '%GENERATED%' THEN 1
		           ELSE 0
		       END AS IS_VIRTUAL
		FROM information_schema.COLUMNS 
		WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? 
		ORDER BY ORDINAL_POSITION
	`

	rows, err := db.Query(query, dbName, tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to query information_schema: %w", err)
	}
	defer rows.Close()

	var columns []Column
	for rows.Next() {
		var col Column
		var isGenerated int
		err := rows.Scan(&col.Name, &isGenerated)
		if err != nil {
			return nil, fmt.Errorf("failed to scan column info: %w", err)
		}
		col.IsGenerated = isGenerated == 1
		columns = append(columns, col)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return columns, nil
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
