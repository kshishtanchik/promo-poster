package database

import (
	"log"
	"os"

	"github.com/kshishtanchik/promo-poster/cmd/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DbInstance struct {
	Db *gorm.DB
}

var Db DbInstance

func Connect() {
	//host := os.Getenv("DB_HOST")
	//user := os.Getenv("DB_USER")
	//password := os.Getenv("DB_USER_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	//port := os.Getenv("DB_PORT")

	//dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", host, user, password, dbName, port)

	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Не удалось подключиться к базе\n", err)
		os.Exit(2)
	}

	log.Println("Подключенно")
	db.Logger = logger.Default.LogMode(logger.Info)

	log.Println("Запуск автоматических миграций")
	db.AutoMigrate(&models.Activity{})
	db.AutoMigrate(&models.Registration{})
	db.AutoMigrate(&models.Member{})
	db.AutoMigrate(&models.UserEvent{})

	Db = DbInstance{
		Db: db,
	}
}
