package services

import (
	"context"

	"github.com/ilya-burinskiy/birthday-notify/internal/models"
)

type UserFetcher interface {
	FetchUsers(ctx context.Context) ([]models.User, error)
}

type FetchUserService struct {
	fetcher UserFetcher
}

func NewFetchUserService(fetcher UserFetcher) FetchUserService {
	return FetchUserService{
		fetcher: fetcher,
	}
}

func (srv FetchUserService) FetchUsers(ctx context.Context) ([]models.User, error) {
	return srv.fetcher.FetchUsers(ctx)
}
