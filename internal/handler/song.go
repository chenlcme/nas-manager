package handler

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"nas-manager/internal/repository"
	"nas-manager/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

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
