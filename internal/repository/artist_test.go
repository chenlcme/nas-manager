package repository

import (
	"testing"

	"nas-manager/internal/model"

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

func TestArtistRepository_GetAllArtistsWithSongCount(t *testing.T) {
	db := createArtistTestDB(t)
	repo := NewArtistRepository(db)

	// 创建测试数据
	songs := []*model.Song{
		{FilePath: "/test/song1.mp3", Title: "Song 1", Artist: "周杰伦"},
		{FilePath: "/test/song2.mp3", Title: "Song 2", Artist: "周杰伦"},
		{FilePath: "/test/song3.mp3", Title: "Song 3", Artist: "周杰伦"},
		{FilePath: "/test/song4.mp3", Title: "Song 4", Artist: "林俊杰"},
		{FilePath: "/test/song5.mp3", Title: "Song 5", Artist: "林俊杰"},
		{FilePath: "/test/song6.mp3", Title: "Song 6", Artist: "王力宏"},
	}
	for _, s := range songs {
		db.Create(s)
	}

	// 测试获取艺术家列表（降序）
	artists, err := repo.GetAllArtistsWithSongCount(false)
	if err != nil {
		t.Fatalf("Failed to get artists: %v", err)
	}

	if len(artists) != 3 {
		t.Errorf("Expected 3 artists, got %d", len(artists))
	}

	// 验证艺术家及其歌曲数量
	artistCounts := make(map[string]int)
	for _, a := range artists {
		artistCounts[a.Name] = a.SongCount
	}

	if artistCounts["周杰伦"] != 3 {
		t.Errorf("Expected 周杰伦 to have 3 songs, got %d", artistCounts["周杰伦"])
	}
	if artistCounts["林俊杰"] != 2 {
		t.Errorf("Expected 林俊杰 to have 2 songs, got %d", artistCounts["林俊杰"])
	}
	if artistCounts["王力宏"] != 1 {
		t.Errorf("Expected 王力宏 to have 1 song, got %d", artistCounts["王力宏"])
	}

	// 测试升序和降序返回的顺序不同（验证排序确实发生了）
	artistsAsc, err := repo.GetAllArtistsWithSongCount(true)
	if err != nil {
		t.Fatalf("Failed to get artists (asc): %v", err)
	}

	artistsDesc, err := repo.GetAllArtistsWithSongCount(false)
	if err != nil {
		t.Fatalf("Failed to get artists (desc): %v", err)
	}

	// 升序和降序的第一个艺术家应该不同（验证排序方向不同）
	if artistsAsc[0].Name == artistsDesc[0].Name {
		t.Error("Expected different first artist between asc and desc order")
	}
}

func TestArtistRepository_GetSongsByArtist(t *testing.T) {
	db := createArtistTestDB(t)
	repo := NewArtistRepository(db)

	// 创建测试数据
	songs := []*model.Song{
		{FilePath: "/test/song1.mp3", Title: "晴天", Artist: "周杰伦"},
		{FilePath: "/test/song2.mp3", Title: "稻香", Artist: "周杰伦"},
		{FilePath: "/test/song3.mp3", Title: "七里香", Artist: "周杰伦"},
		{FilePath: "/test/song4.mp3", Title: "江南", Artist: "林俊杰"},
	}
	for _, s := range songs {
		db.Create(s)
	}

	// 测试获取特定艺术家的歌曲
	found, err := repo.GetSongsByArtist("周杰伦")
	if err != nil {
		t.Fatalf("Failed to get songs by artist: %v", err)
	}

	if len(found) != 3 {
		t.Errorf("Expected 3 songs for 周杰伦, got %d", len(found))
	}

	// 验证返回的歌曲标题
	titles := make(map[string]bool)
	for _, s := range found {
		titles[s.Title] = true
	}

	if !titles["晴天"] || !titles["稻香"] || !titles["七里香"] {
		t.Error("Expected all 周杰伦 songs to be returned")
	}
}

func TestArtistRepository_GetSongsByArtist_NotFound(t *testing.T) {
	db := createArtistTestDB(t)
	repo := NewArtistRepository(db)

	found, err := repo.GetSongsByArtist("不存在的艺术家")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(found) != 0 {
		t.Errorf("Expected 0 songs for non-existent artist, got %d", len(found))
	}
}

func TestArtistRepository_GetAllArtistsWithSongCount_EmptyArtist(t *testing.T) {
	db := createArtistTestDB(t)
	repo := NewArtistRepository(db)

	// 创建测试数据，包含空艺术家
	songs := []*model.Song{
		{FilePath: "/test/song1.mp3", Title: "Song 1", Artist: "周杰伦"},
		{FilePath: "/test/song2.mp3", Title: "Song 2", Artist: ""},
		{FilePath: "/test/song3.mp3", Title: "Song 3", Artist: "   "},
	}
	for _, s := range songs {
		db.Create(s)
	}

	// 空艺术家不应被返回
	artists, err := repo.GetAllArtistsWithSongCount(false)
	if err != nil {
		t.Fatalf("Failed to get artists: %v", err)
	}

	if len(artists) != 1 {
		t.Errorf("Expected 1 artist (空艺术家应被过滤), got %d", len(artists))
	}

	if artists[0].Name != "周杰伦" {
		t.Errorf("Expected artist to be 周杰伦, got %s", artists[0].Name)
	}
}
