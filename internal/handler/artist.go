package handler

import (
	"net/http"
	"strconv"

	"nas-manager/internal/repository"
	"nas-manager/pkg/response"

	"github.com/gin-gonic/gin"
)

// ArtistHandler - 艺术家相关 HTTP 处理
type ArtistHandler struct {
	artistRepo *repository.ArtistRepository
}

// NewArtistHandler - 创建艺术家处理器
func NewArtistHandler(artistRepo *repository.ArtistRepository) *ArtistHandler {
	return &ArtistHandler{artistRepo: artistRepo}
}

// GetArtists - 获取艺术家列表
// GET /api/artists
func (h *ArtistHandler) GetArtists(c *gin.Context) {
	// 获取排序方向，默认为降序
	orderAsc := c.Query("order") == "asc"

	artists, err := h.artistRepo.GetAllArtistsWithSongCount(orderAsc)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "GET_ARTISTS_FAILED", "获取艺术家列表失败")
		return
	}

	response.Success(c, artists)
}

// GetArtistSongs - 获取特定艺术家的歌曲列表
// GET /api/artists/:id/songs?sort_by=title&order=asc
func (h *ArtistHandler) GetArtistSongs(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ARTIST_ID", "无效的艺术家ID")
		return
	}

	// 获取排序参数
	sortBy := c.Query("sort_by")
	order := c.Query("order")

	// 校验 sort_by 参数
	validSortFields := map[string]bool{
		"title":     true,
		"duration":   true,
		"created_at": true,
	}
	validOrders := map[string]bool{"asc": true, "desc": true}

	if !validSortFields[sortBy] {
		sortBy = "title"
	}
	if !validOrders[order] {
		order = "asc"
	}

	// 通过 ID 获取艺术家名（ID 是动态分配的，基于艺术家列表顺序）
	artistName, err := h.artistRepo.GetArtistNameByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "ARTIST_NOT_FOUND", "艺术家不存在")
		return
	}

	songs, err := h.artistRepo.GetSongsByArtist(artistName, sortBy, order)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "GET_ARTIST_SONGS_FAILED", "获取艺术家歌曲失败")
		return
	}

	response.Success(c, songs)
}
