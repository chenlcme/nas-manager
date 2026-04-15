package repository

import (
	"os"
	"testing"

	"nas-manager/internal/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestSettingRepository_GetSetSetting(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	db, err := gorm.Open(sqlite.Open(tmpFile.Name()), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Migrate the Setting model
	db.AutoMigrate(&model.Setting{})

	repo := NewSettingRepository(db)

	// Test SetSetting and GetSetting
	err = repo.SetSetting("test_key", "test_value")
	if err != nil {
		t.Fatalf("Failed to set setting: %v", err)
	}

	value, err := repo.GetSetting("test_key")
	if err != nil {
		t.Fatalf("Failed to get setting: %v", err)
	}

	if value != "test_value" {
		t.Errorf("Expected 'test_value', got '%s'", value)
	}
}

func TestSettingRepository_GetNonExistentSetting(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	db, err := gorm.Open(sqlite.Open(tmpFile.Name()), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	db.AutoMigrate(&model.Setting{})
	repo := NewSettingRepository(db)

	_, err = repo.GetSetting("non_existent")
	if err == nil {
		t.Error("Expected error for non-existent key, got nil")
	}
}

func TestSettingRepository_HasSettings(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	db, err := gorm.Open(sqlite.Open(tmpFile.Name()), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	db.AutoMigrate(&model.Setting{})
	repo := NewSettingRepository(db)

	// Initially no settings
	has, err := repo.HasSettings()
	if err != nil {
		t.Fatalf("Failed to check has settings: %v", err)
	}
	if has {
		t.Error("Expected false for empty database")
	}

	// Add a setting
	repo.SetSetting("key", "value")

	has, err = repo.HasSettings()
	if err != nil {
		t.Fatalf("Failed to check has settings: %v", err)
	}
	if !has {
		t.Error("Expected true after adding setting")
	}
}

func TestSettingRepository_GetAllSettings(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	db, err := gorm.Open(sqlite.Open(tmpFile.Name()), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	db.AutoMigrate(&model.Setting{})
	repo := NewSettingRepository(db)

	repo.SetSetting("key1", "value1")
	repo.SetSetting("key2", "value2")

	all, err := repo.GetAllSettings()
	if err != nil {
		t.Fatalf("Failed to get all settings: %v", err)
	}

	if len(all) != 2 {
		t.Errorf("Expected 2 settings, got %d", len(all))
	}
	if all["key1"] != "value1" || all["key2"] != "value2" {
		t.Error("Settings mismatch")
	}
}
