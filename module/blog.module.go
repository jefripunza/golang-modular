package module

import (
	"github.com/gofiber/fiber/v2"
)

type Blog struct{}

func (ref Blog) Route(api fiber.Router) {
	handler := BlogHandler{}
	route := api.Group("/blog")

	route.Get("/:seo_url", handler.Get)

}

// ---------------------------------------------------------------------------------------------
// ---------------------------------------------------------------------------------------------

type BlogHandler struct{}

func (handler BlogHandler) Get(c *fiber.Ctx) error {
	// var err error

	seo_url := c.Params("seo_url")

	return c.Status(fiber.StatusOK).JSON(map[string]string{
		"message": "OK",
		"seo_url": seo_url,
	})
}
