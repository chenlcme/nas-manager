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
