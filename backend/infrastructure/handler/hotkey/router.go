package hotkey

import (
	"github.com/gin-gonic/gin"
)

func NewRouter(api *gin.RouterGroup) {
	routes := api.Group("hotkey")
	routes.POST("/capture", capture)
}