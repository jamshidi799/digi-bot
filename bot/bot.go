package bot

type Bot interface {
	SendUpdateForUsers(usersId []int, productId int, message string, available bool)
}
