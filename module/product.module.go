package module

import (
	"github.com/gofiber/fiber/v2"
)

type Product struct{}

func (ref Product) Route(api fiber.Router) {
	handler := ProductHandler{}
	route := api.Group("/product")

	route.Get("/best-seller", handler.BestSeller)
	route.Get("/detail/:seo_url", handler.Detail)

}

// ---------------------------------------------------------------------------------------------
// ---------------------------------------------------------------------------------------------

type ProductHandler struct{}

func (handler ProductHandler) BestSeller(c *fiber.Ctx) error {
	// var err error

	return c.Status(fiber.StatusOK).JSON(map[string]string{
		"message": "OK",
	})
}

func (handler ProductHandler) Detail(c *fiber.Ctx) error {
	// var err error

	seo_url := c.Params("seo_url")

	return c.Status(fiber.StatusOK).JSON(map[string]string{
		"message": "OK",
		"seo_url": seo_url,
	})
}
