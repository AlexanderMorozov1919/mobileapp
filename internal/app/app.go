package app

import (
	"context"
	"log"
	"net/http"

	"github.com/AlexanderMorozov1919/mobileapp/internal/adapters/handlers"
	"github.com/AlexanderMorozov1919/mobileapp/internal/adapters/repositories"
	"github.com/AlexanderMorozov1919/mobileapp/internal/adapters/repositories/auth"
	"github.com/AlexanderMorozov1919/mobileapp/internal/config"
	"github.com/AlexanderMorozov1919/mobileapp/internal/services"
	"github.com/AlexanderMorozov1919/mobileapp/internal/usecases"
	"github.com/AlexanderMorozov1919/mobileapp/migrations"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

func New() *fx.App {
	return fx.New(
		fx.Provide(
			config.LoadConfig,
			// Добавляем провайдер для JWT secret
			func(cfg *config.Config) string { return cfg.JWTSecret },
		),
		RepositoryModule,
		ServiceModule,
		UsecaseModule,
		HttpServerModule,
		SeedModule,
	)
}

// TOFIX: cfg.JWTSecret undefined (type *config.Config has no field or method JWTSecret)
var AuthModule = fx.Module("auth_module",
	fx.Provide(
		auth.NewAuthRepository,
		func(cfg *config.Config) string {
			if cfg == nil {
				panic("config is nil")
			}
			return cfg.JWTSecret
		},
		usecases.NewAuthUsecase,
	),
)

var SeedModule = fx.Module("seed_module",
	fx.Invoke(func(db *gorm.DB) {
		if err := migrations.SeedTestDoctors(db); err != nil {
			log.Printf("Warning: failed to seed test doctors: %v", err)
		} else {
			log.Println("Test doctors seeded successfully")
		}
	}),
)

/*
func ProvideLoggers(cfg *config.Config) *logging.Logger {
	loggerCfg := logging.Config(
		cfg.Logging.Enable,
		cfg.Logging.Level,
		cfg.Logging.Format,
		cfg.Logging.LogsDir,
		IntToUint(cfg.Logging.SavingDays),
	)

	logger := logging.NewBaseLogger("VERSION", cfg.App.Version, loggerCfg, logging.WithDailyLogDelete())

	return logger
}

var LoggingModule = fx.Module("logging_module",
	fx.Provide(
		ProvideLoggers,
	),
	fx.Invoke(logging.InvokeBaseLogger),
)


*/
/* -------------------------------------------- */

func InvokeHttpServer(lc fx.Lifecycle, h http.Handler) {
	server := &http.Server{
		Addr:    ":8080", // Упрощаем - используем хардкод порта
		Handler: h,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go server.ListenAndServe()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Close()
		},
	})
}

var HttpServerModule = fx.Module("http_server_module",
	fx.Provide(
		handlers.NewHandler,
		handlers.ProvideRouter,
	),

	fx.Invoke(InvokeHttpServer),
)

/* -------------------------------------------- */

var ServiceModule = fx.Module("service_module",
	fx.Provide(services.NewService),
)

/* -------------------------------------------- */

var RepositoryModule = fx.Module("postgres_module",
	fx.Provide(repositories.NewRepository),
)

/* -------------------------------------------- */

var UsecaseModule = fx.Module("usecases_module",
	fx.Provide(
		usecases.NewUsecases,
		usecases.NewAuthUsecase,
	),
)

/* -------------------------------------------- */

func IntToUint(c int) uint {
	if c < 0 {
		panic([2]any{"a negative number", c})
	}

	return uint(c)
}
