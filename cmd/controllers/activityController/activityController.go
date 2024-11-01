package activityController

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/kshishtanchik/promo-poster/cmd/DB/database"
	"github.com/kshishtanchik/promo-poster/cmd/models"
	"gorm.io/gorm"
)

func RegeserHandlers(app *fiber.App) {
	reg := app.Group("/activity")
	reg.Post("/add", AddActivity)
	reg.Get("/:id", Activity)
	reg.Post("/:id", EditActivity)
}

func AddActivity(c *fiber.Ctx) error {
	activity := new(models.Activity)
	if err := c.BodyParser(activity); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	database.Db.Db.Create(&activity)

	return c.Status(200).JSON(activity)
}

func Activity(c *fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	baseModel := gorm.Model{
		ID: uint(id),
	}
	activity := models.Activity{
		Model: baseModel,
	}
	database.Db.Db.First(&activity)
	return c.Status(200).JSON(activity)
}

func EditActivity(c *fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	baseModel := gorm.Model{
		ID: uint(id),
	}
	activity := models.Activity{
		Model: baseModel,
	}
	database.Db.Db.First(&activity)
	return c.Status(200).JSON(activity)
}
