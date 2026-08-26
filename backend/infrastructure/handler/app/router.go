package app

import (
	"github.com/gin-gonic/gin"
)

func NewRouter(api *gin.RouterGroup) {
	routes := api.Group("apps")
	routes.GET("", getAll)
	routes.GET("/:name/icon", getIcon)
	routes.GET("/icon", getIconByPath)
}