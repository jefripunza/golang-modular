package module

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"project/connection"
	"project/env"
	"project/middleware"
	"project/util"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Pdf struct{}

func (ref Pdf) Route(e *echo.Group) {
	handler := PdfHandler{}

	e.POST("/:project_key/pdf-make", handler.Make, middleware.Onlyproject)

}

// ---------------------------------------------------------------------------------------------
// ---------------------------------------------------------------------------------------------

type PdfHandler struct{}

func (handler PdfHandler) Make(c echo.Context) error {
	var err error

	Validate := util.Validate{}
	CDN := connection.CDN{}

	var body struct {
		Html  *string `json:"html"`
		URL   *string `json:"url"`
		Delay *int    `json:"delay"`
	}
	if err = c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Format JSON tidak valid"})
	}

	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": fmt.Sprintf("Gagal membuat PDF generator: %v", err)})
	}
	if body.Html != nil {
		pdfg.AddPage(wkhtmltopdf.NewPageReader(bytes.NewReader([]byte(*body.Html))))
	} else if body.URL != nil {
		page := wkhtmltopdf.NewPage(*body.URL)
		delaySeconds := 5
		if body.Delay != nil {
			delay, err := Validate.NumberOnly(*body.Delay)
			if err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{"message": "delay bukan integer"})
			}
			delaySeconds = delay
		}
		page.JavascriptDelay.Set(uint(delaySeconds * 1000))
		pdfg.AddPage(page)
	} else {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Tidak ada HTML atau URL yang disediakan"})
	}

	err = pdfg.Create()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": fmt.Sprintf("Gagal membuat PDF: %v", err)})
	}

	tempDir := filepath.Join(env.GetPwd(), "temp")
	uuidv4 := uuid.NewString()
	tempPath := filepath.Join(tempDir, uuidv4+".pdf")
	err = ioutil.WriteFile(tempPath, pdfg.Bytes(), 0644)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": fmt.Sprintf("Gagal menyimpan PDF: %v", err)})
	}

	fileName, err := CDN.Upload(tempPath, "pdf")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": fmt.Sprintf("Gagal upload ke CDN: %v", err)})
	}

	err = os.Remove(tempPath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": fmt.Sprintf("Gagal menghapus file PDF: %v", err)})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "PDF berhasil diupload",
		"result":  fileName,
	})

}
