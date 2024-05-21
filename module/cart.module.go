package module

import (
	"github.com/gofiber/fiber/v2"
)

type Cart struct{}

func (ref Cart) Route(api fiber.Router) {
	handler := CartHandler{}
	route := api.Group("/cart")

	route.Get("/get", handler.Get)

}

// ---------------------------------------------------------------------------------------------
// ---------------------------------------------------------------------------------------------

type CartHandler struct{}

func (handler CartHandler) Get(c *fiber.Ctx) error {
	// var err error

	return c.Status(fiber.StatusOK).JSON(map[string]string{
		"message": "OK",
	})
}
