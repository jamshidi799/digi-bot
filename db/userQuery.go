package db

import (
	"digi-bot/model"
)

func GetUserById(userId int) model.User {
	var user model.User
	DB.Where(model.User{ID: userId}).First(&user)
	return user
}

func GetAllUsersIdByProductId(productId int) []int {
	var usersId []int
	DB.
		Table("pivots").
		Where(model.Pivot{ProductID: productId}).
		Select("user_id").
		Distinct("user_id").
		Find(&usersId)

	return usersId
}
