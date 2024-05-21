package module

import (
	"bytes"
	"core/connection"
	"core/env"
	"core/middleware"
	"core/util"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Pdf struct{}

func (ref Pdf) Route(api fiber.Router) {
	handler := PdfHandler{}
	route := api.Group("/pdf")

	route.Post("/make", handler.Make, middleware.OnIntranetNetwork)

}

// ---------------------------------------------------------------------------------------------
// ---------------------------------------------------------------------------------------------

type PdfHandler struct{}

func (handler PdfHandler) Make(c *fiber.Ctx) error {
	var err error

	Validate := util.Validate{}
	CDN := connection.CDN{}

	var body struct {
		Html  *string `json:"html"`
		URL   *string `json:"url"`
		Delay *int    `json:"delay"`
	}
	if err = c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]string{
			"message": "Format JSON tidak valid",
		})
	}

	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]string{"message": fmt.Sprintf("Gagal membuat PDF generator: %v", err)})
	}
	if body.Html != nil {
		pdfg.AddPage(wkhtmltopdf.NewPageReader(bytes.NewReader([]byte(*body.Html))))
	} else if body.URL != nil {
		page := wkhtmltopdf.NewPage(*body.URL)
		delaySeconds := 5
		if body.Delay != nil {
			delay, err := Validate.NumberOnly(*body.Delay)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(map[string]string{"message": "delay bukan integer"})
			}
			delaySeconds = delay
		}
		page.JavascriptDelay.Set(uint(delaySeconds * 1000))
		pdfg.AddPage(page)
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]string{"message": "Tidak ada HTML atau URL yang disediakan"})
	}

	err = pdfg.Create()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]string{"message": fmt.Sprintf("Gagal membuat PDF: %v", err)})
	}

	tempDir := filepath.Join(env.GetPwd(), "temp")
	uuidv4 := uuid.NewString()
	tempPath := filepath.Join(tempDir, uuidv4+".pdf")
	err = ioutil.WriteFile(tempPath, pdfg.Bytes(), 0644)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]string{"message": fmt.Sprintf("Gagal menyimpan PDF: %v", err)})
	}

	fileName, err := CDN.Upload(tempPath, "pdf")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]string{"message": fmt.Sprintf("Gagal upload ke CDN: %v", err)})
	}

	err = os.Remove(tempPath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]string{"message": fmt.Sprintf("Gagal menghapus file PDF: %v", err)})
	}

	return c.Status(fiber.StatusOK).JSON(map[string]string{
		"message": "PDF berhasil diupload",
		"result":  fileName,
	})

}
