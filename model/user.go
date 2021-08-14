package model

import (
	"github.com/jinzhu/gorm"
	tb "gopkg.in/tucnak/telebot.v2"
)

type UserModel struct {
	gorm.Model
	ID              int
	Username        string `gorm:"unique_index"`
	FirstName       string `gorm:"size:64"`
	LastName        string `gorm:"size:64"`
	LanguageCode    string `gorm:"size:64"`
	IsBot           bool
	CanJoinGroups   bool
	CanReadMessages bool
	SupportsInline  bool
}

func (user UserModel) ToTbUser() *tb.User {
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

func ToUserModel(user *tb.User) UserModel {
	return UserModel{
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
