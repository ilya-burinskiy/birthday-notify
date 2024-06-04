package services

import (
	"context"
	"fmt"
	"time"

	"github.com/ilya-burinskiy/birthday-notify/internal/auth"
	"github.com/ilya-burinskiy/birthday-notify/internal/models"
)

type UserCreator interface {
	CreateUser(
		ctx context.Context,
		login string,
		encryptedPassword []byte,
		birthdayDate time.Time,
	) (models.User, error)
}

type RegisterService struct {
	usrCreator UserCreator
}

func NewRegisterService(usrCreator UserCreator) RegisterService {
	return RegisterService{
		usrCreator: usrCreator,
	}
}

func (srv RegisterService) Register(
	ctx context.Context,
	login string,
	password string,
	birthdayDate time.Time,
) (string, error) {

	encryptedPassword, err := auth.HashPassword(password)
	if err != nil {
		return "", fmt.Errorf("failed to register user: %w", err)
	}

	user, err := srv.usrCreator.CreateUser(ctx, login, encryptedPassword, birthdayDate)
	if err != nil {
		return "", fmt.Errorf("failed to register user: %w", err)
	}

	jwtStr, err := auth.BuildJWTString(user.ID)
	if err != nil {
		return "", fmt.Errorf("failed to register user: %w", err)
	}

	return jwtStr, nil
}
