package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/ilya-burinskiy/birthday-notify/internal/middlewares"
	"github.com/ilya-burinskiy/birthday-notify/internal/models"
	"go.uber.org/zap"
)

type CreateNotificationSettingService interface {
	CreateNotificationSetting(ctx context.Context, userID, daysBeforeNotify int) (models.NotifySetting, error)
}

type UpdateNotificationSettingService interface {
	UpdateNotificationSetting(ctx context.Context, settingID, daysBeforeNotify int) (models.NotifySetting, error)
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

func (h NotificationSettingHandler) Update(updateSrv UpdateNotificationSettingService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		type payload struct {
			DaysBeforeNotify int `json:"days_before_notify"`
		}

		w.Header().Set("Content-Type", "application/json")
		var requestBody payload
		decoder := json.NewDecoder(r.Body)
		encoder := json.NewEncoder(w)
		err := decoder.Decode(&requestBody)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			if err := encoder.Encode("invalid request body"); err != nil {
				h.logger.Info("failed to encode body")
			}
			return
		}

		settingID, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			h.logger.Info("invalid notify setting id", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		notifySetting, err := updateSrv.UpdateNotificationSetting(
			context.Background(),
			settingID,
			requestBody.DaysBeforeNotify,
		)
		if err != nil {
			h.logger.Info("failed to update notification setting", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := encoder.Encode(notifySetting); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.logger.Info("failed to encode response body", zap.Error(err))
			return
		}
	}
}
