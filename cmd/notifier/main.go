package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ilya-burinskiy/birthday-notify/internal/configs"
	"github.com/ilya-burinskiy/birthday-notify/internal/handlers"
	"github.com/ilya-burinskiy/birthday-notify/internal/middlewares"
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
	subscribeSrv := services.NewSubscribeService(store)
	unsubscribeSrv := services.NewUnsubscribeService(store, store)
	notifySettingCreator := services.NewCreateNotificationSettingService(store)

	emailSender := services.NewEmailSender(config)
	notifier := services.NewNotifier(logger, store, emailSender)
	notifier.Start()

	router := chi.NewRouter()
	configureUserRouter(logger, registerSrv, authSrv, router)
	configureSubscriptionRouter(logger, subscribeSrv, unsubscribeSrv, router)
	configureNotificationSettingRouter(logger, notifySettingCreator, router)

	server := http.Server{
		Handler: router,
		Addr:    config.RunAddr,
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

func configureSubscriptionRouter(
	logger *zap.Logger,
	subscribeSrv services.SubscribeService,
	unsubscribeSrv services.UnsubscribeService,
	mainRouter chi.Router) {

	handler := handlers.NewSubscriptionHandler(logger)
	mainRouter.Group(func(router chi.Router) {
		router.Use(middlewares.Authenticate)
		router.Post("/api/users/{id}/subscribe", handler.Subscribe(subscribeSrv))
		router.Delete("/api/users/{id}/unsubscribe", handler.Unsubscribe(unsubscribeSrv))
	})
}

func configureNotificationSettingRouter(
	logger *zap.Logger,
	createSrv services.CreateNotificationSettingService,
	mainRouter chi.Router) {

	handler := handlers.NewNotificationSettingHandler(logger)
	mainRouter.Group(func(router chi.Router) {
		router.Use(middlewares.Authenticate)
		router.Post("/api/notify_settings", handler.Create(createSrv))
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
