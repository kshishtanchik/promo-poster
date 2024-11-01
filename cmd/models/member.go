package models

import "gorm.io/gorm"

type Member struct {
	gorm.Model
	Name          string `json:"memberName"`
	Phone         string `json:"phone"`
	TelegramAlias string `json:"telegramAlias"`
}
