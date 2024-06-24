package client

type UserProfileClient interface {
	ProfileExists(id int) bool
	GuildJoin(id int) error
	GuildLeave(id int) error
	SendMessage(id int) error
}
