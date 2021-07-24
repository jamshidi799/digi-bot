package db

import "digi-bot/model"

func GetUserById(userId int) model.UserModel {
	var user model.UserModel
	DB.Where(model.UserModel{ID: userId}).First(&user)
	return user
}
