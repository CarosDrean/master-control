package system

import (
	"github.com/gin-gonic/gin"

	"master-control/model"
)

func NewRouter(api *gin.RouterGroup, useCase model.SystemUseCase) {
	routes := api.Group("system")
	routes.GET("actions", newHandler(useCase).getCatalog)
	routes.GET("volume", newHandler(useCase).getVolume)
	routes.GET("brightness", newHandler(useCase).getBrightness)
	routes.POST("volume", newHandler(useCase).setVolume)
	routes.POST("brightness", newHandler(useCase).setBrightness)
}