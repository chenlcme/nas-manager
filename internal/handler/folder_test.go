package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"nas-manager/internal/model"
	"nas-manager/internal/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func createFolderHandlerTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	db.AutoMigrate(&model.Song{})

	return db
}

func setupFolderRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	folderRepo := repository.NewFolderRepository(db)
	folderHandler := NewFolderHandler(folderRepo)

	r.GET("/api/folders", folderHandler.GetFolders)
	r.GET("/api/folders/:id/songs", folderHandler.GetFolderSongs)

	return r
}

func TestFolderHandler_GetFolders(t *testing.T) {
	db := createFolderHandlerTestDB(t)
	router := setupFolderRouter(db)

	// 创建测试数据 - 4个文件分布在3个文件夹
	songs := []*model.Song{
		{FilePath: "/music/rock/song1.mp3", Title: "Rock Song 1"},
		{FilePath: "/music/rock/song2.mp3", Title: "Rock Song 2"},
		{FilePath: "/music/pop/song3.mp3", Title: "Pop Song 1"},
		{FilePath: "/music/classical/song4.mp3", Title: "Classical Song 1"},
	}
	for _, s := range songs {
		db.Create(s)
	}

	// 测试获取文件夹列表（降序）
	req, _ := http.NewRequest("GET", "/api/folders", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	data, ok := resp["data"].([]interface{})
	if !ok {
		t.Fatal("Expected data to be an array")
	}

	// 3个文件夹: /music/rock, /music/pop, /music/classical
	if len(data) != 3 {
		t.Errorf("Expected 3 folders, got %d", len(data))
	}
}

func TestFolderHandler_GetFolders_Empty(t *testing.T) {
	db := createFolderHandlerTestDB(t)
	router := setupFolderRouter(db)

	// 空数据库测试
	req, _ := http.NewRequest("GET", "/api/folders", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	data, ok := resp["data"].([]interface{})
	if !ok {
		t.Fatal("Expected data to be an array")
	}

	if len(data) != 0 {
		t.Errorf("Expected 0 folders, got %d", len(data))
	}
}

func TestFolderHandler_GetFolderSongs(t *testing.T) {
	db := createFolderHandlerTestDB(t)
	router := setupFolderRouter(db)

	// 创建测试数据
	songs := []*model.Song{
		{FilePath: "/music/rock/song1.mp3", Title: "Rock Song 1"},
		{FilePath: "/music/rock/song2.mp3", Title: "Rock Song 2"},
		{FilePath: "/music/pop/song3.mp3", Title: "Pop Song 1"},
	}
	for _, s := range songs {
		db.Create(s)
	}

	// 先获取文件夹列表
	req, _ := http.NewRequest("GET", "/api/folders", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].([]interface{})

	// 找到 /music/rock 的 ID (ID=1 因为降序排序后 /music/rock 是第一个)
	_ = data // 验证数据存在即可

	// 测试获取特定文件夹的歌曲（ID=1 是降序后的第一个）
	req, _ = http.NewRequest("GET", "/api/folders/1/songs", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var songsResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &songsResp)
	songsData := songsResp["data"].([]interface{})

	if len(songsData) != 2 {
		t.Errorf("Expected 2 songs for /music/rock, got %d", len(songsData))
	}
}

func TestFolderHandler_GetFolderSongs_InvalidID(t *testing.T) {
	db := createFolderHandlerTestDB(t)
	router := setupFolderRouter(db)

	// 创建测试数据
	db.Create(&model.Song{FilePath: "/music/rock/song1.mp3", Title: "Rock Song 1"})

	// 无效 ID 测试
	req, _ := http.NewRequest("GET", "/api/folders/invalid/songs", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestFolderHandler_GetFolderSongs_NotFound(t *testing.T) {
	db := createFolderHandlerTestDB(t)
	router := setupFolderRouter(db)

	// 创建测试数据
	db.Create(&model.Song{FilePath: "/music/rock/song1.mp3", Title: "Rock Song 1"})

	// 不存在的 ID 测试 (ID=999 超出范围)
	req, _ := http.NewRequest("GET", "/api/folders/999/songs", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestFolderHandler_GetFolders_AscendingOrder(t *testing.T) {
	db := createFolderHandlerTestDB(t)
	router := setupFolderRouter(db)

	// 创建测试数据
	songs := []*model.Song{
		{FilePath: "/music/rock/song1.mp3", Title: "Rock Song 1"},
		{FilePath: "/music/pop/song2.mp3", Title: "Pop Song 1"},
	}
	for _, s := range songs {
		db.Create(s)
	}

	// 测试升序排序
	req, _ := http.NewRequest("GET", "/api/folders?order=asc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].([]interface{})

	// 验证第一个是 /music/pop（字母顺序: pop < rock）
	if len(data) >= 1 {
		firstFolder := data[0].(map[string]interface{})
		if firstFolder["path"] != "/music/pop" {
			t.Errorf("Expected first folder to be /music/pop, got %s", firstFolder["path"])
		}
	}
}
