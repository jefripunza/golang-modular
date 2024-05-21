package module

import (
	"core/connection"
	"core/env"
	"core/interfaces"
	"core/middleware"
	"core/util"
	"fmt"
	"net/smtp"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Email struct{}

func (ref Email) Route(api fiber.Router) {
	handler := EmailHandler{}
	route := api.Group("/email")

	route.Post("/send/:email", handler.Send, middleware.OnIntranetNetwork)
}

// ---------------------------------------------------------------------------------------------
// ---------------------------------------------------------------------------------------------

type EmailHandler struct{}

func (handler EmailHandler) Send(c *fiber.Ctx) error {
	var err error

	MongoDB := connection.MongoDB{}
	Validate := util.Validate{}

	projectID := c.Locals("project_id")
	emailSender := c.Params("email")

	var body struct {
		To      string `json:"to"`
		Subject string `json:"subject"`
		Body    string `json:"body"`
	}
	if err = c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Format JSON tidak valid"})
	}

	client, ctx, err := MongoDB.Connect()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": fmt.Sprintf("Error connect core database: %v", err)})
	}
	defer client.Disconnect(ctx)
	database := client.Database(env.GetMongoName())

	var emailCredential interfaces.IEmailCredential
	emailCredentials := database.Collection("email_credentials")
	filter := bson.M{"project_id": projectID, "email": emailSender}
	if err := emailCredentials.FindOne(ctx, filter).Decode(&emailCredential); err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "email sender not found"})
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": fmt.Sprintf("Error finding document: %v", err)})
		}
	}

	to := []string{body.To}
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s\r\n", body.To, body.Subject, body.Body))

	port, err := Validate.NumberOnly(emailCredential.Port)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	auth := smtp.PlainAuth("", emailCredential.Email, emailCredential.Password, emailCredential.Host)
	err = smtp.SendMail(fmt.Sprintf("%s:%d", emailCredential.Host, port), auth, emailCredential.Email, to, msg)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": fmt.Sprintf("Error sending email: %v", err)})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "success send email"})
}
