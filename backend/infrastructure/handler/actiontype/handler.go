package actiontype

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"master-control/model"
)

type UseCase interface {
	Set(ctx context.Context, actionType model.ActionType) error
	GetAll(ctx context.Context) ([]model.ActionType, error)
	Get(ctx context.Context, ID string) (model.ActionType, error)
	Delete(ctx context.Context, ID string) error
}

type handler struct {
	useCase UseCase
}

func newHandler(useCase UseCase) handler {
	return handler{
		useCase: useCase,
	}
}

func (h handler) set(c *gin.Context) {
	var actionType model.ActionType
	if err := c.ShouldBindJSON(&actionType); err != nil {
		c.JSON(http.StatusBadRequest, model.Error{Error: err.Error()})
		return
	}

	if err := h.useCase.Set(c.Request.Context(), actionType); err != nil {
		c.JSON(http.StatusInternalServerError, model.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.Message{Message: "success"})
}

func (h handler) getAll(c *gin.Context) {
	actionTypes, err := h.useCase.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, actionTypes)
}

func (h handler) get(c *gin.Context) {
	ID := c.Param("id")

	actionType, err := h.useCase.Get(c.Request.Context(), ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, actionType)
}

func (h handler) delete(c *gin.Context) {
	ID := c.Param("id")
	if err := h.useCase.Delete(c.Request.Context(), ID); err != nil {
		c.JSON(http.StatusInternalServerError, model.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.Message{Message: "success"})
}
