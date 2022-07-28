package db

import (
	"digi-bot/model"
	tb "gopkg.in/telebot.v3"
)

func SaveUser(user *tb.User) {
	userModel := model.ToUser(user)
	database.Create(&userModel)
}

func GetUserById(userId int64) model.User {
	var user model.User
	database.Where(model.User{ID: userId}).First(&user)
	return user
}

func GetAllUsersIdByProductId(productId int, notificationSetting int) []int64 {
	var usersId []int64
	database.
		Table("pivots").
		Where("product_id = ? AND notification_setting <= ?", productId, notificationSetting).
		Select("user_id").
		Distinct("user_id").
		Find(&usersId)

	return usersId
}
