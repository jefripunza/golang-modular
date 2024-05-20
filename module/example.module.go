package module

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Example struct{}

func (ref Example) Route(e *echo.Group) {
	handler := ExampleHandler{}

	e.GET("/:project_key/example-trigger/:value", handler.Trigger)

}

// ---------------------------------------------------------------------------------------------
// ---------------------------------------------------------------------------------------------

type ExampleHandler struct{}

func (handler ExampleHandler) Trigger(c echo.Context) error {
	// var err error

	value := c.Param("value")

	return c.JSON(http.StatusOK, map[string]string{"message": "OK", "value": value})
}
