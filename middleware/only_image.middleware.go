package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func OnlyImage(c *fiber.Ctx) error {
	ct := c.Get(fiber.HeaderContentType)
	if strings.HasPrefix(ct, fiber.MIMEApplicationJSON) {
		return c.Next()
	}
	if strings.HasPrefix(ct, fiber.MIMEMultipartForm) {
		file, err := c.FormFile("image")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "image harus disematkan",
			})
		}
		contentType := file.Header.Get("Content-Type")
		if !strings.HasPrefix(contentType, "image/") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "file harus merupakan gambar",
			})
		}
	}
	return c.Next()
}
