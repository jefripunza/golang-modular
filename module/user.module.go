package module

import (
	"github.com/gofiber/fiber/v2"
)

type User struct{}

func (ref User) Route(api fiber.Router) {
	handler := UserHandler{}
	route := api.Group("/user")

	route.Get("/detail", handler.Detail)

}

// ---------------------------------------------------------------------------------------------
// ---------------------------------------------------------------------------------------------

type UserHandler struct{}

func (handler UserHandler) Detail(c *fiber.Ctx) error {
	// var err error

	return c.Status(fiber.StatusOK).JSON(map[string]string{
		"message": "OK",
	})
}
