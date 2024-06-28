package repository

type Auth interface{}

type Profile interface{}

type Playlist interface{}

type Chat interface{}

type Song interface{}

type Message interface{}

type Post interface{}

type Repository struct {
	Auth
	Profile
	Playlist
	Chat
	Song
	Message
	Post
}

func NewRepository() *Repository {
	return &Repository{}
}
