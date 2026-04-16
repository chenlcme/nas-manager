package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"

	"nas-manager/internal/handler"
	"nas-manager/internal/model"
	"nas-manager/internal/repository"
	"nas-manager/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	dbPath string
)

func init() {
	flag.StringVar(&dbPath, "db", "", "SQLite database path")
}

func getDefaultDBPath() string {
	usr, err := user.Current()
	if err != nil {
		return ".nas-manager.db"
	}
	homeDir := filepath.Join(usr.HomeDir, ".nas-manager")
	if err := os.MkdirAll(homeDir, 0755); err != nil {
		return filepath.Join(homeDir, "nas-manager.db")
	}
	return filepath.Join(homeDir, "nas-manager.db")
}

func main() {
	flag.Parse()

	// Resolve database path
	if dbPath == "" {
		if envPath := os.Getenv("NAS_MANAGER_DB"); envPath != "" {
			dbPath = envPath
		} else {
			dbPath = getDefaultDBPath()
		}
	}

	// Connect to SQLite database
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate all models
	if err := db.AutoMigrate(
		&model.Song{},
		&model.Artist{},
		&model.Album{},
		&model.Setting{},
		&model.BatchOperation{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	fmt.Printf("Database initialized at: %s\n", dbPath)

	// Initialize repositories
	settingRepo := repository.NewSettingRepository(db)
	songRepo := repository.NewSongRepository(db)
	artistRepo := repository.NewArtistRepository(db)
	albumRepo := repository.NewAlbumRepository(db)
	folderRepo := repository.NewFolderRepository(db)

	// Initialize services
	settingService := service.NewSettingService(settingRepo)
	id3Service := service.NewID3Service(songRepo)
	scannerService := service.NewScannerService(id3Service, songRepo)

	// Initialize handlers
	settingHandler := handler.NewSettingHandler(settingService)
	scanHandler := handler.NewScanHandler(scannerService, songRepo, settingRepo)
	encryptService := service.NewEncryptService(settingRepo)
	encryptHandler := handler.NewEncryptHandler(encryptService)
	artistHandler := handler.NewArtistHandler(artistRepo)
	albumHandler := handler.NewAlbumHandler(albumRepo)
	folderHandler := handler.NewFolderHandler(folderRepo)

	// Setup Gin router
	r := gin.Default()

	// API routes
	api := r.Group("/api")
	{
		// Setup routes
		api.GET("/setup/status", settingHandler.GetSetupStatus)
		api.POST("/setup", settingHandler.SaveSetup)

		// Scan routes
		api.POST("/songs/scan", scanHandler.Scan)
		api.POST("/songs/cleanup", scanHandler.Cleanup)

		// Auth routes
		api.POST("/auth/setup", encryptHandler.SetupPassword)
		api.POST("/auth/verify", encryptHandler.VerifyPassword)
		api.POST("/auth/change", encryptHandler.ChangePassword)

		// Artist routes
		api.GET("/artists", artistHandler.GetArtists)
		api.GET("/artists/:id/songs", artistHandler.GetArtistSongs)

		// Album routes
		api.GET("/albums", albumHandler.GetAlbums)
		api.GET("/albums/:id/songs", albumHandler.GetAlbumSongs)

		// Folder routes
		api.GET("/folders", folderHandler.GetFolders)
		api.GET("/folders/:id/songs", folderHandler.GetFolderSongs)
	}

	fmt.Println("nas-manager server starting on :8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
