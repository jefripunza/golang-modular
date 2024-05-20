package module

import (
	"fmt"
	"net/http"
	"net/smtp"
	"project/connection"
	"project/env"
	"project/interfaces"
	"project/middleware"
	"project/util"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Email struct{}

func (ref Email) Route(e *echo.Group) {
	handler := EmailHandler{}

	e.POST("/:project_key/email-send/:email", handler.Send, middleware.Onlyproject)

}

// ---------------------------------------------------------------------------------------------
// ---------------------------------------------------------------------------------------------

type EmailHandler struct{}

func (handler EmailHandler) Send(c echo.Context) error {
	var err error

	MongoDB := connection.MongoDB{}

	Validate := util.Validate{}

	projectID := c.Request().Context().Value("project_id")
	emailSender := c.Param("email")

	var body struct {
		To      string `json:"to"`
		Subject string `json:"subject"`
		Body    string `json:"body"`
	}
	if err = c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Format JSON tidak valid"})
	}

	client, ctx, err := MongoDB.Connect()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": fmt.Sprintf("Error connect core database: %v", err)})
	}
	defer client.Disconnect(ctx)
	database := client.Database(env.GetMongoName())
	defer database.Client().Disconnect(ctx)

	var emailCredential interfaces.IEmailCredential
	email_credentials := database.Collection("email_credentials")
	filter := bson.M{"project_id": projectID, "email": emailSender}
	if err := email_credentials.FindOne(ctx, filter).Decode(&emailCredential); err != nil {
		if err == mongo.ErrNoDocuments {
			return c.JSON(http.StatusNotFound, map[string]string{"message": "email sender not found"})
		} else {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": fmt.Sprintf("Error finding document:%v", err)})
		}
	}

	to := []string{body.To}
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s\r\n", body.To, body.Subject, body.Body))

	port, err := Validate.NumberOnly(emailCredential.Port)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	auth := smtp.PlainAuth("", emailCredential.Email, emailCredential.Password, emailCredential.Host)
	err = smtp.SendMail(fmt.Sprintf("%s:%d", emailCredential.Host, port), auth, emailCredential.Email, to, msg)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": fmt.Sprintf("Error sending email: %v", err)})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "success send email"})
}
