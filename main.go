package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"project/connection"
	"project/env"
	"project/util"
	"project/variable"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	// var err error

	Dir := util.Dir{}
	Dir.Make(variable.DatabasePath)
	Dir.Make(variable.TempPath)

	Env := util.Env{}
	Env.Load()

	MongoDB := connection.MongoDB{}
	MongoDB.Connect()

	// RabbitMQ := connection.RabbitMQ{}
	// RabbitMQ.Connect()

	WhatsApp := connection.WhatsApp{}
	WhatsApp.Connect()

	// ---------------------------------

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Renderer = &Template{
		templates: template.Must(template.ParseGlob("view/*.html")),
	}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	api := e.Group("/api")

	ModuleRegister(api)

	e.Any("*", func(c echo.Context) error {
		return c.JSON(http.StatusNotFound, map[string]string{
			"message": "Endpoint tidak ditemukan!",
		})
	})

	port := env.GetPort()
	go func() {
		log.Printf("✅ Server started on port http://localhost:%s\n", port)
		if err := e.Start(":" + port); err != nil {
			e.Logger.Fatal("error start server:", err)
		}
	}()

	// Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
	log.Println("🚦 Listen to Ctrl+C ...")
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

}
