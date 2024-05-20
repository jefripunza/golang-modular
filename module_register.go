package main

import (
	"project/module"

	"github.com/labstack/echo/v4"
)

func ModuleRegister(e *echo.Group) {

	Example := module.Example{}
	Example.Route(e)

	// --------------------------
	// --------------------------

	Temp := module.Temp{}
	Temp.Route(e)

	// --------------------------

	Email := module.Email{}
	Email.Route(e)

	WhatsApp := module.WhatsApp{}
	WhatsApp.Route(e)

	Pdf := module.Pdf{}
	Pdf.Route(e)

	// --------------------------

	// --------------------------
	// --------------------------

}
