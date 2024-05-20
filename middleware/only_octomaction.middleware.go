package middleware

import (
	"context"
	"fmt"
	"net/http"
	"project/connection"
	"project/env"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IProject struct {
	ID        primitive.ObjectID `bson:"_id"`
	SecretKey map[string]string  `bson:"secret_key"`
}

func Onlyproject(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var err error

		MongoDB := connection.MongoDB{}

		projectKey := c.Param("project_key")
		secretKey := c.Request().Header.Get("x-secret-key")
		environment := c.Request().Header.Get("x-environment")
		// fmt.Printf("project_key: %s | x-secret-key: %s | x-environment: %s\n", projectKey, secretKey, environment)

		if secretKey == "" {
			return c.JSON(http.StatusForbidden, map[string]string{"message": "secret-key is required!"})
		}

		client, ctx, err := MongoDB.Connect()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error"})
		}
		defer client.Disconnect(ctx)
		database := client.Database(env.GetMongoName())
		defer database.Client().Disconnect(ctx)
		projectsCollection := database.Collection("projects")

		var project IProject
		err = projectsCollection.FindOne(ctx, bson.M{
			"key":                       projectKey,
			"secret_key." + environment: secretKey,
		}).Decode(&project)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "only project !!!"})
		}
		project_id := project.ID.Hex()
		fmt.Println("project_id:", project_id)

		// Set project_id in request context
		ctx = context.WithValue(c.Request().Context(), "project_id", project_id)
		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}
