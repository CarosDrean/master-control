package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"master-control/infrastructure/handler/action"
	"master-control/infrastructure/handler/actiontype"
	appHandler "master-control/infrastructure/handler/app"
	"master-control/infrastructure/handler/hotkey"
	"master-control/infrastructure/handler/page"
	"master-control/model"
)

func InitRoutes(specification model.RouterSpecification) {
	healthRoutes(specification.Api)
	action.NewRouter(specification)
	actiontype.NewRouter(specification)
	appHandler.NewRouter(specification.Api)
	hotkey.NewRouter(specification.Api)
	page.NewRouter(specification)
}

func healthRoutes(api *gin.RouterGroup) {
	routes := api.Group("")

	routes.GET("health", health)
}

func health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}
