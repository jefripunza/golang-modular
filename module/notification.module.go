package module

import (
	"core/middleware"

	"github.com/gofiber/fiber/v2"
)

type Notification struct{}

func (ref Notification) Route(api fiber.Router) {
	handler := NotificationHandler{}
	route := api.Group("/notification")

	route.Post("/send", handler.Send, middleware.OnIntranetNetwork)
	route.Post("/blast", handler.Blast, middleware.OnIntranetNetwork)

}

// ---------------------------------------------------------------------------------------------
// ---------------------------------------------------------------------------------------------

type NotificationHandler struct{}

func (handler NotificationHandler) Send(c *fiber.Ctx) error {
	// var err error

	return c.Status(fiber.StatusOK).JSON(map[string]string{
		"message": "OK",
	})
}

func (handler NotificationHandler) Blast(c *fiber.Ctx) error {
	// var err error

	return c.Status(fiber.StatusOK).JSON(map[string]string{
		"message": "OK",
	})
}
