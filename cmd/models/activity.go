package models

import "gorm.io/gorm"

// модель активности
type Activity struct {
	gorm.Model
	Title        string `json:"title" description:"Заголовок события" gorm:"text;default:занятие"`
	Description  string `json:"description" description:"Краткое описание мероприятия" gorm:"text;default:null"`
	StartDate    string `json:"startDate" description:"Дата мероприятия" gorm:"text"`
	ActivityType string `json:"activityType" description:"Тип мероприятия" gorm:"text;not null;default:Мероприятие"`
	Duration     string `json:"duration" description:"Продолжительность" gorm:"text;not null;default:2"`
	Address      string `json:"address" description:"Адресс мероприятия" gorm:"text;not null;default:null"`
	Periodicity  string `json:"periodicity" description:"Периодичность"  gorm:"text;not null;default:7"`
	// todo: подумать в чем. base64 или путь на диске
	Baner  string `json:"baner" gorm:"text"`
	ChatId string `json:"chatId" gorm:"text"`
}
