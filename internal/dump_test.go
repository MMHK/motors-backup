package internal

import (
	"motors-backup/internal/config"
	"os"
	"testing"
)

func TestDumpTable(t *testing.T) {
	cfg := config.LoadTestConfig()

	// Skip test if not running in test environment
	if cfg.DBHost == "" {
		t.Skip("Skipping integration test: TEST_DB_HOST not set")
	}

	testTableName := os.Getenv("TEST_TABLE")
	if testTableName == "" {
		t.Skip("Skipping integration test: TEST_TABLE not set")
	}

	err := DumpTable(cfg, testTableName)
	if err != nil {
		t.Errorf("DumpTable failed: %v", err)
	}
}

func TestDumpTableWithInvalidTable(t *testing.T) {
	cfg := config.LoadTestConfig()
	// Skip test if not running in test environment
	if cfg.DBHost == "" {
		t.Skip("Skipping integration test: TEST_DB_HOST not set")
	}

	err := DumpTable(cfg, "non_existent_table")
	if err == nil {
		t.Error("Expected error for non-existent table, got nil")
	}
}
