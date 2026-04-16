package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"nas-manager/internal/handler"
	"nas-manager/internal/model"
	"nas-manager/pkg/response"
	"nas-manager/internal/repository"
	"nas-manager/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	dbPath         string
	frontendDistDir string
)

func init() {
	flag.StringVar(&dbPath, "db", "", "SQLite database path")
	flag.StringVar(&frontendDistDir, "frontend-dist", "", "Frontend dist directory (default: ./frontend/dist)")
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
	batchRepo := repository.NewBatchRepository(db)

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
	songHandler := handler.NewSongHandler(songRepo)
	batchHandler := handler.NewBatchHandler(songRepo, batchRepo)

	// Setup Gin router
	r := gin.Default()

	// Determine frontend dist directory
	if frontendDistDir == "" {
		frontendDistDir = filepath.Join(".", "frontend", "dist")
	}

	// Serve frontend static files
	if _, err := os.Stat(frontendDistDir); err == nil {
		r.GET("/", func(c *gin.Context) {
			c.File(filepath.Join(frontendDistDir, "index.html"))
		})
		r.Static("/assets", filepath.Join(frontendDistDir, "assets"))

		// SPA fallback: serve index.html for non-API routes
		r.NoRoute(func(c *gin.Context) {
			path := c.Request.URL.Path
			if !strings.HasPrefix(path, "/api") {
				// Check if request is for a static file
				if strings.HasSuffix(path, ".js") || strings.HasSuffix(path, ".css") || strings.HasSuffix(path, ".ico") || strings.HasSuffix(path, ".png") || strings.HasSuffix(path, ".jpg") || strings.HasSuffix(path, ".svg") {
					c.File(filepath.Join(frontendDistDir, path))
					return
				}
				// Otherwise serve index.html for SPA routing
				c.File(filepath.Join(frontendDistDir, "index.html"))
			} else {
				response.Error(c, http.StatusNotFound, "NOT_FOUND", "请求的接口不存在")
			}
		})
	}

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

		// Song routes
		api.GET("/songs/search", songHandler.SearchSongs)
		api.GET("/songs/search/by-tag", songHandler.SearchSongsByTag)
		api.GET("/songs/:id", songHandler.GetSong)
		api.PUT("/songs/:id", songHandler.UpdateSong)
		api.GET("/songs/:id/stream", songHandler.StreamSong)
		api.POST("/songs/delete", songHandler.DeleteSongs)
		api.POST("/songs/batch-get", songHandler.GetSongs)
		api.POST("/songs/batch-update", batchHandler.BatchUpdate)
		api.POST("/songs/undo/:batchId", batchHandler.UndoBatch)
		api.GET("/batches/latest", batchHandler.GetLatestBatch)
	}

	fmt.Println("nas-manager server starting on :8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
