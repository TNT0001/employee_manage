package router

import (
	"github.com/gin-gonic/gin"
	"tungnt/emmployee_manage/internal/service/main_service/pkg/handler"
	"tungnt/emmployee_manage/internal/service/main_service/pkg/middleware"
	"tungnt/emmployee_manage/pkg/share/keycloak"
)

func Mount(router *gin.Engine, client *keycloak.Client, h *handler.Handler) {
	API := router.Group("/api/v1")

	{
		API.GET("/users", h.GetUser)
		API.POST("/:realm/admin_user_token", h.PostAdminUserLogin)
		API.POST("/:realm/user_login", h.PostUserLogin)
		API.POST("/:realm/user_login_otp", h.PostUserLoginOTP)
		API.GET("/teams", h.GetTeams)
		API.GET("/roles", h.GetRoles)
		API.GET("/permissions", h.GetPermissions)

		API.Use(middleware.NewAuthMiddleware(client))
		API.POST("/team", h.PostCreateTeam)
		API.POST("/permissions", h.PostCreatePermissions)
		API.POST("/role", h.PostCreateRole)
		API.PATCH("/role/:role_id", h.PatchUpdateRole)
		API.DELETE("/role/:role_id", h.DeleteRole)
		API.POST("/user", h.PostCreateUser)
	}
}
