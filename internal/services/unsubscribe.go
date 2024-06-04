package services

import (
	"context"
	"fmt"

	"github.com/ilya-burinskiy/birthday-notify/internal/models"
)

type SubscriptionFinder interface {
	FindSubscription(ctx context.Context, subscribedUserID, subscribingUserID int) (models.Subscription, error)
}

type SubscriptionDeleter interface {
	DeleteSubscription(ctx context.Context, subscriptionID int) error
}

type UnsubscribeService struct {
	finder  SubscriptionFinder
	deleter SubscriptionDeleter
}

func NewUnsubscribeService(finder SubscriptionFinder, deleter SubscriptionDeleter) UnsubscribeService {
	return UnsubscribeService{
		finder:  finder,
		deleter: deleter,
	}
}

func (srv UnsubscribeService) Unsubscribe(ctx context.Context, subscribedUserID, subscribingUserID int) error {
	subscription, err := srv.finder.FindSubscription(ctx, subscribedUserID, subscribingUserID)
	if err != nil {
		return fmt.Errorf("failed to find subscription: %w", err)
	}

	if err := srv.deleter.DeleteSubscription(ctx, subscription.ID); err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}

	return nil
}
