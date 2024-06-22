package client

type UserClient interface {
	UserExists(id int) error
	IncrementUserGuildsQuantity(id int, quantityToInc int) error
	DecrementUserGuildsQuantity(id int, quantityToDec int) error
}
