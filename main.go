package main

import (
	"core/connection"
	"core/env"
	"core/initialize"
	"core/util"
	"core/variable"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c *fiber.Ctx) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	var err error

	Dir := util.Dir{}
	Dir.Make(variable.DatabasePath)
	Dir.Make(variable.TempPath)

	Env := util.Env{}
	Env.Load()
	err = Env.SetTimezone()
	if err != nil {
		log.Fatalf("error on set timezone: %s", err.Error())
		return
	}
	server_name := env.GetServerName()

	MongoDB := connection.MongoDB{}
	MongoDB.Connect()

	// RabbitMQ := connection.RabbitMQ{}
	// RabbitMQ.Connect()

	WhatsApp := connection.WhatsApp{}
	WhatsApp.Connect()

	// ---------------------------------

	initialize.MongoDB()

	// ---------------------------------

	app := fiber.New(fiber.Config{
		//Prefork:               true,
		ServerHeader:          server_name,
		DisableStartupMessage: true,
		CaseSensitive:         true,
		BodyLimit:             10 * 1024 * 1024, // 10 MB / max file size
	})
	// app.Renderer = &Template{
	// 	templates: template.Must(template.ParseGlob("view/*.html")),
	// }
	app.Use(helmet.New())
	app.Use(cors.New(cors.Config{
		AllowMethods:  "GET,PUT,POST,DELETE,OPTIONS",
		ExposeHeaders: "Content-Type,Authorization,Accept",
	}))
	app.Use(requestid.New())
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	api := app.Group("/api")

	ModuleRegister(api)

	api.Use("*", func(c *fiber.Ctx) error {
		fmt.Println("disini...")
		return c.Status(fiber.StatusNotFound).JSON(map[string]string{
			"message": "Endpoint tidak ditemukan!",
		})
	})

	port := env.GetServerPort()
	go func() {
		log.Printf("✅ Server \"%s\" started on port http://localhost:%s\n", server_name, port)
		if err := app.Listen(":" + port); err != nil {
			log.Fatalln("error start server:", err)
		}
	}()

	// Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
	log.Println("🚦 Listen to Ctrl+C ...")
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

}
