package handler

import (
	"github.com/gin-gonic/gin"

	"password-timer-manager/internal/model"
	"password-timer-manager/internal/pkg/response"
)

type PackageHandler struct{}

func NewPackageHandler() *PackageHandler {
	return &PackageHandler{}
}

func (h *PackageHandler) List(c *gin.Context) {
	response.Success(c, []model.PackageResponse{})
}

func (h *PackageHandler) Create(c *gin.Context) {
	response.Created(c, nil)
}

func (h *PackageHandler) Get(c *gin.Context) {
	response.Success(c, nil)
}

func (h *PackageHandler) Update(c *gin.Context) {
	response.Success(c, nil)
}

func (h *PackageHandler) Delete(c *gin.Context) {
	response.SuccessNoContent(c)
}

func (h *PackageHandler) Download(c *gin.Context) {
	response.Success(c, nil)
}

func (h *PackageHandler) GetDownloadURL(c *gin.Context) {
	response.Success(c, nil)
}
