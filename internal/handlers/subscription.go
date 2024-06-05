package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/ilya-burinskiy/birthday-notify/internal/middlewares"
	"github.com/ilya-burinskiy/birthday-notify/internal/models"
	"go.uber.org/zap"
)

type SubscribeService interface {
	Subscribe(ctx context.Context, subscribedUserID, subscribingUserID int) (models.Subscription, error)
}

type UnsubscribeService interface {
	Unsubscribe(ctx context.Context, subscribedUserID, subscribingUserID int) error
}

type SubscriptionHandler struct {
	logger *zap.Logger
}

func NewSubscriptionHandler(logger *zap.Logger) SubscriptionHandler {
	return SubscriptionHandler{
		logger: logger,
	}
}

func (h SubscriptionHandler) Subscribe(subscribeSrv SubscribeService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		subscribingUserID, _ := middlewares.UserIDFromContext(r.Context())
		subscribedUserID, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			h.logger.Info("invalid user id", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, err = subscribeSrv.Subscribe(context.Background(), subscribedUserID, subscribingUserID)
		// TODO: respond with proper http status
		if err != nil {
			h.logger.Info("failed to subscribe user", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func (h SubscriptionHandler) Unsubscribe(unsubscribeSrv UnsubscribeService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		subscribingUserID, _ := middlewares.UserIDFromContext(r.Context())
		subscribedUserID, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			h.logger.Info("invalid user id", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = unsubscribeSrv.Unsubscribe(context.Background(), subscribedUserID, subscribingUserID)
		// TODO: respond with proper http status
		if err != nil {
			h.logger.Info("failed to unsubscribe user", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
