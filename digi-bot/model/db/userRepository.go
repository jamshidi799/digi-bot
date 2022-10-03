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

func GetAllPivotsByProductId(productId int) []model.UserIdAndDiscountDto {
	var pivots []model.UserIdAndDiscountDto

	database.
		Table("pivots").
		Where("product_id = ?", productId).
		Select("user_id", "discount").
		Distinct("user_id").
		Find(&pivots)

	return pivots
}
