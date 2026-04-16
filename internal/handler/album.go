package handler

import (
	"net/http"
	"strconv"

	"nas-manager/internal/repository"
	"nas-manager/pkg/response"

	"github.com/gin-gonic/gin"
)

// AlbumHandler - 专辑相关 HTTP 处理
type AlbumHandler struct {
	albumRepo *repository.AlbumRepository
}

// NewAlbumHandler - 创建专辑处理器
func NewAlbumHandler(albumRepo *repository.AlbumRepository) *AlbumHandler {
	return &AlbumHandler{albumRepo: albumRepo}
}

// GetAlbums - 获取专辑列表
// GET /api/albums
func (h *AlbumHandler) GetAlbums(c *gin.Context) {
	// 获取排序方向，默认为降序
	orderAsc := c.Query("order") == "asc"

	albums, err := h.albumRepo.GetAllAlbumsWithSongCount(orderAsc)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "GET_ALBUMS_FAILED", "获取专辑列表失败")
		return
	}

	response.Success(c, albums)
}

// GetAlbumSongs - 获取特定专辑的歌曲列表
// GET /api/albums/:id/songs?sort_by=title&order=asc
func (h *AlbumHandler) GetAlbumSongs(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ALBUM_ID", "无效的专辑ID")
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

	// 通过 ID 获取专辑名和艺术家名（ID 是动态分配的，基于专辑列表顺序）
	albumName, artistName, err := h.albumRepo.GetAlbumNameAndArtistByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "ALBUM_NOT_FOUND", "专辑不存在")
		return
	}

	songs, err := h.albumRepo.GetSongsByAlbum(albumName, artistName, sortBy, order)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "GET_ALBUM_SONGS_FAILED", "获取专辑歌曲失败")
		return
	}

	response.Success(c, songs)
}
