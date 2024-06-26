package client

type UserClient interface {
	UserExists(id int) bool
	ProfileExists(id int) bool
	IsTheProfileOwner(userID, profileID int) bool
	GuildJoin(id int) error
	GuildLeave(id int) error
	SendMessage(id int) error
}
