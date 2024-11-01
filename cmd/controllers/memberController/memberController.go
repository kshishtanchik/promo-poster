package memberController

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/kshishtanchik/promo-poster/cmd/DB/database"
	"github.com/kshishtanchik/promo-poster/cmd/models"
	"gorm.io/gorm"
)

func RegeserHandlers(app *fiber.App) {
	app.Get("/member/:id", memberById)
}

func memberById(c *fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	baseModel := gorm.Model{
		ID: uint(id),
	}
	member := models.Member{
		Model: baseModel,
	}

	database.Db.Db.First(&member)

	return c.Status(200).JSON(member)

}
