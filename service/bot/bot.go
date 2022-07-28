package bot

import (
	"sync"
)

type Bot interface {
	SendUpdateForUsers(usersId []int, productId int, message string, available bool)
}

func StartBot(group *sync.WaitGroup) {
	InitTelegramBot(group)
}
