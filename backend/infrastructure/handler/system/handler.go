package system

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"master-control/model"
)

type useCase interface {
	GetCatalog(ctx context.Context) ([]model.SystemAction, error)
	GetVolume(ctx context.Context) (model.SystemVolume, error)
	GetBrightness(ctx context.Context) (model.SystemBrightness, error)
	SetVolume(ctx context.Context, volume int) error
	SetBrightness(ctx context.Context, brightness float64) error
}

type handler struct {
	useCase useCase
}

func newHandler(useCase useCase) handler {
	return handler{
		useCase: useCase,
	}
}

func (h handler) getCatalog(c *gin.Context) {
	catalog, err := h.useCase.GetCatalog(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, catalog)
}

func (h handler) getVolume(c *gin.Context) {
	volume, err := h.useCase.GetVolume(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, volume)
}

func (h handler) getBrightness(c *gin.Context) {
	brightness, err := h.useCase.GetBrightness(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, brightness)
}

func (h handler) setVolume(c *gin.Context) {
	var req model.SystemVolume
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Error{Error: err.Error()})
		return
	}

	if err := h.useCase.SetVolume(c.Request.Context(), req.Volume); err != nil {
		c.JSON(http.StatusInternalServerError, model.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.Message{Message: "success"})
}

func (h handler) setBrightness(c *gin.Context) {
	var req model.SystemBrightness
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Error{Error: err.Error()})
		return
	}

	if err := h.useCase.SetBrightness(c.Request.Context(), req.Brightness); err != nil {
		c.JSON(http.StatusInternalServerError, model.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.Message{Message: "success"})
}