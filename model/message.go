package model

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	FromUserID uint   `gorm:"not null; index "`
	ToUserID   uint   `gorm:"not null; index "`
	Content    string `gorm:"type:text; not null"`
}
