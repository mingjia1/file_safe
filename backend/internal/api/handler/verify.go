package handler

import (
	"github.com/gin-gonic/gin"

	"password-timer-manager/internal/model"
	"password-timer-manager/internal/pkg/response"
)

type VerifyHandler struct{}

func NewVerifyHandler() *VerifyHandler {
	return &VerifyHandler{}
}

func (h *VerifyHandler) Verify(c *gin.Context) {
	var req model.VerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}
	response.Success(c, model.VerifyResponse{})
}

func (h *VerifyHandler) BatchVerify(c *gin.Context) {
	var req model.BatchVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}
	response.Success(c, model.BatchVerifyResponse{})
}

func (h *VerifyHandler) GetStatus(c *gin.Context) {
	response.Success(c, model.PackageStatusResponse{})
}
