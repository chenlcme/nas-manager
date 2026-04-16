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

func createScanTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	db.AutoMigrate(&model.Song{}, &model.Setting{})

	return db
}

func setupScanRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	songRepo := repository.NewSongRepository(db)
	settingRepo := repository.NewSettingRepository(db)
	id3Svc := service.NewID3Service(songRepo)
	scannerSvc := service.NewScannerService(id3Svc, songRepo)
	scanHandler := NewScanHandler(scannerSvc, songRepo, settingRepo)

	r.POST("/api/songs/scan", scanHandler.Scan)
	r.POST("/api/songs/cleanup", scanHandler.Cleanup)

	return r
}

// TestScanHandler_Scan_NoMusicDir tests scan when no music directory is configured
func TestScanHandler_Scan_NoMusicDir(t *testing.T) {
	db := createScanTestDB(t)
	router := setupScanRouter(db)

	req, _ := http.NewRequest("POST", "/api/songs/scan", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data := response["error"].(map[string]interface{})
	if data["code"] != "NO_MUSIC_DIR" {
		t.Errorf("Expected error code 'NO_MUSIC_DIR', got '%s'", data["code"])
	}
}

// TestScanHandler_Scan_DirNotExist tests scan when configured music directory doesn't exist
func TestScanHandler_Scan_DirNotExist(t *testing.T) {
	db := createScanTestDB(t)

	// Set a non-existent music directory
	settingRepo := repository.NewSettingRepository(db)
	settingRepo.SetMusicDir("/non/existent/directory")

	router := setupScanRouter(db)

	req, _ := http.NewRequest("POST", "/api/songs/scan", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data := response["error"].(map[string]interface{})
	if data["code"] != "DIR_NOT_EXIST" {
		t.Errorf("Expected error code 'DIR_NOT_EXIST', got '%s'", data["code"])
	}
}

// TestScanHandler_Scan_ValidDir tests successful scan with valid music directory
func TestScanHandler_Scan_ValidDir(t *testing.T) {
	db := createScanTestDB(t)

	// Create temp music directory with some files
	tmpDir, err := os.MkdirTemp("", "music-scan-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create some music files
	os.WriteFile(filepath.Join(tmpDir, "song1.mp3"), []byte("fake mp3 content"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "song2.flac"), []byte("fake flac content"), 0644)

	// Set music directory
	settingRepo := repository.NewSettingRepository(db)
	settingRepo.SetMusicDir(tmpDir)

	router := setupScanRouter(db)

	reqBody := `{"mode": "incremental"}`
	req, _ := http.NewRequest("POST", "/api/songs/scan", nil)
	req.Body = createRequestBody(reqBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["error"] != nil {
		t.Errorf("Expected no error, got %v", response["error"])
	}
}

// TestScanHandler_Scan_FullMode tests scan with full mode
func TestScanHandler_Scan_FullMode(t *testing.T) {
	db := createScanTestDB(t)

	// Create temp music directory
	tmpDir, err := os.MkdirTemp("", "music-scan-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a music file
	os.WriteFile(filepath.Join(tmpDir, "song1.mp3"), []byte("fake mp3 content"), 0644)

	// Set music directory
	settingRepo := repository.NewSettingRepository(db)
	settingRepo.SetMusicDir(tmpDir)

	router := setupScanRouter(db)

	reqBody := `{"mode": "full"}`
	req, _ := http.NewRequest("POST", "/api/songs/scan", nil)
	req.Body = createRequestBody(reqBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

// TestScanHandler_Scan_DefaultMode tests scan with default (incremental) mode
func TestScanHandler_Scan_DefaultMode(t *testing.T) {
	db := createScanTestDB(t)

	// Create temp music directory
	tmpDir, err := os.MkdirTemp("", "music-scan-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a music file
	os.WriteFile(filepath.Join(tmpDir, "song1.mp3"), []byte("fake mp3 content"), 0644)

	// Set music directory
	settingRepo := repository.NewSettingRepository(db)
	settingRepo.SetMusicDir(tmpDir)

	router := setupScanRouter(db)

	// Empty body should default to incremental
	req, _ := http.NewRequest("POST", "/api/songs/scan", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

// TestScanHandler_Cleanup_Success tests successful cleanup of orphan records
func TestScanHandler_Cleanup_Success(t *testing.T) {
	db := createScanTestDB(t)

	// Create some songs in the database
	songs := []*model.Song{
		{FilePath: "/music/song1.mp3", Title: "Song 1"},
		{FilePath: "/music/song2.mp3", Title: "Song 2"},
	}
	for _, s := range songs {
		db.Create(s)
	}

	router := setupScanRouter(db)

	req, _ := http.NewRequest("POST", "/api/songs/cleanup", nil)
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

// TestScanHandler_Cleanup_EmptyDatabase tests cleanup with empty database
func TestScanHandler_Cleanup_EmptyDatabase(t *testing.T) {
	db := createScanTestDB(t)
	router := setupScanRouter(db)

	req, _ := http.NewRequest("POST", "/api/songs/cleanup", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

// TestScanHandler_Scan_SubdirectoryFiles tests scanning with nested subdirectories
func TestScanHandler_Scan_SubdirectoryFiles(t *testing.T) {
	db := createScanTestDB(t)

	// Create temp music directory with subdirectories
	tmpDir, err := os.MkdirTemp("", "music-scan-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create subdirectories and files
	os.MkdirAll(filepath.Join(tmpDir, "rock"), 0755)
	os.MkdirAll(filepath.Join(tmpDir, "pop"), 0755)
	os.WriteFile(filepath.Join(tmpDir, "rock", "song1.mp3"), []byte("fake mp3"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "pop", "song2.mp3"), []byte("fake mp3"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "song3.mp3"), []byte("fake mp3"), 0644)

	// Set music directory
	settingRepo := repository.NewSettingRepository(db)
	settingRepo.SetMusicDir(tmpDir)

	router := setupScanRouter(db)

	reqBody := `{"mode": "full"}`
	req, _ := http.NewRequest("POST", "/api/songs/scan", nil)
	req.Body = createRequestBody(reqBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
	}
}