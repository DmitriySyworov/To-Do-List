package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `gorm:"not null"`
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	UserId   uint   `gorm:"unique;not null"`
	Error    string `gorm:"-"`
}
type TempUser struct {
	Name     string
	Email    string
	Password string
	UserId   uint
}
type Session struct {
	SessionId     string
	TemporaryCode uint
}
type TaskForm struct {
	gorm.Model
	Header     string
	Task       string `gorm:"not null"`
	Deadline   time.Time
	StatusDone bool   `gorm:"not null"`
	TaskId     uint   `gorm:"not null"`
	UserId     uint   `gorm:"not null"`
	Error      string `gorm:"-"`
}
