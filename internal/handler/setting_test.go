package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"nas-manager/internal/model"
	"nas-manager/internal/repository"
	"nas-manager/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func createSettingTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	db.AutoMigrate(&model.Setting{})

	return db
}

func setupSettingRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	settingRepo := repository.NewSettingRepository(db)
	settingSvc := service.NewSettingService(settingRepo)
	settingHandler := NewSettingHandler(settingSvc)

	r.GET("/api/setup/status", settingHandler.GetSetupStatus)
	r.POST("/api/setup", settingHandler.SaveSetup)

	return r
}

// TestSetupHandler_GetSetupStatus_NeedsSetup tests initial setup status when no music_dir is set
func TestSetupHandler_GetSetupStatus_NeedsSetup(t *testing.T) {
	db := createSettingTestDB(t)
	router := setupSettingRouter(db)

	req, _ := http.NewRequest("GET", "/api/setup/status", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	data, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected data to be a map")
	}

	if data["needs_setup"] != true {
		t.Errorf("Expected needs_setup to be true, got %v", data["needs_setup"])
	}
}

// TestSetupHandler_GetSetupStatus_AlreadyConfigured tests setup status when music_dir is already set
func TestSetupHandler_GetSetupStatus_AlreadyConfigured(t *testing.T) {
	db := createSettingTestDB(t)

	// Pre-configure music directory
	settingRepo := repository.NewSettingRepository(db)
	settingRepo.SetMusicDir("/test/music")

	router := setupSettingRouter(db)

	req, _ := http.NewRequest("GET", "/api/setup/status", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data := response["data"].(map[string]interface{})
	if data["needs_setup"] != false {
		t.Errorf("Expected needs_setup to be false, got %v", data["needs_setup"])
	}
	if data["music_dir"] != "/test/music" {
		t.Errorf("Expected music_dir to be '/test/music', got %v", data["music_dir"])
	}
}

// TestSetupHandler_SaveSetup_ValidConfig tests saving a valid configuration
func TestSetupHandler_SaveSetup_ValidConfig(t *testing.T) {
	db := createSettingTestDB(t)

	// Create a real temp directory for testing
	tmpDir, err := os.MkdirTemp("", "music-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	router := setupSettingRouter(db)

	reqBody := `{"music_dir": "` + tmpDir + `"}`
	req, _ := http.NewRequest("POST", "/api/setup", nil)
	req.Body = createRequestBody(reqBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["error"] != nil {
		t.Errorf("Expected no error, got %v", response["error"])
	}
}

// TestSetupHandler_SaveSetup_EmptyMusicDir tests saving with empty music_dir
func TestSetupHandler_SaveSetup_EmptyMusicDir(t *testing.T) {
	db := createSettingTestDB(t)
	router := setupSettingRouter(db)

	reqBody := `{"music_dir": ""}`
	req, _ := http.NewRequest("POST", "/api/setup", nil)
	req.Body = createRequestBody(reqBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// TestSetupHandler_SaveSetup_NonExistentDir tests saving with non-existent directory
func TestSetupHandler_SaveSetup_NonExistentDir(t *testing.T) {
	db := createSettingTestDB(t)
	router := setupSettingRouter(db)

	reqBody := `{"music_dir": "/non/existent/directory"}`
	req, _ := http.NewRequest("POST", "/api/setup", nil)
	req.Body = createRequestBody(reqBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// TestSetupHandler_SaveSetup_InvalidJSON tests saving with invalid JSON
func TestSetupHandler_SaveSetup_InvalidJSON(t *testing.T) {
	db := createSettingTestDB(t)
	router := setupSettingRouter(db)

	reqBody := `{invalid json}`
	req, _ := http.NewRequest("POST", "/api/setup", nil)
	req.Body = createRequestBody(reqBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// TestSetupHandler_SaveSetup_FileInsteadOfDirectory tests saving with a file path instead of directory
func TestSetupHandler_SaveSetup_FileInsteadOfDirectory(t *testing.T) {
	db := createSettingTestDB(t)

	// Create a temp file (not directory)
	tmpFile, err := os.CreateTemp("", "music-test-")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	router := setupSettingRouter(db)

	reqBody := `{"music_dir": "` + tmpFile.Name() + `"}`
	req, _ := http.NewRequest("POST", "/api/setup", nil)
	req.Body = createRequestBody(reqBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// TestSetupHandler_SaveSetup_WithDBPath tests saving configuration with custom DB path
func TestSetupHandler_SaveSetup_WithDBPath(t *testing.T) {
	db := createSettingTestDB(t)

	// Create temp directories
	tmpMusicDir, err := os.MkdirTemp("", "music-test-")
	if err != nil {
		t.Fatalf("Failed to create temp music dir: %v", err)
	}
	defer os.RemoveAll(tmpMusicDir)

	tmpDBDir, err := os.MkdirTemp("", "db-test-")
	if err != nil {
		t.Fatalf("Failed to create temp db dir: %v", err)
	}
	defer os.RemoveAll(tmpDBDir)

	router := setupSettingRouter(db)

	// Test with writable db path
	tmpDBFile := filepath.Join(tmpDBDir, "test.db")
	reqBody := `{"music_dir": "` + tmpMusicDir + `", "db_path": "` + tmpDBFile + `"}`
	req, _ := http.NewRequest("POST", "/api/setup", nil)
	req.Body = createRequestBody(reqBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}