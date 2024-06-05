package services

import (
	"context"

	"github.com/ilya-burinskiy/birthday-notify/internal/models"
)

type NotificationSettingUpdater interface {
	UpdateNotificationSetting(ctx context.Context, settingID, daysBeforeNotify int) (models.NotifySetting, error)
}

type UpdateNotificationSettingService struct {
	updater NotificationSettingUpdater
}

func NewUpdateNotificationService(updater NotificationSettingUpdater) UpdateNotificationSettingService {
	return UpdateNotificationSettingService{
		updater: updater,
	}
}

func (srv UpdateNotificationSettingService) UpdateNotificationSetting(
	ctx context.Context,
	settingID,
	daysBeforeNotify int,
) (models.NotifySetting, error) {

	// TODO: add authorization
	return srv.updater.UpdateNotificationSetting(ctx, settingID, daysBeforeNotify)
}
