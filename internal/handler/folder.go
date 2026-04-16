package handler

import (
	"net/http"
	"strconv"

	"nas-manager/internal/repository"
	"nas-manager/pkg/response"

	"github.com/gin-gonic/gin"
)

// FolderHandler - 文件夹相关 HTTP 处理
type FolderHandler struct {
	folderRepo *repository.FolderRepository
}

// NewFolderHandler - 创建文件夹处理器
func NewFolderHandler(folderRepo *repository.FolderRepository) *FolderHandler {
	return &FolderHandler{folderRepo: folderRepo}
}

// GetFolders - 获取文件夹列表
// GET /api/folders
func (h *FolderHandler) GetFolders(c *gin.Context) {
	// 获取排序方向，默认为降序
	orderAsc := c.Query("order") == "asc"

	folders, err := h.folderRepo.GetAllFoldersWithSongCount(orderAsc)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "GET_FOLDERS_FAILED", "获取文件夹列表失败")
		return
	}

	response.Success(c, folders)
}

// GetFolderSongs - 获取特定文件夹的歌曲列表
// GET /api/folders/:id/songs?sort_by=title&order=asc
func (h *FolderHandler) GetFolderSongs(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_FOLDER_ID", "无效的文件夹ID")
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

	// 通过 ID 获取文件夹路径（ID 是动态分配的，基于文件夹列表顺序）
	folderPath, err := h.folderRepo.GetFolderPathByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "FOLDER_NOT_FOUND", "文件夹不存在")
		return
	}

	songs, err := h.folderRepo.GetSongsByFolder(folderPath, sortBy, order)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "GET_FOLDER_SONGS_FAILED", "获取文件夹歌曲失败")
		return
	}

	response.Success(c, songs)
}
