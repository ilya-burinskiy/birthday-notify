package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ilya-burinskiy/birthday-notify/internal/middlewares"
	"github.com/ilya-burinskiy/birthday-notify/internal/models"
	"go.uber.org/zap"
)

type CreateNotificationSettingService interface {
	CreateNotificationSetting(ctx context.Context, userID, daysBeforeNotify int) (models.NotifySetting, error)
}

type NotificationSettingHandler struct {
	logger *zap.Logger
}

func NewNotificationSettingHandler(logger *zap.Logger) NotificationSettingHandler {
	return NotificationSettingHandler{
		logger: logger,
	}
}

func (h NotificationSettingHandler) Create(createSrv CreateNotificationSettingService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		type payload struct {
			DaysBeforeNotify int `json:"days_before_notify"`
		}

		w.Header().Set("Content-Type", "application-json")
		var requestBody payload
		decoder := json.NewDecoder(r.Body)
		encoder := json.NewEncoder(w)
		err := decoder.Decode(&requestBody)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			if err := encoder.Encode("invalid request body"); err != nil {
				h.logger.Info("failed to encode response", zap.Error(err))
			}
			return
		}

		userID, _ := middlewares.UserIDFromContext(r.Context())
		notifSetting, err := createSrv.CreateNotificationSetting(
			context.Background(),
			userID,
			requestBody.DaysBeforeNotify,
		)

		// TODO: response with proper http status
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			h.logger.Info("failed to create notification setting", zap.Error(err))
			return
		}

		if err := encoder.Encode(notifSetting); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.logger.Info("failed to encode response", zap.Error(err))
			return
		}
	}
}
