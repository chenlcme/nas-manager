package repository

import (
	"testing"

	"nas-manager/internal/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func createSongTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	db.AutoMigrate(&model.Song{})

	return db
}

func TestSongRepository_Create(t *testing.T) {
	db := createSongTestDB(t)
	repo := NewSongRepository(db)

	song := &model.Song{
		FilePath: "/test/song.mp3",
		Title:    "Test Song",
		Artist:   "Test Artist",
	}

	if err := repo.Create(song); err != nil {
		t.Fatalf("Failed to create song: %v", err)
	}

	if song.ID == 0 {
		t.Error("Expected song ID to be set")
	}
}

func TestSongRepository_ExistsByFilePath(t *testing.T) {
	db := createSongTestDB(t)
	repo := NewSongRepository(db)

	exists, err := repo.ExistsByFilePath("/test/song.mp3")
	if err != nil {
		t.Fatalf("Failed to check existence: %v", err)
	}
	if exists {
		t.Error("Expected song to not exist initially")
	}

	// 创建歌曲
	song := &model.Song{FilePath: "/test/song.mp3", Title: "Test"}
	repo.Create(song)

	exists, err = repo.ExistsByFilePath("/test/song.mp3")
	if err != nil {
		t.Fatalf("Failed to check existence: %v", err)
	}
	if !exists {
		t.Error("Expected song to exist after creation")
	}
}

func TestSongRepository_GetByFilePath(t *testing.T) {
	db := createSongTestDB(t)
	repo := NewSongRepository(db)

	song := &model.Song{
		FilePath: "/test/song.mp3",
		Title:    "Test Song",
		Artist:   "Test Artist",
	}
	repo.Create(song)

	found, err := repo.GetByFilePath("/test/song.mp3")
	if err != nil {
		t.Fatalf("Failed to get song: %v", err)
	}

	if found.Title != song.Title {
		t.Errorf("Expected title %s, got %s", song.Title, found.Title)
	}
}

func TestSongRepository_GetByArtist(t *testing.T) {
	db := createSongTestDB(t)
	repo := NewSongRepository(db)

	// 创建多个歌曲
	songs := []*model.Song{
		{FilePath: "/test/song1.mp3", Title: "Song 1", Artist: "Artist A"},
		{FilePath: "/test/song2.mp3", Title: "Song 2", Artist: "Artist A"},
		{FilePath: "/test/song3.mp3", Title: "Song 3", Artist: "Artist B"},
	}
	for _, s := range songs {
		repo.Create(s)
	}

	found, err := repo.GetByArtist("Artist A")
	if err != nil {
		t.Fatalf("Failed to get songs by artist: %v", err)
	}

	if len(found) != 2 {
		t.Errorf("Expected 2 songs, got %d", len(found))
	}
}

func TestSongRepository_SearchByFileName(t *testing.T) {
	db := createSongTestDB(t)
	repo := NewSongRepository(db)

	// Create test songs
	songs := []*model.Song{
		{FilePath: "/music/rock/晴天.mp3", Title: "晴天"},
		{FilePath: "/music/pop/夜曲.mp3", Title: "夜曲"},
		{FilePath: "/music/rock/七里香.mp3", Title: "七里香"},
		{FilePath: "/music/classic/梁祝.mp3", Title: "梁祝"},
	}
	for _, s := range songs {
		repo.Create(s)
	}

	// Search by Chinese filename
	results, err := repo.SearchByFileName("晴天", 20, 0)
	if err != nil {
		t.Fatalf("Failed to search: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results))
	}
	if len(results) > 0 && results[0].Title != "晴天" {
		t.Errorf("Expected title '晴天', got '%s'", results[0].Title)
	}
}

func TestSongRepository_SearchByFileName_MultipleResults(t *testing.T) {
	db := createSongTestDB(t)
	repo := NewSongRepository(db)

	// Create test songs - all in same folder
	songs := []*model.Song{
		{FilePath: "/music/周杰伦/晴天.mp3", Title: "晴天"},
		{FilePath: "/music/周杰伦/夜曲.mp3", Title: "夜曲"},
		{FilePath: "/music/周杰伦/七里香.mp3", Title: "七里香"},
		{FilePath: "/music/林俊傑/江南.mp3", Title: "江南"},
	}
	for _, s := range songs {
		repo.Create(s)
	}

	// Search by folder name (should match 3 songs)
	results, err := repo.SearchByFileName("周杰伦", 20, 0)
	if err != nil {
		t.Fatalf("Failed to search: %v", err)
	}
	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	}
}

func TestSongRepository_SearchByFileName_NoResults(t *testing.T) {
	db := createSongTestDB(t)
	repo := NewSongRepository(db)

	// Create test songs
	songs := []*model.Song{
		{FilePath: "/music/rock/晴天.mp3", Title: "晴天"},
		{FilePath: "/music/pop/夜曲.mp3", Title: "夜曲"},
	}
	for _, s := range songs {
		repo.Create(s)
	}

	// Search for non-existent keyword
	results, err := repo.SearchByFileName("不存在", 20, 0)
	if err != nil {
		t.Fatalf("Failed to search: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("Expected 0 results, got %d", len(results))
	}
}

func TestSongRepository_SearchByFileName_EmptyKeyword(t *testing.T) {
	db := createSongTestDB(t)
	repo := NewSongRepository(db)

	// Create test songs
	songs := []*model.Song{
		{FilePath: "/music/rock/晴天.mp3", Title: "晴天"},
	}
	for _, s := range songs {
		repo.Create(s)
	}

	// Search with empty keyword - repository doesn't validate, handler rejects empty
	// Repository just executes the query, so %% matches all (handler validates before calling)
	results, err := repo.SearchByFileName("", 20, 0)
	if err != nil {
		t.Fatalf("Failed to search: %v", err)
	}
	// Repository behavior: empty keyword with %% matches all
	if len(results) != 1 {
		t.Errorf("Expected 1 result for empty keyword at repository level, got %d", len(results))
	}
}
