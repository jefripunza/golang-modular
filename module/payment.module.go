package module

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Payment struct{}

func (ref Payment) Route(e *echo.Group) {
	handler := PaymentHandler{}

	e.GET("/:project_key/example-trigger/:value", handler.Trigger)

}

// ---------------------------------------------------------------------------------------------
// ---------------------------------------------------------------------------------------------

type PaymentHandler struct{}

func (handler PaymentHandler) Trigger(c echo.Context) error {
	// var err error

	value := c.Param("value")

	return c.JSON(http.StatusOK, map[string]string{"message": "OK", "value": value})
}
