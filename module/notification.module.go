package module

import (
	"net/http"
	"project/middleware"

	"github.com/labstack/echo/v4"
)

type Notification struct{}

func (ref Notification) Route(e *echo.Group) {
	handler := NotificationHandler{}

	e.POST("/:project_key/notification-send", handler.Send, middleware.Onlyproject)

}

// ---------------------------------------------------------------------------------------------
// ---------------------------------------------------------------------------------------------

type NotificationHandler struct{}

func (handler NotificationHandler) Send(c echo.Context) error {
	// var err error

	return c.JSON(http.StatusOK, map[string]string{"message": "OK"})
}
