package models

import "github.com/jinzhu/gorm"

type UserFriend struct {
	gorm.Model
	UserId     uint
	UserNote   string
	UserStauts string
}
