package services

import (
	"context"

	"github.com/ilya-burinskiy/birthday-notify/internal/models"
)

type UsersFetcher interface {
	FetchUsers(ctx context.Context) ([]models.User, error)
}

type FetchUsersService struct {
	fetcher UsersFetcher
}

func NewFetchUsersService(fetcher UsersFetcher) FetchUsersService {
	return FetchUsersService{
		fetcher: fetcher,
	}
}

func (srv FetchUsersService) FetchUsers(ctx context.Context) ([]models.User, error) {
	return srv.fetcher.FetchUsers(ctx)
}
