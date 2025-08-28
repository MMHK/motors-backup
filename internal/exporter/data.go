package exporter

import (
	"database/sql"
	"fmt"
	"strings"
)

// ExportData exports table data as INSERT statements
func ExportData(db *sql.DB, tableName string, columns []string) error {
	// 构建查询语句
	columnList := "`" + strings.Join(columns, "`, `") + "`"
	query := fmt.Sprintf("SELECT %s FROM `%s`", columnList, tableName)

	rows, err := db.Query(query)
	if err != nil {
		return fmt.Errorf("failed to query table data: %w", err)
	}
	defer rows.Close()

	// 获取列信息
	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return fmt.Errorf("failed to get column types: %w", err)
	}

	// 输出表头信息
	fmt.Printf("--\n-- Dumping data for table `%s`\n--\n\n", tableName)

	// 准备用于Scan的值
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	// 遍历每一行数据
	for rows.Next() {
		// Scan数据
		err := rows.Scan(valuePtrs...)
		if err != nil {
			return fmt.Errorf("failed to scan row: %w", err)
		}

		// 构建INSERT语句
		insertStmt := buildInsertStatement(tableName, columns, values, columnTypes)
		fmt.Println(insertStmt)
	}

	if err = rows.Err(); err != nil {
		return fmt.Errorf("error iterating rows: %w", err)
	}

	return nil
}

// buildInsertStatement builds an INSERT statement for a row of data
func buildInsertStatement(tableName string, columns []string, values []interface{}, columnTypes []*sql.ColumnType) string {
	// 构建列名部分
	columnList := "`" + strings.Join(columns, "`, `") + "`"

	// 构建值部分
	var valueList []string
	for i, value := range values {
		// 根据列类型处理值
		strValue := formatValue(value, columnTypes[i])
		valueList = append(valueList, strValue)
	}

	// 组合INSERT语句
	return fmt.Sprintf("INSERT INTO `%s` (%s) VALUES (%s);", tableName, columnList, strings.Join(valueList, ", "))
}

// EscapeSQLString 接收一个字符串，并返回一个符合SQL字面量规范的安全字符串。
// 它会用单引号包裹结果，并对内部的特殊字符进行转义。
func EscapeSQLString(value string) string {
	var sb strings.Builder
	// SQL字符串以单引号开始
	sb.WriteByte('\'')

	for _, r := range value {
		switch r {
		case '\'':
			// 将单引号转义为两个单引号
			sb.WriteString("''")
		case '\\':
			// 将反斜杠转义为两个反斜杠
			sb.WriteString("\\\\")
		case '\n':
			// 将换行符转义为 \n
			sb.WriteString("\\n")
		case '\r':
			// 将回车符转义为 \r
			sb.WriteString("\\r")
		case '\t':
			// 将制表符转义为 \t
			sb.WriteString("\\t")
		case '\x00':
			// 将空字节转义为 \0
			sb.WriteString("\\0")
		case '\x1a':
			sb.WriteString("\\Z") // Ctrl+Z
		default:
			sb.WriteRune(r)
		}
	}

	// SQL字符串以单引号结束
	sb.WriteByte('\'')
	return sb.String()
}

// formatValue formats a value for use in an SQL statement
func formatValue(value interface{}, columnType *sql.ColumnType) string {
	if value == nil {
		return "NULL"
	}

	switch v := value.(type) {
	case int64:
		return fmt.Sprintf("'%d'", v)
	case float64:
		return fmt.Sprintf("'%.2f'", v) // 保留 2 位小數
	}

	// 字符串类型需要引号
	strValue := fmt.Sprintf("%s", value)
	return EscapeSQLString(strValue)
}
