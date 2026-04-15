package handler

import (
	"net/http"

	"nas-manager/internal/service"
	"nas-manager/pkg/response"

	"github.com/gin-gonic/gin"
)

// EncryptHandler - 加密处理器
type EncryptHandler struct {
	encryptService *service.EncryptService
}

// NewEncryptHandler - 创建加密处理器
func NewEncryptHandler(encryptService *service.EncryptService) *EncryptHandler {
	return &EncryptHandler{encryptService: encryptService}
}

// SetupPassword - 设置密码
// POST /api/auth/setup
func (h *EncryptHandler) SetupPassword(c *gin.Context) {
	var req service.SetupPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	if err := h.encryptService.SetupPassword(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "SETUP_FAILED", err.Error())
		return
	}

	response.Success(c, gin.H{"success": true})
}

// VerifyPassword - 验证密码
// POST /api/auth/verify
func (h *EncryptHandler) VerifyPassword(c *gin.Context) {
	var req service.VerifyPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	valid, err := h.encryptService.VerifyPassword(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "VERIFY_FAILED", err.Error())
		return
	}

	response.Success(c, gin.H{"valid": valid})
}

// ChangePassword - 修改密码
// POST /api/auth/change
func (h *EncryptHandler) ChangePassword(c *gin.Context) {
	var req service.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	if err := h.encryptService.ChangePassword(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "CHANGE_FAILED", err.Error())
		return
	}

	response.Success(c, gin.H{"success": true})
}
