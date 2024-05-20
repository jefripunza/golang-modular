package module

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"project/env"
	"project/middleware"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Temp struct{}

func (ref Temp) Route(e *echo.Group) {
	handler := TempHandler{}

	e.GET("/:project_key/temp-clear", handler.Clear, middleware.Onlyproject)
	e.POST("/:project_key/temp-upload-image", handler.UploadImage, middleware.Onlyproject, middleware.OnlyImage)

}

// ---------------------------------------------------------------------------------------------
// ---------------------------------------------------------------------------------------------

type TempHandler struct{}

func (handler TempHandler) Clear(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"message": "temp clear"})
}

func (handler TempHandler) UploadImage(c echo.Context) error {
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

	return c.JSON(http.StatusOK, map[string]string{"message": "uploaded", "file": newFile})
}
