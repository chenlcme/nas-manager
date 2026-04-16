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

func createArtistTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	db.AutoMigrate(&model.Song{})

	return db
}

func setupArtistRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	artistRepo := repository.NewArtistRepository(db)
	artistHandler := NewArtistHandler(artistRepo)

	r.GET("/api/artists", artistHandler.GetArtists)
	r.GET("/api/artists/:id/songs", artistHandler.GetArtistSongs)

	return r
}

func TestArtistHandler_GetArtists(t *testing.T) {
	db := createArtistTestDB(t)

	// 创建测试数据
	songs := []*model.Song{
		{FilePath: "/test/song1.mp3", Title: "Song 1", Artist: "周杰伦"},
		{FilePath: "/test/song2.mp3", Title: "Song 2", Artist: "周杰伦"},
		{FilePath: "/test/song3.mp3", Title: "Song 3", Artist: "林俊杰"},
	}
	for _, s := range songs {
		db.Create(s)
	}

	router := setupArtistRouter(db)

	// 测试获取艺术家列表
	req, _ := http.NewRequest("GET", "/api/artists", nil)
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
		t.Errorf("Expected 2 artists, got %d", len(data))
	}
}

func TestArtistHandler_GetArtists_Empty(t *testing.T) {
	db := createArtistTestDB(t)
	router := setupArtistRouter(db)

	req, _ := http.NewRequest("GET", "/api/artists", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// 当没有艺术家时，response.Success 被调用返回空数组
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	// 检查响应结构
	t.Logf("Response body: %s", w.Body.String())

	data := response["data"]
	if data == nil {
		// 空数组返回 nil 而不是 []
		t.Log("Data is nil, which is equivalent to empty array - PASS")
		return
	}

	dataSlice, ok := data.([]interface{})
	if !ok {
		t.Fatalf("Expected data to be an array, got %T", data)
	}

	if len(dataSlice) != 0 {
		t.Errorf("Expected 0 artists, got %d", len(dataSlice))
	}
}

func TestArtistHandler_GetArtists_SortAsc(t *testing.T) {
	db := createArtistTestDB(t)

	// 创建测试数据 - 使用全大写艺术家名以便测试排序
	songs := []*model.Song{
		{FilePath: "/test/song1.mp3", Title: "Song 1", Artist: "ZEBRA"},
		{FilePath: "/test/song2.mp3", Title: "Song 2", Artist: "APPLE"},
		{FilePath: "/test/song3.mp3", Title: "Song 3", Artist: "MANGO"},
	}
	for _, s := range songs {
		db.Create(s)
	}

	router := setupArtistRouter(db)

	// 测试升序排序 - ASCII 顺序: APPLE < MANGO < ZEBRA
	req, _ := http.NewRequest("GET", "/api/artists?order=asc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data := response["data"].([]interface{})
	artist0 := data[0].(map[string]interface{})

	// 升序时 APPLE 应该在最前面
	if artist0["name"] != "APPLE" {
		t.Errorf("Expected first artist to be APPLE (asc), got %v", artist0["name"])
	}

	// 降序排序 - ASCII 顺序: ZEBRA > MANGO > APPLE
	req2, _ := http.NewRequest("GET", "/api/artists?order=desc", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	var response2 map[string]interface{}
	json.Unmarshal(w2.Body.Bytes(), &response2)

	data2 := response2["data"].([]interface{})
	artist02 := data2[0].(map[string]interface{})

	// 降序时 ZEBRA 应该在最前面
	if artist02["name"] != "ZEBRA" {
		t.Errorf("Expected first artist to be ZEBRA (desc), got %v", artist02["name"])
	}
}

func TestArtistHandler_GetArtistSongs(t *testing.T) {
	db := createArtistTestDB(t)

	// 创建测试数据
	songs := []*model.Song{
		{FilePath: "/test/song1.mp3", Title: "晴天", Artist: "周杰伦"},
		{FilePath: "/test/song2.mp3", Title: "稻香", Artist: "周杰伦"},
		{FilePath: "/test/song3.mp3", Title: "江南", Artist: "林俊杰"},
	}
	for _, s := range songs {
		db.Create(s)
	}

	router := setupArtistRouter(db)

	// 先获取艺术家列表以确认 ID
	artistsReq, _ := http.NewRequest("GET", "/api/artists", nil)
	artistsW := httptest.NewRecorder()
	router.ServeHTTP(artistsW, artistsReq)

	var artistsResponse map[string]interface{}
	json.Unmarshal(artistsW.Body.Bytes(), &artistsResponse)
	artistsData := artistsResponse["data"].([]interface{})

	// 找到周杰伦的 ID (应该是 1，因为按字母排序周杰伦排在最后)
	var zhouID float64
	for _, a := range artistsData {
		artist := a.(map[string]interface{})
		if artist["name"] == "周杰伦" {
			zhouID = artist["id"].(float64)
			break
		}
	}

	// 测试获取特定艺术家的歌曲
	req, _ := http.NewRequest("GET", "/api/artists/"+bytesmarshalID(int(zhouID))+"/songs", nil)
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
		t.Errorf("Expected 2 songs for 周杰伦, got %d", len(songsData))
	}
}

func TestArtistHandler_GetArtistSongs_InvalidID(t *testing.T) {
	db := createArtistTestDB(t)
	router := setupArtistRouter(db)

	req, _ := http.NewRequest("GET", "/api/artists/999/songs", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestArtistHandler_GetArtistSongs_InvalidIDFormat(t *testing.T) {
	db := createArtistTestDB(t)
	router := setupArtistRouter(db)

	req, _ := http.NewRequest("GET", "/api/artists/abc/songs", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func bytesmarshalID(id int) string {
	return fmt.Sprintf("%d", id)
}
