package handler

import (
	"github.com/gin-gonic/gin"

	"password-timer-manager/internal/model"
	"password-timer-manager/internal/pkg/response"
)

type AuditHandler struct{}

func NewAuditHandler() *AuditHandler {
	return &AuditHandler{}
}

func (h *AuditHandler) List(c *gin.Context) {
	response.Success(c, []model.AuditLogResponse{})
}

func (h *AuditHandler) ListByPackage(c *gin.Context) {
	response.Success(c, []model.AuditLogResponse{})
}

func (h *AuditHandler) ListByUser(c *gin.Context) {
	response.Success(c, []model.AuditLogResponse{})
}

func (h *AuditHandler) ListVerifyFails(c *gin.Context) {
	response.Success(c, []model.AuditLogResponse{})
}

func (h *AuditHandler) Export(c *gin.Context) {
	response.Success(c, nil)
}
