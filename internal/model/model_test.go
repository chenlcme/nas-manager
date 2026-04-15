package model

import (
	"os"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestSongModel(t *testing.T) {
	// Create a temporary database for testing
	tmpFile, err := os.CreateTemp("", "test-*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	db, err := gorm.Open(sqlite.Open(tmpFile.Name()), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Migrate the Song model
	if err := db.AutoMigrate(&Song{}); err != nil {
		t.Fatalf("Failed to migrate Song model: %v", err)
	}

	// Create a song
	song := Song{
		FilePath: "/music/test.mp3",
		Title:    "Test Song",
		Artist:   "Test Artist",
		Album:    "Test Album",
		Year:     2024,
		Genre:    "Test",
		Duration: 180,
	}

	if result := db.Create(&song); result.Error != nil {
		t.Fatalf("Failed to create song: %v", result.Error)
	}

	// Verify the song was created
	var found Song
	if err := db.First(&found, song.ID).Error; err != nil {
		t.Fatalf("Failed to find created song: %v", err)
	}

	if found.Title != song.Title {
		t.Errorf("Expected title %s, got %s", song.Title, found.Title)
	}
	if found.Artist != song.Artist {
		t.Errorf("Expected artist %s, got %s", song.Artist, found.Artist)
	}
}

func TestArtistModel(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	db, err := gorm.Open(sqlite.Open(tmpFile.Name()), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	if err := db.AutoMigrate(&Artist{}); err != nil {
		t.Fatalf("Failed to migrate Artist model: %v", err)
	}

	artist := Artist{Name: "Test Artist"}
	if result := db.Create(&artist); result.Error != nil {
		t.Fatalf("Failed to create artist: %v", result.Error)
	}

	var found Artist
	if err := db.First(&found, artist.ID).Error; err != nil {
		t.Fatalf("Failed to find created artist: %v", err)
	}

	if found.Name != artist.Name {
		t.Errorf("Expected name %s, got %s", artist.Name, found.Name)
	}
}

func TestAlbumModel(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	db, err := gorm.Open(sqlite.Open(tmpFile.Name()), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	if err := db.AutoMigrate(&Album{}); err != nil {
		t.Fatalf("Failed to migrate Album model: %v", err)
	}

	album := Album{Name: "Test Album", Artist: "Test Artist"}
	if result := db.Create(&album); result.Error != nil {
		t.Fatalf("Failed to create album: %v", result.Error)
	}

	var found Album
	if err := db.First(&found, album.ID).Error; err != nil {
		t.Fatalf("Failed to find created album: %v", err)
	}

	if found.Name != album.Name {
		t.Errorf("Expected name %s, got %s", album.Name, found.Name)
	}
}

func TestSettingModel(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	db, err := gorm.Open(sqlite.Open(tmpFile.Name()), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	if err := db.AutoMigrate(&Setting{}); err != nil {
		t.Fatalf("Failed to migrate Setting model: %v", err)
	}

	setting := Setting{Key: "music_dir", Value: "/music"}
	if result := db.Create(&setting); result.Error != nil {
		t.Fatalf("Failed to create setting: %v", result.Error)
	}

	var found Setting
	if err := db.Where("key = ?", setting.Key).First(&found).Error; err != nil {
		t.Fatalf("Failed to find created setting: %v", err)
	}

	if found.Value != setting.Value {
		t.Errorf("Expected value %s, got %s", setting.Value, found.Value)
	}
}

func TestBatchOperationModel(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	db, err := gorm.Open(sqlite.Open(tmpFile.Name()), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	if err := db.AutoMigrate(&BatchOperation{}); err != nil {
		t.Fatalf("Failed to migrate BatchOperation model: %v", err)
	}

	batch := BatchOperation{
		Type:      "update",
		TargetIDs: "[1, 2, 3]",
		OldValues: `{"artist": "Old Artist"}`,
		NewValues: `{"artist": "New Artist"}`,
	}

	if result := db.Create(&batch); result.Error != nil {
		t.Fatalf("Failed to create batch operation: %v", result.Error)
	}

	var found BatchOperation
	if err := db.First(&found, batch.ID).Error; err != nil {
		t.Fatalf("Failed to find created batch operation: %v", err)
	}

	if found.Type != batch.Type {
		t.Errorf("Expected type %s, got %s", batch.Type, found.Type)
	}
}

func TestAllModelsAutoMigrate(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	db, err := gorm.Open(sqlite.Open(tmpFile.Name()), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto-migrate all models
	if err := db.AutoMigrate(
		&Song{},
		&Artist{},
		&Album{},
		&Setting{},
		&BatchOperation{},
	); err != nil {
		t.Fatalf("Failed to migrate all models: %v", err)
	}

	// Verify all tables exist by inserting and querying
	song := Song{FilePath: "/test/song.mp3", Title: "Test"}
	if err := db.Create(&song).Error; err != nil {
		t.Errorf("Failed to create song after migration: %v", err)
	}

	artist := Artist{Name: "Test Artist"}
	if err := db.Create(&artist).Error; err != nil {
		t.Errorf("Failed to create artist after migration: %v", err)
	}

	album := Album{Name: "Test Album", Artist: "Test Artist"}
	if err := db.Create(&album).Error; err != nil {
		t.Errorf("Failed to create album after migration: %v", err)
	}

	setting := Setting{Key: "test_key", Value: "test_value"}
	if err := db.Create(&setting).Error; err != nil {
		t.Errorf("Failed to create setting after migration: %v", err)
	}

	batch := BatchOperation{Type: "test", TargetIDs: "[]"}
	if err := db.Create(&batch).Error; err != nil {
		t.Errorf("Failed to create batch after migration: %v", err)
	}
}
