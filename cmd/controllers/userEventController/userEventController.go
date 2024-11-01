package userEventController

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/kshishtanchik/promo-poster/cmd/DB/database"
	"github.com/kshishtanchik/promo-poster/cmd/models"
	"gorm.io/gorm"
)

type UserEventRequest struct {
	ChatId     string `json:"chatId"`
	ActivityId string `json:"activityId"`
	UserId     string `json:"userId"`
	ActionType string `json:"actionType"`
}

func RegeserHandlers(app *fiber.App) {
	app.Post("/userEvent", recordEvent)
}

func recordEvent(c *fiber.Ctx) error {

	userEvent := new(UserEventRequest)
	if err := c.BodyParser(userEvent); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// если нет пользователя - создать болванку для пользователя, id отправить в ответ
	if userEvent.UserId == "" {
		student := models.Member{}
		database.Db.Db.Create(&student)

		userEvent := models.UserEvent{
			Member:     student,
			ActivityId: userEvent.ActivityId,
			ActionType: "initial",
		}
		database.Db.Db.Create(&userEvent)

		return c.Status(200).JSON(student)
	}
	// если есть userId - добавить событие
	id, _ := strconv.ParseUint(userEvent.UserId, 10, 32)
	baseModel := gorm.Model{
		ID: uint(id),
	}
	student := models.Member{
		Model: baseModel,
	}
	database.Db.Db.First(&student)

	newUserEvent := models.UserEvent{
		Member:     student,
		ActivityId: userEvent.ActivityId,
		ActionType: userEvent.ActionType,
	}

	database.Db.Db.Create(&newUserEvent)
	return c.Status(200).JSON(student)
}
