package handler

import (
	"net/http"

	"nas-manager/internal/service"
	"nas-manager/pkg/response"

	"github.com/gin-gonic/gin"
)

// SettingHandler - 设置处理器
type SettingHandler struct {
	settingService *service.SettingService
}

// NewSettingHandler - 创建设置处理器
func NewSettingHandler(settingService *service.SettingService) *SettingHandler {
	return &SettingHandler{settingService: settingService}
}

// GetSetupStatus - 获取设置状态
// GET /api/setup/status
func (h *SettingHandler) GetSetupStatus(c *gin.Context) {
	status, err := h.settingService.CheckSetupRequired()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "FAILED_TO_CHECK_SETUP", err.Error())
		return
	}
	response.Success(c, status)
}

// SaveSetup - 保存配置
// POST /api/setup
func (h *SettingHandler) SaveSetup(c *gin.Context) {
	var cfg service.SetupConfig
	if err := c.ShouldBindJSON(&cfg); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	if err := h.settingService.SaveSetupConfig(&cfg); err != nil {
		response.Error(c, http.StatusBadRequest, "SETUP_FAILED", err.Error())
		return
	}

	response.Success(c, gin.H{"success": true})
}
