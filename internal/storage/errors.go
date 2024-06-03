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
