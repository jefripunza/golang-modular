package module

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"project/connection"
	"project/env"
	"project/interfaces"
	"project/middleware"
	"project/util"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/streadway/amqp"
)

type WhatsApp struct{}

func (ref WhatsApp) Route(e *echo.Group) {
	handler := WhatsAppHandler{}

	e.GET("/whatsapp-qr-code", handler.QrCode)
	e.GET("/:project_key/whatsapp-is-register/:target_number", handler.IsRegister)
	e.POST("/:project_key/whatsapp-send", handler.Send, middleware.Onlyproject)

}

// ---------------------------------------------------------------------------------------------
// ---------------------------------------------------------------------------------------------

type WhatsAppHandler struct{}

func (handler WhatsAppHandler) QrCode(c echo.Context) error {
	if connection.WhatsAppQrCode != "" {
		return c.Render(http.StatusOK, "qr-code.html", map[string]interface{}{
			"qr_code_base64_img": connection.WhatsAppQrCode,
		})
	}

	return c.Render(http.StatusOK, "qr-code-connected.html", nil)
}

func (handler WhatsAppHandler) IsRegister(c echo.Context) error {
	var err error

	Validate := util.Validate{}
	target_number := c.Param("target_number")
	targetNumber, err := Validate.NumberOnly(target_number)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "phone number is not number"})
	}
	target_number = string(rune(targetNumber))
	if connection.WhatsAppClient.IsConnected() && connection.WhatsAppClient.IsLoggedIn() {
		is_registered, err := connection.WhatsAppClient.IsOnWhatsApp([]string{target_number})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		}

		return c.JSON(http.StatusOK, map[string]bool{"is_registered": is_registered[0].IsIn})
	}

	return c.JSON(http.StatusBadRequest, map[string]string{"message": "WhatsApp client not connected or logged in"})
}

func (handler WhatsAppHandler) Send(c echo.Context) error {
	var err error

	// projectID := c.Request().Context().Value("project_id")

	String := util.String{}

	WhatsAppClient := connection.WhatsAppClient
	WhatsAppMessage := connection.WhatsAppMessage

	var body struct {
		Type          string                  `json:"type"`
		TargetNumbers []string                `json:"target_numbers"` // support blast
		Message       *string                 `json:"message,omitempty"`
		Data          *map[string]interface{} `json:"data,omitempty"`
		FileName      *string                 `json:"filename,omitempty"`
	}
	if err = c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Format JSON tidak valid"})
	}

	if WhatsAppClient.IsConnected() && WhatsAppClient.IsLoggedIn() {
		rabbit_url := env.GetRabbitUrl()
		Connection, err := amqp.Dial(rabbit_url)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": fmt.Sprintf("Failed to connect to RabbitMQ: %s", err)})
		}
		defer Connection.Close()
		Channel, err := Connection.Channel()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": fmt.Sprintf("Failed to open a channel: %s", err)})
		}
		defer Channel.Close()

		registers, err := WhatsAppClient.IsOnWhatsApp(body.TargetNumbers)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		}

		if !(body.Type == "text" || body.Type == "image" || body.Type == "file") {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "type not found"})
		}

		messageObj := interfaces.IWhatsAppSendQueueRabbitMQ{
			Type: body.Type,
		}

		if body.Message != nil {
			messageObj.Message = body.Message
			if body.Data != nil {
				message := String.ReplaceMessageWithDynamicData(*body.Message, *body.Data)
				messageObj.Message = &message
			}
		}

		if body.Type == "image" || body.Type == "file" {
			if body.FileName == nil {
				return c.JSON(http.StatusBadRequest, map[string]string{"message": "filename is required"})
			}

			tempFile := filepath.Join(env.GetPwd(), "temp", *body.FileName)
			data, err := os.ReadFile(tempFile)
			if err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{"message": "file not found"})
			}

			if body.Type == "image" {
				contentType := http.DetectContentType(data)
				if !strings.HasPrefix(contentType, "image/") {
					return c.JSON(http.StatusBadRequest, map[string]string{"message": "file is not image"})
				}
			}

			messageObj.FileName = body.FileName
			fmt.Println("OK...")
		}
		fmt.Printf("messageObj: %+v\n", messageObj)

		success := make(map[string]string)
		errors := make(map[string]string)

		for _, register := range registers {
			if register.IsIn {
				jid := register.JID.String()
				messageObj.TargetNumber = jid

				jsonMsg, err := json.Marshal(messageObj)
				if err != nil {
					errors[register.Query] = "error serializing message to json"
					continue
				}

				err = Channel.Publish(
					WhatsAppMessage.Exchange, // exchange
					"/",                      // routing key
					false,                    // mandatory
					false,                    // immediate
					amqp.Publishing{
						ContentType: "application/json",
						Body:        jsonMsg,
					},
				)
				if err != nil {
					errors[register.Query] = fmt.Sprintf("Error on send queue: %+v", err)
					continue
				}

				success[register.Query] = "success"
			} else {
				errors[register.Query] = "not registered"
			}
		}

		return c.JSON(http.StatusOK, map[string]any{
			"message": "success send message",
			"success": success,
			"errors":  errors,
		})
	}

	return c.JSON(http.StatusBadRequest, map[string]string{"message": "WhatsApp client not connected or logged in"})
}
