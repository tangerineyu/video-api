package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username      string `gorm:"unique;not null"`
	Password      string `gorm:"not null"`
	Name          string `gorm:"default:'新用户'"`
	Avatar        string
	FollowCount   int64 `gorm:"default:0"`
	FollowerCount int64 `gorm:"default:0"`
	//MFA
	MfaSecret  string `gorm:"default:''"`
	MfaEnabled bool   `gorm:"default:false"`
}
