package service

import (
	"os"
	"testing"

	"nas-manager/internal/model"
	"nas-manager/internal/repository"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestSettingService_CheckSetupRequired(t *testing.T) {
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

	repo := repository.NewSettingRepository(db)
	svc := NewSettingService(repo)

	// Initially needs setup
	status, err := svc.CheckSetupRequired()
	if err != nil {
		t.Fatalf("Failed to check setup required: %v", err)
	}
	if !status.NeedsSetup {
		t.Error("Expected needs setup to be true initially")
	}

	// After setting music_dir, should not need setup
	repo.SetMusicDir("/test/music")
	status, err = svc.CheckSetupRequired()
	if err != nil {
		t.Fatalf("Failed to check setup required: %v", err)
	}
	if status.NeedsSetup {
		t.Error("Expected needs setup to be false after setting music_dir")
	}
}

func TestSettingService_SaveSetupConfig_InvalidMusicDir(t *testing.T) {
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

	repo := repository.NewSettingRepository(db)
	svc := NewSettingService(repo)

	// Empty music dir should fail
	err = svc.SaveSetupConfig(&SetupConfig{
		MusicDir: "",
	})
	if err == nil {
		t.Error("Expected error for empty music_dir")
	}
}

func TestSettingService_SaveSetupConfig_NonExistentDir(t *testing.T) {
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

	repo := repository.NewSettingRepository(db)
	svc := NewSettingService(repo)

	// Non-existent directory should fail
	err = svc.SaveSetupConfig(&SetupConfig{
		MusicDir: "/non/existent/directory",
	})
	if err == nil {
		t.Error("Expected error for non-existent directory")
	}
}

func TestSettingService_SaveSetupConfig_ValidDir(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	// Create a temp directory for music
	tmpMusicDir, err := os.MkdirTemp("", "music-")
	if err != nil {
		t.Fatalf("Failed to create temp music dir: %v", err)
	}
	defer os.RemoveAll(tmpMusicDir)

	db, err := gorm.Open(sqlite.Open(tmpFile.Name()), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	db.AutoMigrate(&model.Setting{})

	repo := repository.NewSettingRepository(db)
	svc := NewSettingService(repo)

	// Valid directory should succeed
	err = svc.SaveSetupConfig(&SetupConfig{
		MusicDir: tmpMusicDir,
	})
	if err != nil {
		t.Errorf("Unexpected error for valid directory: %v", err)
	}

	// Verify the setting was saved
	savedDir, err := repo.GetMusicDir()
	if err != nil {
		t.Fatalf("Failed to get saved music dir: %v", err)
	}
	if savedDir != tmpMusicDir {
		t.Errorf("Expected '%s', got '%s'", tmpMusicDir, savedDir)
	}
}
