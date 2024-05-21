package module

import (
	"core/middleware"

	"github.com/gofiber/fiber/v2"
)

type CronJob struct{}

func (ref CronJob) Route(api fiber.Router) {
	handler := CronJobHandler{}
	route := api.Group("/cron-job", middleware.OnIntranetNetwork)

	route.Post("/trigger/:cron_name", handler.Trigger)

}

// ---------------------------------------------------------------------------------------------
// ---------------------------------------------------------------------------------------------

type CronJobHandler struct{}

func (handler CronJobHandler) Trigger(c *fiber.Ctx) error {
	// var err error

	return c.Status(fiber.StatusOK).JSON(map[string]string{
		"message": "OK",
	})
}
