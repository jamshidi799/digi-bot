package model

import (
	"github.com/jinzhu/gorm"
	tb "gopkg.in/telebot.v3"
)

type User struct {
	gorm.Model
	ID              int64
	Username        string `gorm:"unique_index"`
	FirstName       string `gorm:"size:64"`
	LastName        string `gorm:"size:64"`
	LanguageCode    string `gorm:"size:64"`
	IsBot           bool
	CanJoinGroups   bool
	CanReadMessages bool
	SupportsInline  bool
	Pivots          []Pivot
}

func (user User) ToTbUser() *tb.User {
	return &tb.User{
		ID:              user.ID,
		Username:        user.Username,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		LanguageCode:    user.LanguageCode,
		IsBot:           user.IsBot,
		CanJoinGroups:   user.CanJoinGroups,
		CanReadMessages: user.CanReadMessages,
		SupportsInline:  user.SupportsInline,
	}
}

func ToUser(user *tb.User) User {
	return User{
		ID:              user.ID,
		Username:        user.Username,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		LanguageCode:    user.LanguageCode,
		IsBot:           user.IsBot,
		CanJoinGroups:   user.CanJoinGroups,
		CanReadMessages: user.CanReadMessages,
		SupportsInline:  user.SupportsInline,
	}
}
