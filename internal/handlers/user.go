package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/ilya-burinskiy/birthday-notify/internal/auth"
	"github.com/ilya-burinskiy/birthday-notify/internal/storage"
	"go.uber.org/zap"
)

type RegisterService interface {
	Register(ctx context.Context, email, password string, birthDate time.Time) (string, error)
}

type UserHandler struct {
	logger *zap.Logger
}

func NewUserHandlers(logger *zap.Logger) UserHandler {
	return UserHandler{}
}

func (h UserHandler) Register(regSrv RegisterService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		type payload struct {
			Email     string    `json:"email"`
			Password  string    `json:"password"`
			Birthdate time.Time `json:"birthdate"`
		}

		w.Header().Set("Content-Type", "application/json")
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
		jwtStr, err := regSrv.Register(
			r.Context(),
			requestBody.Email,
			requestBody.Password,
			requestBody.Birthdate,
		)
		if err != nil {
			var notUniqErr storage.ErrUserNotUniq
			if errors.As(err, &notUniqErr) {
				w.WriteHeader((http.StatusConflict))
				if err := encoder.Encode(err.Error()); err != nil {
					h.logger.Info("failed to encode response", zap.Error(err))
				}
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			if err := encoder.Encode(err.Error()); err != nil {
				h.logger.Info("failed to encode response", zap.Error(err))
			}
			return
		}
		auth.SetJWTCookie(w, jwtStr)
		w.WriteHeader(http.StatusOK)
	}
}
