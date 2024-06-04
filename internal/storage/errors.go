package storage

import (
	"fmt"

	"github.com/ilya-burinskiy/birthday-notify/internal/models"
)

type ErrUserNotUniq struct {
	User models.User
}

func (err ErrUserNotUniq) Error() string {
	return fmt.Sprintf("user with email \"%s\" already exists", err.User.Email)
}

type ErrUserNotFound struct {
	User models.User
}

func (err ErrUserNotFound) Error() string {
	return fmt.Sprintf("user with email \"%s\" not found", err.User.Email)
}
