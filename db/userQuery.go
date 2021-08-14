package db

import (
	"digi-bot/model"
)

func GetUserById(userId int) model.UserModel {
	var user model.UserModel
	DB.Where(model.UserModel{ID: userId}).First(&user)
	return user
}

func GetAllUsersIdByProductId(productId int) []int {
	var usersId []int
	DB.
		Table("pivot_models").
		Where(model.PivotModel{ProductId: productId}).
		Select("user_id").
		Distinct("user_id").
		Find(&usersId)

	return usersId
}
