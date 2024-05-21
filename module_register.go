package main

import (
	"core/module"

	"github.com/gofiber/fiber/v2"
)

func ModuleRegister(api fiber.Router) {

	Example := module.Example{}
	Example.Route(api)

	// --------------------------
	// --------------------------

	Temp := module.Temp{}
	Temp.Route(api)

	// --------------------------

	Email := module.Email{}
	Email.Route(api)

	WhatsApp := module.WhatsApp{}
	WhatsApp.Route(api)

	Pdf := module.Pdf{}
	Pdf.Route(api)

	// --------------------------

	Auth := module.Auth{}
	Auth.Route(api)

	// --------------------------
	// --------------------------

}
