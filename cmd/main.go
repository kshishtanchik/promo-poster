package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/kshishtanchik/promo-poster/cmd/DB/database"
	"github.com/kshishtanchik/promo-poster/cmd/controllers"

	promobot "github.com/kshishtanchik/promo-poster/cmd/bot"
)

const (
	FINISH_TAG           = "#Завершено"
	USER_ERROR           = "Что-то пошло не так.. Уже смотрим что."
	EVENT_FINISH_MESSAGE = "Регистрация на мероприятие закончена"
)

func main() {

	godotenv.Load(".env", "secret.env")

	database.Connect()

	promobot.New(os.Getenv("BOT_TOKEN"))

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
		AllowOrigins:     "*",
		AllowCredentials: false,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))

	controllers.RegisterHandlers(app)
	app.Static("/", "../pubic")

	PORT := os.Getenv("PORT")
	log.Printf("Service started on %v", PORT)

	log.Fatal(app.Listen(":" + PORT))
}
