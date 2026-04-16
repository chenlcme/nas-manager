package handler

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

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

	log.Printf("[SongHandler] DeleteSongs count=%d", len(req.IDs))

	result := &DeleteResult{
		Total:     len(req.IDs),
		Succeeded: 0,
		Failed:    0,
		Results:   make([]SongDeleteResult, 0, len(req.IDs)),
	}

	// Use mutex to protect results slice
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Process deletions concurrently with timeout control
	for _, id := range req.IDs {
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

			// Delete file with timeout
			deleteCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			done := make(chan error, 1)
			go func() {
				done <- os.Remove(song.FilePath)
			}()

			select {
			case err := <-done:
				if err != nil {
					log.Printf("[SongHandler] DeleteSongs file deletion failed: id=%d", songID)
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
			case <-deleteCtx.Done():
				log.Printf("[SongHandler] DeleteSongs file deletion timeout: id=%d", songID)
				mu.Lock()
				result.Results = append(result.Results, SongDeleteResult{
					ID:       songID,
					FilePath: song.FilePath,
					Status:   "failed",
					Error:    "file deletion timeout",
				})
				result.Failed++
				mu.Unlock()
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
