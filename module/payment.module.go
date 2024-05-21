package module

import (
	"github.com/gofiber/fiber/v2"
)

type Payment struct{}

func (ref Payment) Route(api fiber.Router) {
	handler := PaymentHandler{}
	route := api.Group("/payment")

	route.Post("/make", handler.Make)

}

// ---------------------------------------------------------------------------------------------
// ---------------------------------------------------------------------------------------------

type PaymentHandler struct{}

func (handler PaymentHandler) Make(c *fiber.Ctx) error {
	// var err error

	return c.Status(fiber.StatusOK).JSON(map[string]string{
		"message": "OK",
	})
}
