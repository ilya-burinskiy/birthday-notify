package services

import (
	"context"

	"github.com/ilya-burinskiy/birthday-notify/internal/models"
)

type SubscriptionCreator interface {
	CreateSubscription(ctx context.Context, subscribedUserID, subscribingUserID int) (models.Subscription, error)
}

type SubscribeService struct {
	creator SubscriptionCreator
}

func NewSubscribeService(creator SubscriptionCreator) SubscribeService {
	return SubscribeService{
		creator: creator,
	}
}

func (srv SubscribeService) Subscribe(
	ctx context.Context,
	subscribedUserID,
	subscribingUserID int,
) (models.Subscription, error) {

	return srv.creator.CreateSubscription(ctx, subscribedUserID, subscribingUserID)
}
