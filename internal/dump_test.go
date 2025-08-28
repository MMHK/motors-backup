package internal

import (
	"database/sql"
	"motors-backup/internal/config"
	"os"
	"testing"
)

func TestDumpCreateDatabase(t *testing.T) {
	cfg := config.LoadTestConfig()
	// Skip test if not running in test environment
	if cfg.DBHost == "" {
		t.Skip("Skipping integration test: DB_HOST not set")
	}

	err := StartExport(cfg, func(database *sql.DB, info *MySQLInfo) error {
		err := DumpCreateDatabase(cfg, database, true)
		if err != nil {
			t.Errorf("DumpCreateDatabase failed: %v", err)
		}

		return err
	})

	if err != nil {
		t.Errorf("StartExport failed: %v", err)
	}
}

func TestDumpTable(t *testing.T) {
	cfg := config.LoadTestConfig()

	// Skip test if not running in test environment
	if cfg.DBHost == "" {
		t.Skip("Skipping integration test: DB_HOST not set")
	}

	testTableName := os.Getenv("TEST_TABLE")
	if testTableName == "" {
		t.Skip("Skipping integration test: TEST_TABLE not set")
	}

	err := StartExport(cfg, func(database *sql.DB, info *MySQLInfo) error {
		err := DumpTable(cfg, database, testTableName, "")
		if err != nil {
			t.Errorf("DumpTable failed: %v", err)
		}

		return err
	})

	if err != nil {
		t.Errorf("StartExport failed: %v", err)
	}
}

func TestDumpTableStructure(t *testing.T) {
	cfg := config.LoadTestConfig()

	testTableName := os.Getenv("TEST_TABLE")
	if testTableName == "" {
		t.Skip("Skipping integration test: TEST_TABLE not set")
	}

	err := StartExport(cfg, func(database *sql.DB, info *MySQLInfo) error {
		err := DumpTableStructure(cfg, database, testTableName)
		if err != nil {
			t.Errorf("DumpTableStructure failed: %v", err)
		}

		return err
	})

	if err != nil {
		t.Errorf("StartExport failed: %v", err)
	}
}

func TestDumpTableWithInvalidTable(t *testing.T) {
	cfg := config.LoadTestConfig()
	// Skip test if not running in test environment
	if cfg.DBHost == "" {
		t.Skip("Skipping integration test: DB_HOST not set")
	}

	err := StartExport(cfg, func(database *sql.DB, info *MySQLInfo) error {
		err := DumpTable(cfg, database, "non_existent_table", "")
		if err != nil {
			t.Errorf("DumpTable failed: %v", err)
		}

		return err
	})

	if err == nil {
		t.Error("Expected error for non-existent table, got nil")
	}
}

func TestPrintEnvironmentSettings(t *testing.T) {
	cfg := config.LoadTestConfig()
	err := StartExport(cfg, func(database *sql.DB, info *MySQLInfo) error {

		PrintEnvironmentSettings(cfg, info)
		PrintRestoreConnectionSettings()

		return nil
	})

	if err != nil {
		t.Errorf("StartExport failed: %v", err)
	}
}
