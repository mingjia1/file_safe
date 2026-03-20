package handler

import (
	"github.com/gin-gonic/gin"

	"password-timer-manager/internal/model"
	"password-timer-manager/internal/pkg/response"
)

type AdminHandler struct{}

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

func (h *AdminHandler) GetEncryptionConfig(c *gin.Context) {
	response.Success(c, model.EncryptionConfig{})
}

func (h *AdminHandler) UpdateEncryptionConfig(c *gin.Context) {
	response.Success(c, nil)
}

func (h *AdminHandler) ValidateEncryptionConfig(c *gin.Context) {
	response.Success(c, model.EncryptionValidationResponse{})
}

func (h *AdminHandler) RegenerateKeys(c *gin.Context) {
	response.Success(c, model.RegenerateKeysResponse{})
}

func (h *AdminHandler) GetSystemInfo(c *gin.Context) {
	response.Success(c, gin.H{
		"version": "1.0.0",
		"status":  "healthy",
	})
}

func (h *AdminHandler) ListRoles(c *gin.Context) {
	response.Success(c, []model.UserRole{})
}

func (h *AdminHandler) GetDashboardStats(c *gin.Context) {
	response.Success(c, gin.H{
		"total_packages":      0,
		"active_packages":     0,
		"total_passwords":     0,
		"active_passwords":    0,
		"total_downloads":     0,
		"total_verifies":      0,
		"verify_success_rate": 0.0,
	})
}
