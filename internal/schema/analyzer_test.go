package schema

import (
	"database/sql"
	"motors-backup/internal/config"
	"motors-backup/internal/db"
	"os"
	"testing"
)

func GetTestConfig() (*sql.DB, string, string, error) {
	conf := config.LoadTestConfig()

	//log.Logger.Debugf("%+v", conf)

	db, err := db.Connect(conf)
	if err != nil {
		return nil, "", "", err
	}

	testTableName := os.Getenv("TEST_TABLE")

	return db, conf.DBName, testTableName, nil
}

func TestListAllTables(t *testing.T) {
	dbConn, _, _, err := GetTestConfig()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	tables, err := ListAllTables(dbConn)
	if err != nil {
		t.Errorf("Failed to list tables: %v", err)
	}
	t.Logf("Tables: %+v", tables)
}

func TestAnalyzeColumns(t *testing.T) {
	dbConn, dbName, tableName, err := GetTestConfig()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	columns, err := AnalyzeColumns(dbConn, dbName, tableName)
	if err != nil {
		t.Errorf("Failed to analyze columns: %v", err)
	}

	t.Logf("Columns: %+v", columns)
}

func TestGetTableDDL(t *testing.T) {
	dbConn, dbName, tableName, err := GetTestConfig()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	ddl, err := GetTableDDL(dbConn, dbName, tableName)
	if err != nil {
		t.Errorf("Failed to get table DDL: %v", err)
	}

	t.Logf("Table DDL: %s", ddl)
}

func TestGetNonGeneratedColumns(t *testing.T) {
	dbConn, dbName, tableName, err := GetTestConfig()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	columns, err := AnalyzeColumns(dbConn, dbName, tableName)
	if err != nil {
		t.Errorf("Failed to analyze columns: %v", err)
	}
	nonGeneratedColumns := GetNonGeneratedColumns(columns)
	t.Logf("Non-generated columns: %+v", nonGeneratedColumns)
}
