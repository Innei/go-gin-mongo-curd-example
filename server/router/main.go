package router

import (
	"clipboard/controllers"
	"clipboard/middlewares"

	"github.com/gin-gonic/gin"
)

var Clip = controllers.Clip{}
var Auth = controllers.Auth{}

func RegisterRoutes(e *gin.Engine) {
	auth := e.Group("/auth")
	{
		auth.POST("/login", Auth.LoginRoute)
		auth.POST("/register", Auth.RegisterRoute)
		auth.Use(middlewares.PermissionMiddleware()).PATCH("/", Auth.PatchUserRoute)
	}
	clip := e.Group("/clip")
	clip.Use(middlewares.PermissionMiddleware())
	{
		clip.GET("/", Clip.GetClipsRoute)
		clip.GET("/:id", Clip.GetClipOneRoute)
		clip.POST("/", Clip.CreateClipRoute)
		clip.PATCH("/:id", Clip.PatchClipRoute)
		clip.DELETE("/:id", Clip.DeleteClipRoute)
	}
}
