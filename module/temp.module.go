package module

import (
	"core/env"
	"core/middleware"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Temp struct{}

func (ref Temp) Route(api fiber.Router) {
	handler := TempHandler{}
	route := api.Group("/temp")

	route.Get("/clear", handler.Clear, middleware.OnIntranetNetwork)
	route.Post("/upload-image", handler.UploadImage, middleware.OnIntranetNetwork, middleware.OnlyImage)

}

// ---------------------------------------------------------------------------------------------
// ---------------------------------------------------------------------------------------------

type TempHandler struct{}

func (handler TempHandler) Clear(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(map[string]string{
		"message": "temp clear",
	})
}

func (handler TempHandler) UploadImage(c *fiber.Ctx) error {
	var err error

	// Menerima file dari form
	image, err := c.FormFile("image")
	if err != nil {
		return err
	}

	// Membuka file yang diupload
	src, err := image.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Mendapatkan ekstensi dari nama file asli
	ext := filepath.Ext(image.Filename)

	tempDir := filepath.Join(env.GetPwd(), "temp")
	uuidv4 := uuid.NewString()
	newFile := uuidv4 + ext
	tempPath := filepath.Join(tempDir, newFile)
	fmt.Println("tempPath:", tempPath)
	dst, err := os.Create(tempPath)
	if err != nil {
		return err
	}
	defer dst.Close()
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(map[string]string{
		"message": "uploaded",
		"file":    newFile,
	})
}
