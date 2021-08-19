package pivot

import (
	"digi-bot/db"
	"digi-bot/messageCreator"
	"log"
)

func UpdateStatus(status int, productId string, userId int) string {
	result := db.DB.
		Table("pivots").
		Where("user_id = ? AND product_id = ?", userId, productId).
		Update("notification_setting", status)
	if result.Error != nil {
		log.Println(result.Error)
		return messageCreator.CreateErrorText()
	}

	return messageCreator.CreateSuccessText()
}
