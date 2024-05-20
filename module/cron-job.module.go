package module

import (
	"project/middleware"

	"github.com/labstack/echo/v4"
)

type CronJob struct{}

func (ref CronJob) Route(e *echo.Group) {
	handler := CronJobHandler{}

	e.POST("/:project_key/cron-job-trigger/:cron_name", handler.Trigger, middleware.Onlyproject)

}

// ---------------------------------------------------------------------------------------------
// ---------------------------------------------------------------------------------------------

type CronJobHandler struct{}

func (handler CronJobHandler) Trigger(c echo.Context) error {
	// var err error

	return nil
}
