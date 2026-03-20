package api

import (
	"github.com/gin-gonic/gin"

	"password-timer-manager/internal/api/handler"
	"password-timer-manager/internal/api/middleware"
	"password-timer-manager/internal/model"
)

type Router struct {
	authHandler     *handler.AuthHandler
	packageHandler  *handler.PackageHandler
	passwordHandler *handler.PasswordHandler
	verifyHandler   *handler.VerifyHandler
	auditHandler    *handler.AuditHandler
	adminHandler    *handler.AdminHandler
	authMiddleware  *middleware.AuthMiddleware
}

func NewRouter(
	authHandler *handler.AuthHandler,
	packageHandler *handler.PackageHandler,
	passwordHandler *handler.PasswordHandler,
	verifyHandler *handler.VerifyHandler,
	auditHandler *handler.AuditHandler,
	adminHandler *handler.AdminHandler,
	authMiddleware *middleware.AuthMiddleware,
) *Router {
	return &Router{
		authHandler:     authHandler,
		packageHandler:  packageHandler,
		passwordHandler: passwordHandler,
		verifyHandler:   verifyHandler,
		auditHandler:    auditHandler,
		adminHandler:    adminHandler,
		authMiddleware:  authMiddleware,
	}
}

func (r *Router) SetupRoutes(e *gin.Engine) {
	e.Use(middleware.CORS())

	e.GET("/api/v1/health", r.authHandler.HealthCheck)

	api := e.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", r.authHandler.Login)
			auth.POST("/register", r.authHandler.Register)
			auth.POST("/logout", r.authHandler.Logout)
			auth.GET("/me", r.authMiddleware.JWT(), r.authHandler.Me)
			auth.PUT("/password", r.authMiddleware.JWT(), r.authHandler.ChangePassword)
			auth.POST("/api-keys", r.authMiddleware.JWT(), r.authHandler.CreateAPIKey)
			auth.GET("/api-keys", r.authMiddleware.JWT(), r.authHandler.ListAPIKeys)
			auth.DELETE("/api-keys/:key_id", r.authMiddleware.JWT(), r.authHandler.DeleteAPIKey)
		}

		packages := api.Group("/packages")
		{
			packages.Use(r.authMiddleware.JWT())
			packages.GET("", r.packageHandler.List)
			packages.POST("", middleware.RequirePermission(model.PermissionPackageCreate), r.packageHandler.Create)
			packages.GET("/:id", r.packageHandler.Get)
			packages.PUT("/:id", middleware.RequirePermission(model.PermissionPackageUpdate), r.packageHandler.Update)
			packages.DELETE("/:id", middleware.RequirePermission(model.PermissionPackageDelete), r.packageHandler.Delete)
			packages.GET("/:id/download", r.packageHandler.Download)
			packages.POST("/:id/download-url", r.packageHandler.GetDownloadURL)
			packages.GET("/:id/passwords", r.passwordHandler.List)
			packages.POST("/:id/passwords", middleware.RequirePermission(model.PermissionPasswordManage), r.passwordHandler.Create)
			packages.POST("/:id/passwords/batch", middleware.RequirePermission(model.PermissionPasswordManage), r.passwordHandler.CreateBatch)
			packages.GET("/:id/passwords/current", r.passwordHandler.GetCurrent)
		}

		passwords := api.Group("/passwords")
		passwords.Use(r.authMiddleware.JWT())
		{
			passwords.GET("/:id", r.passwordHandler.Get)
			passwords.PUT("/:id", middleware.RequirePermission(model.PermissionPasswordManage), r.passwordHandler.Update)
			passwords.DELETE("/:id", middleware.RequirePermission(model.PermissionPasswordManage), r.passwordHandler.Delete)
			passwords.POST("/:id/activate", middleware.RequirePermission(model.PermissionPasswordActivate), r.passwordHandler.Activate)
			passwords.POST("/:id/deactivate", middleware.RequirePermission(model.PermissionPasswordActivate), r.passwordHandler.Deactivate)
		}

		verify := api.Group("/verify")
		{
			verify.POST("", r.verifyHandler.Verify)
			verify.POST("/batch", r.verifyHandler.BatchVerify)
			verify.GET("/status/:id", r.verifyHandler.GetStatus)
		}

		audit := api.Group("/audit")
		audit.Use(r.authMiddleware.JWT(), middleware.RequirePermission(model.PermissionAuditRead))
		{
			audit.GET("", r.auditHandler.List)
			audit.GET("/package/:id", r.auditHandler.ListByPackage)
			audit.GET("/user/:id", r.auditHandler.ListByUser)
			audit.GET("/verify-fails", r.auditHandler.ListVerifyFails)
			audit.GET("/export", r.auditHandler.Export)
		}

		admin := api.Group("/admin")
		admin.Use(r.authMiddleware.JWT())
		{
			admin.GET("/encryption/config", r.adminHandler.GetEncryptionConfig)
			admin.PUT("/encryption/config", middleware.RequirePermission(model.PermissionEncryptionManage), r.adminHandler.UpdateEncryptionConfig)
			admin.POST("/encryption/validate", middleware.RequirePermission(model.PermissionEncryptionManage), r.adminHandler.ValidateEncryptionConfig)
			admin.POST("/encryption/regenerate-keys", middleware.RequirePermission(model.PermissionEncryptionManage), r.adminHandler.RegenerateKeys)
			admin.GET("/system", r.adminHandler.GetSystemInfo)
			admin.GET("/roles", r.adminHandler.ListRoles)
			admin.GET("/stats/dashboard", r.adminHandler.GetDashboardStats)
		}
	}
}
