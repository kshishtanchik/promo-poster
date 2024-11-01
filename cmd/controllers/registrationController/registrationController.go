package registrationController

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kshishtanchik/promo-poster/cmd/DB/database"
	"github.com/kshishtanchik/promo-poster/cmd/models"
	"github.com/peteprogrammer/go-automapper"
)

type RegistrationRequest struct {
	ChatId        string `json:"chatId"`
	MemberId      string `json:"memberId"`
	ActivityId    string `json:"activityId"`
	Name          string `json:"memberName"`
	Phone         string `json:"phone"`
	TelegramAlias string `json:"telegramAlias"`
}

func RegeserHandlers(app *fiber.App) {
	reg := app.Group("/registration")
	reg.Post("/:activityId", RegisterMember)
	reg.Get("/:activityId/:memberId/cancel", CancelRegistation)
}

func RegisterMember(c *fiber.Ctx) error {
	activityId := c.Params("activityId")

	registrationRequest := new(RegistrationRequest)
	if err := c.BodyParser(&registrationRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	member := models.Member{}
	database.Db.Db.First(&member, "ID = ?", registrationRequest.MemberId)
	automapper.Map(registrationRequest, member)

	activity := models.Activity{}
	database.Db.Db.First(&activity, "ID = ?", activityId)

	eventRegistration := models.Registration{
		Activity:  activity,
		Member:    member,
		Registred: true,
	}
	database.Db.Db.Create(&eventRegistration)

	return c.Status(200).JSON(activity)
}

func CancelRegistation(c *fiber.Ctx) error {
	activityId := c.Params("activityId")
	memberId := c.Params("memberId")

	eventRegistration := models.Registration{}

	database.Db.Db.Where("ActivityId = ? AND MemberId >= ?", activityId, memberId).Find(&eventRegistration)

	eventRegistration.Registred = false
	database.Db.Db.Save(&eventRegistration)

	return c.Status(200).JSON(eventRegistration)
}
