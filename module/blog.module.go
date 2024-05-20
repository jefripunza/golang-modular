package module

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Blog struct{}

func (ref Blog) Route(e *echo.Group) {
	handler := BlogHandler{}

	e.GET("/:project_key/example-trigger/:value", handler.Trigger)

}

// ---------------------------------------------------------------------------------------------
// ---------------------------------------------------------------------------------------------

type BlogHandler struct{}

func (handler BlogHandler) Trigger(c echo.Context) error {
	// var err error

	value := c.Param("value")

	return c.JSON(http.StatusOK, map[string]string{"message": "OK", "value": value})
}
