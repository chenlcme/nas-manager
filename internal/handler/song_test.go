package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"nas-manager/internal/model"
	"nas-manager/internal/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// API response types for testing
type APIResponse struct {
	Data  interface{} `json:"data"`
	Error *APIError  `json:"error,omitempty"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto-migrate the schema
	if err := db.AutoMigrate(&model.Song{}); err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func TestGetSong_Success(t *testing.T) {
	db := setupTestDB(t)

	// Create a test song
	testSong := &model.Song{
		FilePath:  "/test/song.mp3",
		Title:    "Test Song",
		Artist:   "Test Artist",
		Album:    "Test Album",
		Year:     2024,
		Genre:    "Pop",
		TrackNum: 1,
		Duration: 180,
		FileSize: 1234567,
	}
	if err := db.Create(testSong).Error; err != nil {
		t.Fatalf("Failed to create test song: %v", err)
	}

	// Setup handler
	songRepo := repository.NewSongRepository(db)
	handler := NewSongHandler(songRepo)

	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/songs/:id", handler.GetSong)

	// Make request
	req, _ := http.NewRequest("GET", "/songs/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Debug: print raw response
	t.Logf("Response body: %s", w.Body.String())

	// Parse response into map structure
	var resp map[string]json.RawMessage
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Get the data field
	dataRaw, ok := resp["data"]
	if !ok {
		t.Fatal("Expected 'data' field in response")
	}

	// Unmarshal data into song struct
	var song model.Song
	if err := json.Unmarshal(dataRaw, &song); err != nil {
		t.Fatalf("Failed to unmarshal song data: %v", err)
	}

	if song.Title != "Test Song" {
		t.Errorf("Expected title 'Test Song', got '%s'", song.Title)
	}
	if song.Artist != "Test Artist" {
		t.Errorf("Expected artist 'Test Artist', got '%s'", song.Artist)
	}
	if song.Album != "Test Album" {
		t.Errorf("Expected album 'Test Album', got '%s'", song.Album)
	}
	if song.Year != 2024 {
		t.Errorf("Expected year 2024, got %d", song.Year)
	}
	if song.Genre != "Pop" {
		t.Errorf("Expected genre 'Pop', got '%s'", song.Genre)
	}
}

func TestGetSong_NotFound(t *testing.T) {
	db := setupTestDB(t)

	// Setup handler
	songRepo := repository.NewSongRepository(db)
	handler := NewSongHandler(songRepo)

	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/songs/:id", handler.GetSong)

	// Make request for non-existent song
	req, _ := http.NewRequest("GET", "/songs/999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check status code
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}

	// Verify error response format
	var resp APIResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.Error == nil {
		t.Fatal("Expected error response")
	}
	if resp.Error.Code != "SONG_NOT_FOUND" {
		t.Errorf("Expected error code 'SONG_NOT_FOUND', got '%s'", resp.Error.Code)
	}
}

func TestGetSong_InvalidID(t *testing.T) {
	db := setupTestDB(t)

	// Setup handler
	songRepo := repository.NewSongRepository(db)
	handler := NewSongHandler(songRepo)

	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/songs/:id", handler.GetSong)

	// Make request with invalid ID
	req, _ := http.NewRequest("GET", "/songs/invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check status code
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// TestGetSong_DBError tests handling of DB errors other than not-found
func TestGetSong_DBError(t *testing.T) {
	db := setupTestDB(t)

	// Create and save a song first
	testSong := &model.Song{
		FilePath: "/test/song.mp3",
		Title:    "Test Song",
	}
	if err := db.Create(testSong).Error; err != nil {
		t.Fatalf("Failed to create test song: %v", err)
	}

	songRepo := repository.NewSongRepository(db)
	handler := NewSongHandler(songRepo)

	// Close the underlying DB to trigger an error on subsequent calls
	dbSQL, err := db.DB()
	if err != nil {
		t.Fatalf("Failed to get underlying DB: %v", err)
	}
	dbSQL.Close()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/songs/:id", handler.GetSong)

	req, _ := http.NewRequest("GET", "/songs/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should return 500 Internal Server Error for DB errors
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d for DB error, got %d", http.StatusInternalServerError, w.Code)
	}
}

// TestGetSong_NullZeroFields tests song with null/zero fields
func TestGetSong_NullZeroFields(t *testing.T) {
	db := setupTestDB(t)

	// Create a song with minimal/zero fields
	testSong := &model.Song{
		FilePath: "/test/song.mp3",
		Title:    "",
		Artist:   "",
		Album:    "",
		Year:     0,
		Genre:    "",
		TrackNum: 0,
		Duration: 0,
		FileSize: 0,
	}
	if err := db.Create(testSong).Error; err != nil {
		t.Fatalf("Failed to create test song: %v", err)
	}

	songRepo := repository.NewSongRepository(db)
	handler := NewSongHandler(songRepo)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/songs/:id", handler.GetSong)

	req, _ := http.NewRequest("GET", "/songs/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Parse response
	var resp map[string]json.RawMessage
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	dataRaw, ok := resp["data"]
	if !ok {
		t.Fatal("Expected 'data' field in response")
	}

	var song model.Song
	if err := json.Unmarshal(dataRaw, &song); err != nil {
		t.Fatalf("Failed to unmarshal song data: %v", err)
	}

	// Zero values should be returned as empty/zero
	if song.Title != "" {
		t.Errorf("Expected empty title, got '%s'", song.Title)
	}
	if song.Year != 0 {
		t.Errorf("Expected year 0, got %d", song.Year)
	}
	if song.Duration != 0 {
		t.Errorf("Expected duration 0, got %d", song.Duration)
	}
}

func TestDeleteSongs_Success(t *testing.T) {
	db := setupTestDB(t)

	// Create test songs with actual files
	tmpDir := t.TempDir()
	song1Path := tmpDir + "/song1.mp3"
	song2Path := tmpDir + "/song2.mp3"

	// Create actual files
	if err := createTestFile(song1Path); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	if err := createTestFile(song2Path); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	testSong1 := &model.Song{
		FilePath: song1Path,
		Title:    "Test Song 1",
		Artist:   "Test Artist",
	}
	testSong2 := &model.Song{
		FilePath: song2Path,
		Title:    "Test Song 2",
		Artist:   "Test Artist",
	}
	if err := db.Create(testSong1).Error; err != nil {
		t.Fatalf("Failed to create test song 1: %v", err)
	}
	if err := db.Create(testSong2).Error; err != nil {
		t.Fatalf("Failed to create test song 2: %v", err)
	}

	// Setup handler
	songRepo := repository.NewSongRepository(db)
	handler := NewSongHandler(songRepo)

	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/songs/delete", handler.DeleteSongs)

	// Make request
	reqBody := `{"ids": [1, 2]}`
	req, _ := http.NewRequest("POST", "/songs/delete", nil)
	req.Body = createRequestBody(reqBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Parse response
	var resp map[string]json.RawMessage
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	dataRaw, ok := resp["data"]
	if !ok {
		t.Fatal("Expected 'data' field in response")
	}

	var result DeleteResult
	if err := json.Unmarshal(dataRaw, &result); err != nil {
		t.Fatalf("Failed to unmarshal result: %v", err)
	}

	if result.Total != 2 {
		t.Errorf("Expected total 2, got %d", result.Total)
	}
	if result.Succeeded != 2 {
		t.Errorf("Expected succeeded 2, got %d", result.Succeeded)
	}
	if result.Failed != 0 {
		t.Errorf("Expected failed 0, got %d", result.Failed)
	}

	// Verify database records are deleted
	var count int64
	db.Model(&model.Song{}).Count(&count)
	if count != 0 {
		t.Errorf("Expected 0 songs in DB, got %d", count)
	}

	// Verify files are deleted
	if fileExists(song1Path) {
		t.Errorf("Expected file %s to be deleted", song1Path)
	}
	if fileExists(song2Path) {
		t.Errorf("Expected file %s to be deleted", song2Path)
	}
}

func TestDeleteSongs_EmptyIDs(t *testing.T) {
	db := setupTestDB(t)

	// Setup handler
	songRepo := repository.NewSongRepository(db)
	handler := NewSongHandler(songRepo)

	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/songs/delete", handler.DeleteSongs)

	// Make request with empty IDs
	reqBody := `{"ids": []}`
	req, _ := http.NewRequest("POST", "/songs/delete", nil)
	req.Body = createRequestBody(reqBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should return 400 Bad Request
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestDeleteSongs_PartialFailure(t *testing.T) {
	db := setupTestDB(t)

	// Create one test song with an actual file
	tmpDir := t.TempDir()
	song1Path := tmpDir + "/song1.mp3"
	if err := createTestFile(song1Path); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	testSong1 := &model.Song{
		FilePath: song1Path,
		Title:    "Test Song 1",
	}
	if err := db.Create(testSong1).Error; err != nil {
		t.Fatalf("Failed to create test song 1: %v", err)
	}

	// Setup handler
	songRepo := repository.NewSongRepository(db)
	handler := NewSongHandler(songRepo)

	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/songs/delete", handler.DeleteSongs)

	// Make request with one valid ID and one non-existent ID
	reqBody := `{"ids": [1, 999]}`
	req, _ := http.NewRequest("POST", "/songs/delete", nil)
	req.Body = createRequestBody(reqBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should still return 200 OK
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Parse response
	var resp map[string]json.RawMessage
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	dataRaw, ok := resp["data"]
	if !ok {
		t.Fatal("Expected 'data' field in response")
	}

	var result DeleteResult
	if err := json.Unmarshal(dataRaw, &result); err != nil {
		t.Fatalf("Failed to unmarshal result: %v", err)
	}

	if result.Total != 2 {
		t.Errorf("Expected total 2, got %d", result.Total)
	}
	if result.Succeeded != 1 {
		t.Errorf("Expected succeeded 1, got %d", result.Succeeded)
	}
	if result.Failed != 1 {
		t.Errorf("Expected failed 1, got %d", result.Failed)
	}
}

// Helper function to create test file
func createTestFile(path string) error {
	content := []byte("fake mp3 content")
	return os.WriteFile(path, content, 0644)
}

func createRequestBody(content string) *requestBody {
	return &requestBody{content: content}
}

type requestBody struct {
	content string
	pos     int
}

func (rb *requestBody) Read(p []byte) (int, error) {
	if rb.pos >= len(rb.content) {
		return 0, nil
	}
	n := copy(p, rb.content[rb.pos:])
	rb.pos += n
	return n, nil
}

func (rb *requestBody) Close() error {
	return nil
}

// fileExists is a simple check for file existence
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

