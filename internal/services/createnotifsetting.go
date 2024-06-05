package services

import (
	"context"

	"github.com/ilya-burinskiy/birthday-notify/internal/models"
)

type NotificationSettingCreator interface {
	CreateNotificationSetting(ctx context.Context, userID, daysBeforeNotify int) (models.NotifySetting, error)
}

type CreateNotificationSettingService struct {
	creator NotificationSettingCreator
}

func NewCreateNotificationSettingService(creator NotificationSettingCreator) CreateNotificationSettingService {
	return CreateNotificationSettingService{
		creator: creator,
	}
}

func (srv CreateNotificationSettingService) CreateNotificationSetting(
	ctx context.Context,
	userID,
	daysBeforeNotiy int,
) (models.NotifySetting, error) {

	return srv.creator.CreateNotificationSetting(ctx, userID, daysBeforeNotiy)
}
