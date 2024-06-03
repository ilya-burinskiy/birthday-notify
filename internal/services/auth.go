package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/ilya-burinskiy/birthday-notify/internal/auth"
	"github.com/ilya-burinskiy/birthday-notify/internal/models"
)

type UserFinder interface {
	FindUserByEmail(ctx context.Context, email string) (models.User, error)
}

type AuthenticateService struct {
	userFinder UserFinder
}

func NewAuthenticateService(usrFinder UserFinder) AuthenticateService {
	return AuthenticateService{
		userFinder: usrFinder,
	}
}

func (srv AuthenticateService) Authenticate(ctx context.Context, email, password string) (string, error) {
	user, err := srv.userFinder.FindUserByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("failed to authenticate user: %w", err)
	}

	if !auth.ValidatePasswordHash(password, string(user.EncryptedPassword)) {
		return "", errors.New("invalid email or password")
	}

	jwtStr, err := auth.BuildJWTString(user.ID)
	if err != nil {
		return "", fmt.Errorf("failed to authenticate user: %w", err)
	}

	return jwtStr, nil
}
