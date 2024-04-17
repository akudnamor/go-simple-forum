package viewmodel

import "go-simple-forum/internal/storage"

type LoggedIn struct {
	Status string
	User   storage.User
}
