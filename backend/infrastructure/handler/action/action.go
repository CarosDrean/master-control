package action

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"master-control/model"
)

type useCase interface {
	Execute(ctx context.Context, id string) (model.ExecutionResult, error)
	Set(ctx context.Context, action model.Action) error
	GetAll(ctx context.Context) ([]model.Action, error)
	Get(ctx context.Context, id string) (model.Action, error)
	Delete(ctx context.Context, id string) error
}

type handler struct {
	useCase useCase
}

func newHandler(useCase useCase) handler {
	return handler{
		useCase: useCase,
	}
}

func (h handler) execute(c *gin.Context) {
	id := c.Param("id")
	result, err := h.useCase.Execute(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Error{Error: err.Error()})
		return
	}

	// Un exit code distinto de cero no es un error del servidor: viaja en el resultado.
	c.JSON(http.StatusOK, result)
}

func (h handler) set(c *gin.Context) {
	var action model.Action
	if err := c.ShouldBindJSON(&action); err != nil {
		c.JSON(http.StatusBadRequest, model.Error{Error: err.Error()})
		return
	}

	if err := h.useCase.Set(c.Request.Context(), action); err != nil {
		c.JSON(http.StatusInternalServerError, model.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.Message{Message: "success"})
}

func (h handler) getAll(c *gin.Context) {
	actions, err := h.useCase.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, actions)
}

func (h handler) get(c *gin.Context) {
	id := c.Param("id")

	action, err := h.useCase.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, action)
}

func (h handler) delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.useCase.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, model.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.Message{Message: "success"})
}
