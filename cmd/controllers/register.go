package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kshishtanchik/promo-poster/cmd/controllers/activityController"
	"github.com/kshishtanchik/promo-poster/cmd/controllers/memberController"
	"github.com/kshishtanchik/promo-poster/cmd/controllers/registrationController"
	"github.com/kshishtanchik/promo-poster/cmd/controllers/userEventController"
)

func RegisterHandlers(app *fiber.App) {
	activityController.RegeserHandlers(app)
	userEventController.RegeserHandlers(app)
	registrationController.RegeserHandlers(app)
	memberController.RegeserHandlers(app)
}
