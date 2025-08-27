package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// 保存原始环境变量
	originalHost := os.Getenv("DB_HOST")
	originalPort := os.Getenv("DB_PORT")
	originalUser := os.Getenv("DB_USER")
	originalPassword := os.Getenv("DB_PASSWORD")
	originalName := os.Getenv("DB_NAME")

	// 确保在测试后恢复环境变量
	defer func() {
		os.Setenv("DB_HOST", originalHost)
		os.Setenv("DB_PORT", originalPort)
		os.Setenv("DB_USER", originalUser)
		os.Setenv("DB_PASSWORD", originalPassword)
		os.Setenv("DB_NAME", originalName)
	}()

	// 设置测试环境变量
	os.Setenv("DB_HOST", "testhost")
	os.Setenv("DB_PORT", "3307")
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_PASSWORD", "testpass")
	os.Setenv("DB_NAME", "testdb")

	// 加载配置
	cfg := LoadConfig()

	// 验证配置值
	if cfg.DBHost != "testhost" {
		t.Errorf("Expected DB_HOST to be 'testhost', got '%s'", cfg.DBHost)
	}

	if cfg.DBPort != 3307 {
		t.Errorf("Expected DB_PORT to be 3307, got %d", cfg.DBPort)
	}

	if cfg.DBUser != "testuser" {
		t.Errorf("Expected DB_USER to be 'testuser', got '%s'", cfg.DBUser)
	}

	if cfg.DBPassword != "testpass" {
		t.Errorf("Expected DB_PASSWORD to be 'testpass', got '%s'", cfg.DBPassword)
	}

	if cfg.DBName != "testdb" {
		t.Errorf("Expected DB_NAME to be 'testdb', got '%s'", cfg.DBName)
	}
}

func TestLoadConfigDefaultValues(t *testing.T) {
	// 保存原始环境变量
	originalHost := os.Getenv("DB_HOST")
	originalPort := os.Getenv("DB_PORT")
	originalUser := os.Getenv("DB_USER")
	originalPassword := os.Getenv("DB_PASSWORD")
	originalName := os.Getenv("DB_NAME")

	// 确保在测试后恢复环境变量
	defer func() {
		os.Setenv("DB_HOST", originalHost)
		os.Setenv("DB_PORT", originalPort)
		os.Setenv("DB_USER", originalUser)
		os.Setenv("DB_PASSWORD", originalPassword)
		os.Setenv("DB_NAME", originalName)
	}()

	// 清除环境变量以测试默认值
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_PASSWORD")
	os.Setenv("DB_NAME", "testdb") // DB_NAME是必需的

	// 加载配置
	cfg := LoadConfig()

	// 验证默认值
	if cfg.DBHost != "localhost" {
		t.Errorf("Expected default DB_HOST to be 'localhost', got '%s'", cfg.DBHost)
	}

	if cfg.DBPort != 3306 {
		t.Errorf("Expected default DB_PORT to be 3306, got %d", cfg.DBPort)
	}

	if cfg.DBUser != "root" {
		t.Errorf("Expected default DB_USER to be 'root', got '%s'", cfg.DBUser)
	}

	if cfg.DBPassword != "" {
		t.Errorf("Expected default DB_PASSWORD to be empty, got '%s'", cfg.DBPassword)
	}

	if cfg.DBName != "testdb" {
		t.Errorf("Expected DB_NAME to be 'testdb', got '%s'", cfg.DBName)
	}
}

func TestLoadConfigInvalidPort(t *testing.T) {
	// 保存原始环境变量
	originalPort := os.Getenv("DB_PORT")
	originalName := os.Getenv("DB_NAME")

	// 确保在测试后恢复环境变量
	defer func() {
		os.Setenv("DB_PORT", originalPort)
		os.Setenv("DB_NAME", originalName)
	}()

	// 设置无效的端口值
	os.Setenv("DB_PORT", "invalid")
	os.Setenv("DB_NAME", "testdb")

	// 加载配置
	cfg := LoadConfig()

	// 验证端口是否回退到默认值
	if cfg.DBPort != 3306 {
		t.Errorf("Expected DB_PORT to fallback to default 3306 when invalid, got %d", cfg.DBPort)
	}
}
