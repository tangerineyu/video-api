package model

import (
	"time"

	"gorm.io/gorm"
)

type Video struct {
	gorm.Model
	UserID        uint   `gorm:"not null; index"`
	Author        User   `gorm:"foreignKey:UserID"`
	Title         string `gorm:"not null; type:varchar(255)"`
	PlayURL       string `gorm:"not null; type:varchar(255)"`
	CoverURL      string `gorm:"not null; type:varchar(255)"`
	FavoriteCount int    `gorm:"default:0"`
	CommentCount  int    `gorm:"default:0"`
	PublishTime   time.Time
	//描述和点击量
	Description string `gorm:"type:text"`
	Views       int    `gorm:"default:0"`
}
