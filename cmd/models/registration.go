package models

import "gorm.io/gorm"

// Модель записи на мероприятие
type Registration struct {
	gorm.Model
	Activity    Activity
	ActivityID  uint
	Member      Member
	MemberID    uint
	Present     bool `json:"present" description:"Отметка присутствия на мероприятии"`
	PaymentMark bool `json:"paymentMark" description:"Отметка об оплате"`
	Registred   bool `json:"registred" description:"Отметка об регистрации"`
}
