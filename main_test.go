package main

import (
	"flag"
	"os"
	"strings"
	"testing"
)

func TestFlagParsing(t *testing.T) {
	// 保存原始的os.Args
	oldArgs := os.Args

	// 恢复原始状态
	defer func() {
		os.Args = oldArgs
	}()

	// 测试用例
	testCases := []struct {
		name                    string
		args                    []string
		expectedIgnoreTables    []string
		expectedIgnoreTableData []string
		expectedCreateDatabase  bool
		expectedTableNames      []string
		expectedWhereCondition  string
	}{
		{
			name:                    "basic table export",
			args:                    []string{"motors-backup", "users"},
			expectedIgnoreTables:    []string{},
			expectedIgnoreTableData: []string{},
			expectedCreateDatabase:  true,
			expectedTableNames:      []string{"users"},
			expectedWhereCondition:  "",
		},
		{
			name:                    "multiple tables export",
			args:                    []string{"motors-backup", "users,orders"},
			expectedIgnoreTables:    []string{},
			expectedIgnoreTableData: []string{},
			expectedCreateDatabase:  true,
			expectedTableNames:      []string{"users", "orders"},
			expectedWhereCondition:  "",
		},
		{
			name:                    "ignore table structure and data",
			args:                    []string{"motors-backup", "--ignore-table=logs", "users,logs"},
			expectedIgnoreTables:    []string{"logs"},
			expectedIgnoreTableData: []string{},
			expectedCreateDatabase:  true,
			expectedTableNames:      []string{"users", "logs"},
			expectedWhereCondition:  "",
		},
		{
			name:                    "ignore table data only",
			args:                    []string{"motors-backup", "--ignore-table-data=logs", "users,logs"},
			expectedIgnoreTables:    []string{},
			expectedIgnoreTableData: []string{"logs"},
			expectedCreateDatabase:  true,
			expectedTableNames:      []string{"users", "logs"},
			expectedWhereCondition:  "",
		},
		{
			name:                    "disable create database",
			args:                    []string{"motors-backup", "--create-database=false", "users"},
			expectedIgnoreTables:    []string{},
			expectedIgnoreTableData: []string{},
			expectedCreateDatabase:  false,
			expectedTableNames:      []string{"users"},
			expectedWhereCondition:  "",
		},
		{
			name:                    "multiple ignore flags",
			args:                    []string{"motors-backup", "--ignore-table=logs", "--ignore-table=temp", "--ignore-table-data=sessions", "users,logs,temp,sessions,orders"},
			expectedIgnoreTables:    []string{"logs", "temp"},
			expectedIgnoreTableData: []string{"sessions"},
			expectedCreateDatabase:  true,
			expectedTableNames:      []string{"users", "logs", "temp", "sessions", "orders"},
			expectedWhereCondition:  "",
		},
		{
			name:                    "where condition",
			args:                    []string{"motors-backup", "--where=id>100", "users"},
			expectedIgnoreTables:    []string{},
			expectedIgnoreTableData: []string{},
			expectedCreateDatabase:  true,
			expectedTableNames:      []string{"users"},
			expectedWhereCondition:  "id>100",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 重置flag状态以避免"flag redefined"错误
			flag.CommandLine = flag.NewFlagSet(tc.args[0], flag.ExitOnError)

			// 设置命令行参数
			os.Args = tc.args

			// 调用parseFlags函数
			tableNames, createDatabase, ignoreTables, ignoreTableDataList, whereCondition, err := parseFlags()
			if err != nil {
				t.Fatalf("parseFlags returned error: %v", err)
			}

			// 验证create-database标志
			if createDatabase != tc.expectedCreateDatabase {
				t.Errorf("createDatabase = %v, want %v", createDatabase, tc.expectedCreateDatabase)
			}

			// 验证ignore-table标志
			if len(ignoreTables) != len(tc.expectedIgnoreTables) {
				t.Errorf("ignoreTables length = %d, want %d", len(ignoreTables), len(tc.expectedIgnoreTables))
			} else {
				for i, expected := range tc.expectedIgnoreTables {
					if ignoreTables[i] != expected {
						t.Errorf("ignoreTables[%d] = %s, want %s", i, ignoreTables[i], expected)
					}
				}
			}

			// 验证ignore-table-data标志
			if len(ignoreTableDataList) != len(tc.expectedIgnoreTableData) {
				t.Errorf("ignoreTableDataList length = %d, want %d", len(ignoreTableDataList), len(tc.expectedIgnoreTableData))
			} else {
				for i, expected := range tc.expectedIgnoreTableData {
					if ignoreTableDataList[i] != expected {
						t.Errorf("ignoreTableDataList[%d] = %s, want %s", i, ignoreTableDataList[i], expected)
					}
				}
			}

			// 验证表名解析
			if len(tableNames) != len(tc.expectedTableNames) {
				t.Errorf("tableNames length = %d, want %d", len(tableNames), len(tc.expectedTableNames))
			} else {
				for i, expected := range tc.expectedTableNames {
					if tableNames[i] != expected {
						t.Errorf("tableNames[%d] = %s, want %s", i, tableNames[i], expected)
					}
				}
			}

			// 验证where条件
			if !strings.EqualFold(whereCondition, tc.expectedWhereCondition) {
				t.Errorf("whereCondition = %s, want %s", whereCondition, tc.expectedWhereCondition)
			}
		})
	}
}

func TestIgnoreListContains(t *testing.T) {
	il := ignoreList{"users", "logs"}

	if !il.Contains("users") {
		t.Error("Expected ignoreList to contain 'users'")
	}

	if !il.Contains("logs") {
		t.Error("Expected ignoreList to contain 'logs'")
	}

	if il.Contains("orders") {
		t.Error("Expected ignoreList not to contain 'orders'")
	}
}
