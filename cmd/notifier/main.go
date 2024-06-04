package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ilya-burinskiy/birthday-notify/internal/configs"
	"github.com/ilya-burinskiy/birthday-notify/internal/handlers"
	"github.com/ilya-burinskiy/birthday-notify/internal/services"
	"github.com/ilya-burinskiy/birthday-notify/internal/storage"
	"go.uber.org/zap"
)

func main() {
	config := configs.Parse()
	store, err := storage.NewDBStorage(config.DSN)
	if err != nil {
		panic(err)
	}
	logger := configureLogger("info")

	registerSrv := services.NewRegisterService(store)
	authSrv := services.NewAuthenticateService(store)

	router := chi.NewRouter()
	configureUserRouter(logger, registerSrv, authSrv, router)

	server := http.Server{
		Handler: router,
		Addr: config.RunAddr,
	}
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}

func configureUserRouter(
	logger *zap.Logger,
	registerSrv services.RegisterService,
	authSrv services.AuthenticateService,
	mainRouter chi.Router) {

	handler := handlers.NewUserHandlers(logger)
	mainRouter.Group(func(router chi.Router) {
		router.Use(middleware.AllowContentType("application/json"))
		router.Post("/api/users/register", handler.Register(registerSrv))
		router.Post("/api/users/login", handler.Authenticate(authSrv))
	})
}

func configureLogger(level string) *zap.Logger {
	logLvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		panic(err)
	}
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.Level = logLvl
	logger, err := loggerConfig.Build()
	if err != nil {
		panic(err)
	}

	return logger
}
