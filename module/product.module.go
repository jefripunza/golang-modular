package module

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Product struct{}

func (ref Product) Route(e *echo.Group) {
	handler := ProductHandler{}

	e.GET("/:project_key/example-trigger/:value", handler.Trigger)

}

// ---------------------------------------------------------------------------------------------
// ---------------------------------------------------------------------------------------------

type ProductHandler struct{}

func (handler ProductHandler) Trigger(c echo.Context) error {
	// var err error

	value := c.Param("value")

	return c.JSON(http.StatusOK, map[string]string{"message": "OK", "value": value})
}
