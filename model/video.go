package model

import "time"

type Video struct {
	UserID        uint   `gorm:"not null; index"`
	Author        User   `gorm:"foreignKey:UserID"`
	Title         string `gorm:"not null; type:varchar(255)"`
	PlayURL       string `gorm:"not null; type:varchar(255)"`
	CoverURL      string `gorm:"not null; type:varchar(255)"`
	FavoriteCount int    `gorm:"default:0"`
	CommentCount  int    `gorm:"default:0"`
	PublishTime   time.Time
}
