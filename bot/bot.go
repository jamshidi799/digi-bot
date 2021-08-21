package bot

import "sync"

type Bot interface {
	Init(*sync.WaitGroup)
	SendUpdateForUsers(usersId []int, productId int, message string)
	handleStart()
	handleDeleteAll()
	handleAdd()
	handleDelete()
	handleHelp()
	handleList()
	handleGraph()
	handleSetting()
}
