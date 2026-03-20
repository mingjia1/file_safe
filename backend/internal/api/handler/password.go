package handler

import (
	"github.com/gin-gonic/gin"

	"password-timer-manager/internal/model"
	"password-timer-manager/internal/pkg/response"
)

type PasswordHandler struct{}

func NewPasswordHandler() *PasswordHandler {
	return &PasswordHandler{}
}

func (h *PasswordHandler) List(c *gin.Context) {
	response.Success(c, []model.PasswordResponse{})
}

func (h *PasswordHandler) Create(c *gin.Context) {
	response.Created(c, nil)
}

func (h *PasswordHandler) CreateBatch(c *gin.Context) {
	response.Created(c, nil)
}

func (h *PasswordHandler) Get(c *gin.Context) {
	response.Success(c, nil)
}

func (h *PasswordHandler) Update(c *gin.Context) {
	response.Success(c, nil)
}

func (h *PasswordHandler) Delete(c *gin.Context) {
	response.SuccessNoContent(c)
}

func (h *PasswordHandler) Activate(c *gin.Context) {
	response.Success(c, nil)
}

func (h *PasswordHandler) Deactivate(c *gin.Context) {
	response.Success(c, nil)
}

func (h *PasswordHandler) GetCurrent(c *gin.Context) {
	response.Success(c, nil)
}
