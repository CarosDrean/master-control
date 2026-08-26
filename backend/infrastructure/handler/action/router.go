package action

import (
	"github.com/gin-gonic/gin"

	"master-control/model"
)

func NewRouter(specification model.RouterSpecification) {
	publicRoutes(specification.Api, newHandler(specification.Dependencies.ActionUseCase))
}

func publicRoutes(api *gin.RouterGroup, h handler) {
	routes := api.Group("actions")
	routes.GET("", h.getAll)
	routes.GET("/:id", h.get)
	routes.POST("", h.set)
	routes.PUT("/:id", h.set)
	routes.DELETE("/:id", h.delete)
	routes.POST("/:id/execute", h.execute)
}
