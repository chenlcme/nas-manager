package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"nas-manager/internal/model"
	"nas-manager/internal/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func createAlbumTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	db.AutoMigrate(&model.Song{})

	return db
}

func setupAlbumRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	albumRepo := repository.NewAlbumRepository(db)
	albumHandler := NewAlbumHandler(albumRepo)

	r.GET("/api/albums", albumHandler.GetAlbums)
	r.GET("/api/albums/:id/songs", albumHandler.GetAlbumSongs)

	return r
}

func TestAlbumHandler_GetAlbums(t *testing.T) {
	db := createAlbumTestDB(t)

	// 创建测试数据
	songs := []*model.Song{
		{FilePath: "/test/song1.mp3", Title: "晴天", Album: "叶惠美", Artist: "周杰伦"},
		{FilePath: "/test/song2.mp3", Title: "稻香", Album: "叶惠美", Artist: "周杰伦"},
		{FilePath: "/test/song3.mp3", Title: "江南", Album: "江南", Artist: "林俊杰"},
	}
	for _, s := range songs {
		db.Create(s)
	}

	router := setupAlbumRouter(db)

	// 测试获取专辑列表
	req, _ := http.NewRequest("GET", "/api/albums", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data, ok := response["data"].([]interface{})
	if !ok {
		t.Fatal("Expected data to be an array")
	}

	if len(data) != 2 {
		t.Errorf("Expected 2 albums, got %d", len(data))
	}
}

func TestAlbumHandler_GetAlbums_Empty(t *testing.T) {
	db := createAlbumTestDB(t)
	router := setupAlbumRouter(db)

	req, _ := http.NewRequest("GET", "/api/albums", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	t.Logf("Response body: %s", w.Body.String())

	data := response["data"]
	if data == nil {
		t.Log("Data is nil, which is equivalent to empty array - PASS")
		return
	}

	dataSlice, ok := data.([]interface{})
	if !ok {
		t.Fatalf("Expected data to be an array, got %T", data)
	}

	if len(dataSlice) != 0 {
		t.Errorf("Expected 0 albums, got %d", len(dataSlice))
	}
}

func TestAlbumHandler_GetAlbums_SortAsc(t *testing.T) {
	db := createAlbumTestDB(t)

	// 创建测试数据 - 使用全大写专辑名以便测试排序
	songs := []*model.Song{
		{FilePath: "/test/song1.mp3", Title: "Song 1", Album: "ZEBRA", Artist: "Artist1"},
		{FilePath: "/test/song2.mp3", Title: "Song 2", Album: "APPLE", Artist: "Artist2"},
		{FilePath: "/test/song3.mp3", Title: "Song 3", Album: "MANGO", Artist: "Artist3"},
	}
	for _, s := range songs {
		db.Create(s)
	}

	router := setupAlbumRouter(db)

	// 测试升序排序 - ASCII 顺序: APPLE < MANGO < ZEBRA
	req, _ := http.NewRequest("GET", "/api/albums?order=asc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data := response["data"].([]interface{})
	album0 := data[0].(map[string]interface{})

	// 升序时 APPLE 应该在最前面
	if album0["name"] != "APPLE" {
		t.Errorf("Expected first album to be APPLE (asc), got %v", album0["name"])
	}

	// 降序排序 - ASCII 顺序: ZEBRA > MANGO > APPLE
	req2, _ := http.NewRequest("GET", "/api/albums?order=desc", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	var response2 map[string]interface{}
	json.Unmarshal(w2.Body.Bytes(), &response2)

	data2 := response2["data"].([]interface{})
	album02 := data2[0].(map[string]interface{})

	// 降序时 ZEBRA 应该在最前面
	if album02["name"] != "ZEBRA" {
		t.Errorf("Expected first album to be ZEBRA (desc), got %v", album02["name"])
	}
}

func TestAlbumHandler_GetAlbumSongs(t *testing.T) {
	db := createAlbumTestDB(t)

	// 创建测试数据
	songs := []*model.Song{
		{FilePath: "/test/song1.mp3", Title: "晴天", Album: "叶惠美", Artist: "周杰伦"},
		{FilePath: "/test/song2.mp3", Title: "稻香", Album: "叶惠美", Artist: "周杰伦"},
		{FilePath: "/test/song3.mp3", Title: "江南", Album: "江南", Artist: "林俊杰"},
	}
	for _, s := range songs {
		db.Create(s)
	}

	router := setupAlbumRouter(db)

	// 先获取专辑列表以确认 ID
	albumsReq, _ := http.NewRequest("GET", "/api/albums", nil)
	albumsW := httptest.NewRecorder()
	router.ServeHTTP(albumsW, albumsReq)

	var albumsResponse map[string]interface{}
	json.Unmarshal(albumsW.Body.Bytes(), &albumsResponse)
	albumsData := albumsResponse["data"].([]interface{})

	// 找到叶惠美的 ID (按字母排序叶惠美应该排第一)
	var yeID float64
	for _, a := range albumsData {
		album := a.(map[string]interface{})
		if album["name"] == "叶惠美" {
			yeID = album["id"].(float64)
			break
		}
	}

	// 测试获取特定专辑的歌曲
	req, _ := http.NewRequest("GET", "/api/albums/"+formatID(int(yeID))+"/songs", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	songsData, ok := response["data"].([]interface{})
	if !ok {
		t.Fatal("Expected data to be an array")
	}

	if len(songsData) != 2 {
		t.Errorf("Expected 2 songs for 叶惠美, got %d", len(songsData))
	}
}

func TestAlbumHandler_GetAlbumSongs_InvalidID(t *testing.T) {
	db := createAlbumTestDB(t)
	router := setupAlbumRouter(db)

	req, _ := http.NewRequest("GET", "/api/albums/999/songs", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestAlbumHandler_GetAlbumSongs_InvalidIDFormat(t *testing.T) {
	db := createAlbumTestDB(t)
	router := setupAlbumRouter(db)

	req, _ := http.NewRequest("GET", "/api/albums/abc/songs", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func formatID(id int) string {
	return fmt.Sprintf("%d", id)
}
