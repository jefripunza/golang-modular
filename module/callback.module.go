package module

import (
	"context"
	"core/config"
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/api/oauth2/v1"
)

type Callback struct{}

func (ref Callback) Route(api fiber.Router) {
	handler := CallbackHandler{}
	route := api.Group("/callback")

	route.Get("/google-login", handler.GoogleLogin)

}

// ---------------------------------------------------------------------------------------------
// ---------------------------------------------------------------------------------------------

type CallbackHandler struct{}

func (handler CallbackHandler) GoogleLogin(c *fiber.Ctx) error {
	state := c.Query("state")
	if state != config.GoogleOauthStateString {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "invalid state"})
	}

	code := c.Query("code")
	token, err := config.GoogleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Printf("oauthConf.Exchange() failed with '%s'\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to exchange token"})
	}

	client := config.GoogleOauthConfig.Client(context.Background(), token)
	oauth2Service, err := oauth2.New(client)
	if err != nil {
		log.Printf("oauth2.New() failed with '%s'\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to create oauth2 service"})
	}

	userInfo, err := oauth2Service.Userinfo.Get().Do()
	if err != nil {
		log.Printf("Userinfo.Get() failed with '%s'\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to get user info"})
	}

	userJson, err := json.Marshal(userInfo)
	if err != nil {
		log.Printf("json.Marshal() failed with '%s'\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to marshal user info"})
	}

	return c.Status(fiber.StatusOK).SendString(string(userJson))
}
