package services

import (
	"context"
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/ilya-burinskiy/birthday-notify/internal/models"
	"go.uber.org/zap"
)

type NotificationsForCurrentDateFetcher interface {
	FetchNotificationsForCurrentDate(ctx context.Context) ([]models.Notification, error)
}

type NotificationSender interface {
	Send(to string, subject string, body string) error
}

type Notifier struct {
	logger  *zap.Logger
	fetcher NotificationsForCurrentDateFetcher
	sender  NotificationSender
}

func NewNotifier(logger *zap.Logger, fetcher NotificationsForCurrentDateFetcher, sender NotificationSender) Notifier {
	return Notifier{
		logger:  logger,
		fetcher: fetcher,
		sender:  sender,
	}
}

func (notifier Notifier) Start() {
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.Cron("0 12 * * *").Do(notifier.notify)
	scheduler.StartAsync()
}

func (notifier Notifier) notify() {
	notifications, err := notifier.fetcher.FetchNotificationsForCurrentDate(context.Background()) 
	if err != nil {
		notifier.logger.Info("failed to fetch notifications", zap.Error(err))
		return
	}

	subject := "Birthday notification"
	body := "The user %s has birthday in %d days"
	for _, notification := range notifications {
		err := notifier.sender.Send(
			notification.SubscribingUserEmail,
			subject,
			fmt.Sprintf(body, notification.SubscribedUserEmail, notification.DaysBeforeNotify),
		)
		if err != nil {
			notifier.logger.Info("failed to send notifiaction", zap.Error(err))
		}
	}
}
