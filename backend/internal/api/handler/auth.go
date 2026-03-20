package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"password-timer-manager/internal/model"
	ptmerrors "password-timer-manager/internal/pkg/errors"
	"password-timer-manager/internal/pkg/response"
	"password-timer-manager/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	result, err := h.authService.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		if ptmerrors.Is(err, ptmerrors.ErrAuthFailed) || ptmerrors.Is(err, ptmerrors.ErrPasswordWrong) {
			response.Unauthorized(c, "invalid username or password")
			return
		}
		response.Error(c, ptmerrors.ErrInternalError)
		return
	}

	response.Success(c, result)
}

func (h *AuthHandler) Me(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		response.Unauthorized(c, "not authenticated")
		return
	}

	u := user.(*model.User)
	resp := u.ToResponse()
	response.Success(c, resp)
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		response.Unauthorized(c, "not authenticated")
		return
	}

	u := user.(*model.User)

	var req model.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	err := h.authService.ChangePassword(c.Request.Context(), u.ID, req.OldPassword, req.NewPassword)
	if err != nil {
		if ptmerrors.Is(err, ptmerrors.ErrPasswordWrong) {
			response.Unauthorized(c, "current password is incorrect")
			return
		}
		response.Error(c, ptmerrors.ErrInternalError)
		return
	}

	response.Success(c, nil)
}

func (h *AuthHandler) CreateAPIKey(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		response.Unauthorized(c, "not authenticated")
		return
	}

	u := user.(*model.User)

	var req model.CreateApiKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	result, err := h.authService.CreateAPIKey(c.Request.Context(), u.ID, req.Name, req.ExpiresIn)
	if err != nil {
		response.Error(c, ptmerrors.ErrInternalError)
		return
	}

	response.Created(c, result)
}

func (h *AuthHandler) ListAPIKeys(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		response.Unauthorized(c, "not authenticated")
		return
	}

	u := user.(*model.User)

	keys, err := h.authService.ListAPIKeys(c.Request.Context(), u.ID)
	if err != nil {
		response.Error(c, ptmerrors.ErrInternalError)
		return
	}

	response.Success(c, keys)
}

func (h *AuthHandler) DeleteAPIKey(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		response.Unauthorized(c, "not authenticated")
		return
	}

	u := user.(*model.User)
	keyIDStr := c.Param("key_id")
	keyID, err := uuid.Parse(keyIDStr)
	if err != nil {
		response.BadRequest(c, "invalid key id")
		return
	}

	if err := h.authService.DeleteAPIKey(c.Request.Context(), u.ID, keyID); err != nil {
		response.Error(c, ptmerrors.ErrInternalError)
		return
	}

	response.SuccessNoContent(c)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req model.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	result, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		if ptmerrors.Is(err, ptmerrors.ErrUserExists) {
			response.Error(c, ptmerrors.ErrUserExists)
			return
		}
		response.Error(c, ptmerrors.ErrInternalError)
		return
	}

	response.Created(c, result)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	response.Success(c, nil)
}

func (h *AuthHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
	})
}
