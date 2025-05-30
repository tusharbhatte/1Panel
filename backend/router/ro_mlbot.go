package router

import (
	v1 "github.com/1Panel-dev/1Panel/backend/app/api/v1"
	"github.com/1Panel-dev/1Panel/backend/middleware"
	"github.com/gin-gonic/gin"
)

type MLBotRouter struct{}

func (s *MLBotRouter) InitRouter(Router *gin.RouterGroup) {
	mlbotRouter := Router.Group("mlbot").
		Use(middleware.JwtAuth()).
		Use(middleware.SessionAuth()).
		Use(middleware.PasswordExpired())
	baseApi := v1.ApiGroupApp.BaseApi
	{
		mlbotRouter.GET("/config", baseApi.GetMLBotConfig)
		mlbotRouter.POST("/config", baseApi.UpdateMLBotConfig)
		mlbotRouter.POST("/test", baseApi.TestMLBotNotification)
	}
}
