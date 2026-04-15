package service

import (
	"testing"

	"nas-manager/internal/model"
	"nas-manager/internal/repository"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func createTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	db.AutoMigrate(&model.Setting{})

	return db
}

func TestEncryptService_SetupPassword(t *testing.T) {
	db := createTestDB(t)
	repo := repository.NewSettingRepository(db)
	svc := NewEncryptService(repo)

	// Setup password
	err := svc.SetupPassword(&SetupPasswordRequest{Password: "testpassword123"})
	if err != nil {
		t.Fatalf("Failed to setup password: %v", err)
	}

	// Verify password
	valid, err := svc.VerifyPassword(&VerifyPasswordRequest{Password: "testpassword123"})
	if err != nil {
		t.Fatalf("Failed to verify password: %v", err)
	}
	if !valid {
		t.Error("Expected password to be valid")
	}

	// Wrong password should not verify
	valid, err = svc.VerifyPassword(&VerifyPasswordRequest{Password: "wrongpassword"})
	if err != nil {
		t.Fatalf("Failed to verify wrong password: %v", err)
	}
	if valid {
		t.Error("Expected wrong password to be invalid")
	}
}

func TestEncryptService_SetupPassword_TooShort(t *testing.T) {
	db := createTestDB(t)
	repo := repository.NewSettingRepository(db)
	svc := NewEncryptService(repo)

	err := svc.SetupPassword(&SetupPasswordRequest{Password: "short"})
	if err == nil {
		t.Error("Expected error for short password")
	}
}

func TestEncryptService_SetupPassword_AlreadySet(t *testing.T) {
	db := createTestDB(t)
	repo := repository.NewSettingRepository(db)
	svc := NewEncryptService(repo)

	// Setup password first time
	err := svc.SetupPassword(&SetupPasswordRequest{Password: "testpassword123"})
	if err != nil {
		t.Fatalf("Failed to setup password first time: %v", err)
	}

	// Try to setup again
	err = svc.SetupPassword(&SetupPasswordRequest{Password: "anotherpassword"})
	if err == nil {
		t.Error("Expected error when password already set")
	}
}

func TestEncryptService_ChangePassword(t *testing.T) {
	db := createTestDB(t)
	repo := repository.NewSettingRepository(db)
	svc := NewEncryptService(repo)

	// Setup initial password
	err := svc.SetupPassword(&SetupPasswordRequest{Password: "oldpassword123"})
	if err != nil {
		t.Fatalf("Failed to setup initial password: %v", err)
	}

	// Change password
	err = svc.ChangePassword(&ChangePasswordRequest{
		OldPassword: "oldpassword123",
		NewPassword: "newpassword456",
	})
	if err != nil {
		t.Fatalf("Failed to change password: %v", err)
	}

	// Old password should no longer work
	valid, _ := svc.VerifyPassword(&VerifyPasswordRequest{Password: "oldpassword123"})
	if valid {
		t.Error("Expected old password to be invalid after change")
	}

	// New password should work
	valid, _ = svc.VerifyPassword(&VerifyPasswordRequest{Password: "newpassword456"})
	if !valid {
		t.Error("Expected new password to be valid")
	}
}

func TestEncryptService_ChangePassword_WrongOldPassword(t *testing.T) {
	db := createTestDB(t)
	repo := repository.NewSettingRepository(db)
	svc := NewEncryptService(repo)

	// Setup password
	err := svc.SetupPassword(&SetupPasswordRequest{Password: "oldpassword123"})
	if err != nil {
		t.Fatalf("Failed to setup password: %v", err)
	}

	// Try to change with wrong old password
	err = svc.ChangePassword(&ChangePasswordRequest{
		OldPassword: "wrongpassword",
		NewPassword: "newpassword456",
	})
	if err == nil {
		t.Error("Expected error when old password is wrong")
	}
}

func TestEncryptService_HasPassword(t *testing.T) {
	db := createTestDB(t)
	repo := repository.NewSettingRepository(db)
	svc := NewEncryptService(repo)

	// Initially no password
	has, _ := svc.HasPassword()
	if has {
		t.Error("Expected no password initially")
	}

	// Set password
	svc.SetupPassword(&SetupPasswordRequest{Password: "testpassword123"})

	// Now has password
	has, _ = svc.HasPassword()
	if !has {
		t.Error("Expected password to be set")
	}
}
