package handler

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"nas-manager/internal/model"
	"nas-manager/internal/repository"
	"nas-manager/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DeleteResult - 批量删除结果
type DeleteResult struct {
	Total     int                 `json:"total"`
	Succeeded int                 `json:"succeeded"`
	Failed    int                 `json:"failed"`
	Results   []SongDeleteResult `json:"results"`
}

// SongDeleteResult - 单个歌曲删除结果
type SongDeleteResult struct {
	ID       uint   `json:"id"`
	FilePath string `json:"file_path"`
	Status   string `json:"status"`
	Error    string `json:"error,omitempty"`
}

// SongHandler - 歌曲相关 HTTP 处理
type SongHandler struct {
	songRepo *repository.SongRepository
}

// NewSongHandler - 创建歌曲处理器
func NewSongHandler(songRepo *repository.SongRepository) *SongHandler {
	return &SongHandler{songRepo: songRepo}
}

// GetSong - 获取单曲详情
// GET /api/songs/:id
func (h *SongHandler) GetSong(c *gin.Context) {
	start := time.Now()
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		log.Printf("[SongHandler] Invalid ID: %s, error: %v", idStr, err)
		response.Error(c, http.StatusBadRequest, "INVALID_SONG_ID", "无效的歌曲ID")
		return
	}

	log.Printf("[SongHandler] GetSong id=%d", id)

	// Create context with timeout for DB query
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	song, err := h.songRepo.GetByIDWithContext(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("[SongHandler] Song not found: id=%d", id)
			response.Error(c, http.StatusNotFound, "SONG_NOT_FOUND", "歌曲不存在")
			return
		}
		log.Printf("[SongHandler] DB error for id=%d: %v", id, err)
		response.Error(c, http.StatusInternalServerError, "DB_ERROR", "数据库错误")
		return
	}

	log.Printf("[SongHandler] GetSong id=%d success, duration=%v", id, time.Since(start))
	response.Success(c, song)
}

// GetAllSongs - 获取所有歌曲（支持排序）
// GET /api/songs?sort_by=title&order=asc
func (h *SongHandler) GetAllSongs(c *gin.Context) {
	start := time.Now()

	sortBy := c.DefaultQuery("sort_by", "title")
	order := c.DefaultQuery("order", "asc")

	log.Printf("[SongHandler] GetAllSongs sort_by=%s order=%s", sortBy, order)

	songs, err := h.songRepo.GetAllSorted(sortBy, order)
	if err != nil {
		log.Printf("[SongHandler] GetAllSongs DB error: %v", err)
		response.Error(c, http.StatusInternalServerError, "DB_ERROR", "数据库错误")
		return
	}

	log.Printf("[SongHandler] GetAllSongs found=%d duration=%v", len(songs), time.Since(start))
	response.Success(c, songs)
}

// GetSongs - 批量获取歌曲详情
// POST /api/songs/batch-get
func (h *SongHandler) GetSongs(c *gin.Context) {
	start := time.Now()

	var req struct {
		IDs []uint `json:"ids" binding:"required,min=1,max=100"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[SongHandler] GetSongs invalid request: %v", err)
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "无效的请求参数")
		return
	}

	log.Printf("[SongHandler] GetSongs count=%d", len(req.IDs))

	// Deduplicate IDs
	seen := make(map[uint]bool)
	uniqueIDs := make([]uint, 0, len(req.IDs))
	for _, id := range req.IDs {
		if !seen[id] {
			seen[id] = true
			uniqueIDs = append(uniqueIDs, id)
		}
	}

	// Fetch songs
	songs := make([]model.Song, 0, len(uniqueIDs))
	for _, id := range uniqueIDs {
		song, err := h.songRepo.GetByID(id)
		if err != nil {
			continue // Skip songs that can't be found
		}
		songs = append(songs, *song)
	}

	log.Printf("[SongHandler] GetSongs found=%d duration=%v", len(songs), time.Since(start))
	response.Success(c, songs)
}

// DeleteSongs - 批量删除歌曲
// POST /api/songs/delete
func (h *SongHandler) DeleteSongs(c *gin.Context) {
	start := time.Now()

	var req struct {
		IDs []uint `json:"ids" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[SongHandler] DeleteSongs invalid request: %v", err)
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "无效的请求参数")
		return
	}

	// Deduplicate IDs to avoid processing the same song multiple times
	seen := make(map[uint]bool)
	uniqueIDs := make([]uint, 0, len(req.IDs))
	for _, id := range req.IDs {
		if !seen[id] {
			seen[id] = true
			uniqueIDs = append(uniqueIDs, id)
		}
	}

	log.Printf("[SongHandler] DeleteSongs count=%d after dedup=%d", len(req.IDs), len(uniqueIDs))

	result := &DeleteResult{
		Total:     len(uniqueIDs),
		Succeeded: 0,
		Failed:    0,
		Results:   make([]SongDeleteResult, 0, len(uniqueIDs)),
	}

	// Use mutex to protect results slice
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Process deletions concurrently with timeout control
	for _, id := range uniqueIDs {
		wg.Add(1)
		go func(songID uint) {
			defer wg.Done()

			song, err := h.songRepo.GetByID(songID)
			if err != nil {
				log.Printf("[SongHandler] DeleteSongs song not found: id=%d", songID)
				mu.Lock()
				result.Results = append(result.Results, SongDeleteResult{
					ID:     songID,
					Status: "failed",
					Error:  "song not found",
				})
				result.Failed++
				mu.Unlock()
				return
			}

			// Validate file path to prevent path traversal
			cleanPath := filepath.Clean(song.FilePath)
			if cleanPath != song.FilePath {
				log.Printf("[SongHandler] DeleteSongs suspicious path detected: id=%d", songID)
				mu.Lock()
				result.Results = append(result.Results, SongDeleteResult{
					ID:       songID,
					FilePath: song.FilePath,
					Status:   "failed",
					Error:    "invalid file path",
				})
				result.Failed++
				mu.Unlock()
				return
			}

			// Delete file with timeout - use buffered channel to prevent goroutine leak
			deleteCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			// Buffered channel so goroutine can exit even if we timeout
			done := make(chan error, 1)
			go func() {
				done <- os.Remove(song.FilePath)
			}()

			var removeErr error
			select {
			case removeErr = <-done:
				// Operation completed
			case <-deleteCtx.Done():
				// Timeout - the goroutine will eventually send and exit (buffered channel)
				// We don't read from done, just use timeout error
				removeErr = deleteCtx.Err()
			}

			if removeErr != nil {
				log.Printf("[SongHandler] DeleteSongs file deletion failed: id=%d, err=%v", songID, removeErr)
				mu.Lock()
				result.Results = append(result.Results, SongDeleteResult{
					ID:       songID,
					FilePath: song.FilePath,
					Status:   "failed",
					Error:    "file delete failed",
				})
				result.Failed++
				mu.Unlock()
				// 即使文件删除失败，仍然尝试删除数据库记录
				if dbErr := h.songRepo.Delete(songID); dbErr != nil {
					log.Printf("[SongHandler] DeleteSongs db record deletion failed: id=%d", songID)
				}
				return
			}

			// 删除数据库记录
			if err := h.songRepo.Delete(songID); err != nil {
				log.Printf("[SongHandler] DeleteSongs db record deletion failed: id=%d", songID)
				mu.Lock()
				result.Results = append(result.Results, SongDeleteResult{
					ID:       songID,
					FilePath: song.FilePath,
					Status:   "failed",
					Error:    "db delete failed",
				})
				result.Failed++
				mu.Unlock()
				return
			}

			log.Printf("[SongHandler] DeleteSongs deleted: id=%d", songID)
			mu.Lock()
			result.Results = append(result.Results, SongDeleteResult{
				ID:       songID,
				FilePath: song.FilePath,
				Status:   "deleted",
			})
			result.Succeeded++
			mu.Unlock()
		}(id)
	}

	wg.Wait()

	log.Printf("[SongHandler] DeleteSongs completed: total=%d, succeeded=%d, failed=%d, duration=%v",
		result.Total, result.Succeeded, result.Failed, time.Since(start))
	response.Success(c, result)
}

// SearchSongs - 按文件名搜索歌曲
// GET /api/songs/search?q=keyword&limit=20&offset=0
func (h *SongHandler) SearchSongs(c *gin.Context) {
	start := time.Now()

	keyword := c.Query("q")
	if keyword == "" {
		log.Printf("[SongHandler] SearchSongs missing query parameter 'q'")
		response.Error(c, http.StatusBadRequest, "MISSING_QUERY", "搜索关键词不能为空")
		return
	}

	// Parse pagination parameters
	limit := 20 // default limit
	offset := 0 // default offset
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	// Enforce maximum limit to prevent excessive queries
	const maxLimit = 1000
	limit = min(limit, maxLimit)

	log.Printf("[SongHandler] SearchSongs keyword=%s limit=%d offset=%d", keyword, limit, offset)

	songs, err := h.songRepo.SearchByFileName(keyword, limit, offset)
	if err != nil {
		log.Printf("[SongHandler] SearchSongs DB error: %v", err)
		response.Error(c, http.StatusInternalServerError, "DB_ERROR", "数据库错误")
		return
	}

	log.Printf("[SongHandler] SearchSongs keyword=%s found=%d duration=%v", keyword, len(songs), time.Since(start))
	response.Success(c, songs)
}

// SearchSongsByTag - 按标签内容搜索歌曲（标题、艺术家、专辑）
// GET /api/songs/search/by-tag?q=keyword&limit=20&offset=0
func (h *SongHandler) SearchSongsByTag(c *gin.Context) {
	start := time.Now()

	keyword := c.Query("q")
	if keyword == "" {
		log.Printf("[SongHandler] SearchSongsByTag missing query parameter 'q'")
		response.Error(c, http.StatusBadRequest, "MISSING_QUERY", "搜索关键词不能为空")
		return
	}

	// Parse pagination parameters
	limit := 20 // default limit
	offset := 0 // default offset
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	// Enforce maximum limit to prevent excessive queries
	const maxLimit = 1000
	limit = min(limit, maxLimit)

	log.Printf("[SongHandler] SearchSongsByTag keyword=%s limit=%d offset=%d", keyword, limit, offset)

	// 检查是否包含空格，支持多条件组合搜索
	keywords := strings.Fields(keyword)
	var songs []model.Song
	var err error

	if len(keywords) > 1 {
		// 多关键词搜索
		songs, err = h.songRepo.SearchByTagContentMulti(keywords, limit, offset)
	} else {
		// 单关键词搜索
		songs, err = h.songRepo.SearchByTagContent(keyword, limit, offset)
	}

	if err != nil {
		log.Printf("[SongHandler] SearchSongsByTag DB error: %v", err)
		response.Error(c, http.StatusInternalServerError, "DB_ERROR", "数据库错误")
		return
	}

	log.Printf("[SongHandler] SearchSongsByTag keyword=%s found=%d duration=%v", keyword, len(songs), time.Since(start))
	response.Success(c, songs)
}

// UpdateSong - 更新歌曲信息
// PUT /api/songs/:id
func (h *SongHandler) UpdateSong(c *gin.Context) {
	start := time.Now()
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		log.Printf("[SongHandler] UpdateSong invalid ID: %s, error: %v", idStr, err)
		response.Error(c, http.StatusBadRequest, "INVALID_SONG_ID", "无效的歌曲ID")
		return
	}

	log.Printf("[SongHandler] UpdateSong id=%d", id)

	// Get existing song
	song, err := h.songRepo.GetByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("[SongHandler] UpdateSong song not found: id=%d", id)
			response.Error(c, http.StatusNotFound, "SONG_NOT_FOUND", "歌曲不存在")
			return
		}
		log.Printf("[SongHandler] UpdateSong DB error for id=%d: %v", id, err)
		response.Error(c, http.StatusInternalServerError, "DB_ERROR", "数据库错误")
		return
	}

	// Parse update request
	var req struct {
		Title     *string `json:"title"`
		Artist    *string `json:"artist"`
		Album     *string `json:"album"`
		Year      *int    `json:"year"`
		Genre     *string `json:"genre"`
		TrackNum  *int    `json:"trackNum"`
		CoverPath *string `json:"coverPath"`
		Lyrics    *string `json:"lyrics"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[SongHandler] UpdateSong invalid request: %v", err)
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "无效的请求参数")
		return
	}

	// Update fields if provided (nil means don't update)
	if req.Title != nil {
		song.Title = *req.Title
	}
	if req.Artist != nil {
		song.Artist = *req.Artist
	}
	if req.Album != nil {
		song.Album = *req.Album
	}
	if req.Year != nil {
		song.Year = *req.Year
	}
	if req.Genre != nil {
		song.Genre = *req.Genre
	}
	if req.TrackNum != nil {
		song.TrackNum = *req.TrackNum
	}
	if req.CoverPath != nil {
		song.CoverPath = *req.CoverPath
	}
	if req.Lyrics != nil {
		song.Lyrics = *req.Lyrics
	}

	// Save updates
	if err := h.songRepo.Update(song); err != nil {
		log.Printf("[SongHandler] UpdateSong failed to update id=%d: %v", id, err)
		response.Error(c, http.StatusInternalServerError, "UPDATE_FAILED", "更新失败")
		return
	}

	log.Printf("[SongHandler] UpdateSong id=%d success, duration=%v", id, time.Since(start))
	response.Success(c, song)
}

// StreamSong - 流式播放歌曲
// GET /api/songs/:id/stream
func (h *SongHandler) StreamSong(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		log.Printf("[SongHandler] StreamSong invalid ID: %s, error: %v", idStr, err)
		response.Error(c, http.StatusBadRequest, "INVALID_SONG_ID", "无效的歌曲ID")
		return
	}

	log.Printf("[SongHandler] StreamSong id=%d", id)

	// Get song to find file path
	song, err := h.songRepo.GetByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("[SongHandler] StreamSong song not found: id=%d", id)
			response.Error(c, http.StatusNotFound, "SONG_NOT_FOUND", "歌曲不存在")
			return
		}
		log.Printf("[SongHandler] StreamSong DB error for id=%d: %v", id, err)
		response.Error(c, http.StatusInternalServerError, "DB_ERROR", "数据库错误")
		return
	}

	// Check if file exists
	if _, err := os.Stat(song.FilePath); os.IsNotExist(err) {
		log.Printf("[SongHandler] StreamSong file not found: %s", song.FilePath)
		response.Error(c, http.StatusNotFound, "FILE_NOT_FOUND", "文件不存在")
		return
	}

	// Determine content type based on file extension
	ext := strings.ToLower(filepath.Ext(song.FilePath))
	var contentType string
	switch ext {
	case ".mp3":
		contentType = "audio/mpeg"
	case ".flac":
		contentType = "audio/flac"
	case ".ogg":
		contentType = "audio/ogg"
	case ".m4a", ".aac":
		contentType = "audio/mp4"
	case ".wav":
		contentType = "audio/wav"
	case ".ape":
		contentType = "audio/ape"
	default:
		contentType = "audio/mpeg"
	}

	// Set headers for streaming
	c.Header("Content-Type", contentType)
	c.Header("Accept-Ranges", "bytes")

	// Stream the file
	c.File(song.FilePath)
}
