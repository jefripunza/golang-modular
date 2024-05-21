package middleware

import (
	"context"
	"core/connection"
	"core/env"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IProject struct {
	ID        primitive.ObjectID `bson:"_id"`
	SecretKey map[string]string  `bson:"secret_key"`
}

func OnIntranetNetwork(c *fiber.Ctx) error {
	var err error

	projectKey := c.Params("project_key")
	secretKey := c.Get("x-secret-key")
	environment := c.Get("x-environment")
	// fmt.Printf("project_key: %s | x-secret-key: %s | x-environment: %s\n", projectKey, secretKey, environment)

	if secretKey == "" {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{
			"message": "secret-key is required!",
		})
	}

	MongoDB := connection.MongoDB{}
	client, ctx, err := MongoDB.Connect()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "internal server error",
		})
	}
	defer client.Disconnect(ctx)
	database := client.Database(env.GetMongoName())
	projectsCollection := database.Collection("projects")

	var core IProject
	err = projectsCollection.FindOne(ctx, bson.M{
		"key":                       projectKey,
		"secret_key." + environment: secretKey,
	}).Decode(&core)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "only core !!!",
		})
	}
	projectID := core.ID.Hex()
	fmt.Println("project_id:", projectID)

	// Set project_id in request context
	ctx = context.WithValue(c.UserContext(), "project_id", projectID)
	c.SetUserContext(ctx)

	return c.Next()
}
