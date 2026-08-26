package hotkey

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	hotkeycapture "master-control/infrastructure/hotkey"
	"master-control/model"
)

func capture(c *gin.Context) {
	result, err := hotkeycapture.Capture(15 * time.Second)
	if err != nil {
		switch {
		case errors.Is(err, hotkeycapture.ErrNoPermission):
			c.JSON(http.StatusForbidden, model.Error{Error: err.Error()})
		case errors.Is(err, hotkeycapture.ErrTimeout):
			c.JSON(http.StatusGatewayTimeout, model.Error{Error: err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, model.Error{Error: err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, result)
}