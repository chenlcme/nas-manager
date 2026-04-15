package handler

import (
	"net/http"
	"os"
	"time"

	"nas-manager/internal/repository"
	"nas-manager/internal/service"
	"nas-manager/pkg/response"

	"github.com/gin-gonic/gin"
)

// ScanHandler - 扫描处理器
type ScanHandler struct {
	scannerService *service.ScannerService
	songRepo       *repository.SongRepository
	settingRepo    *repository.SettingRepository
}

// NewScanHandler - 创建扫描处理器
func NewScanHandler(scannerService *service.ScannerService, songRepo *repository.SongRepository, settingRepo *repository.SettingRepository) *ScanHandler {
	return &ScanHandler{
		scannerService: scannerService,
		songRepo:       songRepo,
		settingRepo:    settingRepo,
	}
}

// ScanRequest - 扫描请求
type ScanRequest struct {
	Mode string `json:"mode"` // "full" or "incremental", defaults to "incremental"
}

// Scan - 触发扫描
// POST /api/songs/scan
func (h *ScanHandler) Scan(c *gin.Context) {
	// 解析请求
	var req ScanRequest
	req.Mode = "incremental" // 默认增量扫描
	if err := c.ShouldBindJSON(&req); err != nil {
		// 忽略错误，使用默认增量模式
	}

	// 获取音乐目录
	musicDir, err := h.settingRepo.GetMusicDir()
	if err != nil {
		response.Error(c, http.StatusBadRequest, "NO_MUSIC_DIR", "Music directory not configured")
		return
	}

	if musicDir == "" {
		response.Error(c, http.StatusBadRequest, "NO_MUSIC_DIR", "Music directory not configured")
		return
	}

	// 检查目录是否存在
	if _, err := os.Stat(musicDir); os.IsNotExist(err) {
		response.Error(c, http.StatusBadRequest, "DIR_NOT_EXIST", "Music directory does not exist")
		return
	}

	// 确定扫描模式
	var mode service.ScanMode
	if req.Mode == "full" {
		mode = service.ScanModeFull
	} else {
		mode = service.ScanModeIncremental
		// 获取上次扫描时间
		lastScanTime, err := h.settingRepo.GetLastScanTime()
		if err == nil {
			h.scannerService.SetLastScanTime(lastScanTime)
		}
	}

	// 执行扫描
	result, err := h.scannerService.ScanFiles(musicDir, mode)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "SCAN_FAILED", err.Error())
		return
	}

	// 更新最后扫描时间
	h.settingRepo.SetLastScanTime(time.Now().Unix())

	response.Success(c, result)
}

// Cleanup - 清理孤立的数据库记录
// POST /api/songs/cleanup
func (h *ScanHandler) Cleanup(c *gin.Context) {
	result, err := h.scannerService.CleanupOrphanRecords()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "CLEANUP_FAILED", err.Error())
		return
	}

	response.Success(c, result)
}
