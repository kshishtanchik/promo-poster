package models

import "gorm.io/gorm"

// Действия пользователя на странице
type UserEvent struct {
	gorm.Model
	Member     Member
	MemberID   uint   `json:"memberId" description:"Идентификатор пользователя"`
	ActivityId string `json:"activityId" description:"Идентификатор мероприятия"`
	ActionType string `json:"actionType" description:"Тип события"`
}
