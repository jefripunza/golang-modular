package module

import (
	"github.com/gofiber/fiber/v2"
)

type Example struct{}

func (ref Example) Route(api fiber.Router) {
	handler := ExampleHandler{}
	route := api.Group("/example")

	route.Get("/trigger/:value", handler.Trigger)

}

// ---------------------------------------------------------------------------------------------
// ---------------------------------------------------------------------------------------------

type ExampleHandler struct{}

func (handler ExampleHandler) Trigger(c *fiber.Ctx) error {
	// var err error

	value := c.Params("value")

	return c.Status(fiber.StatusOK).JSON(map[string]string{
		"message": "OK",
		"value":   value,
	})
}
