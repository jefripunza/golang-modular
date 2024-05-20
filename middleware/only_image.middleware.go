package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func OnlyImage(next echo.HandlerFunc) echo.HandlerFunc {
	// Set maksimum ukuran form untuk mencegah serangan DoS
	MaxMultipartMemory := int64(1024 * 1024 * 10) // MB
	return func(c echo.Context) error {
		ct := c.Request().Header.Get(echo.HeaderContentType)
		if strings.HasPrefix(ct, echo.MIMEApplicationJSON) {
			return next(c)
		}
		if strings.HasPrefix(ct, echo.MIMEMultipartForm) {
			err := c.Request().ParseMultipartForm(MaxMultipartMemory)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			_, fileHeader, err := c.Request().FormFile("image")
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "image harus disematkan")
			}
			contentType := fileHeader.Header.Get("Content-Type")
			if !strings.HasPrefix(contentType, "image/") {
				return echo.NewHTTPError(http.StatusBadRequest, "file harus merupakan gambar")
			}
		}
		return next(c)
	}
}
