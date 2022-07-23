package db

import (
	"digi-bot/service"
	"log"
)

func UpdateStatus(status int, productId string, userId int) string {
	result := database.
		Table("pivots").
		Where("user_id = ? AND product_id = ?", userId, productId).
		Update("notification_setting", status)
	if result.Error != nil {
		log.Println(result.Error)
		return service.CreateErrorText()
	}

	return service.CreateSuccessText()
}
