package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/cmd/web-api/config"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/app"
)

var (
	AppName      = "users"
	AppRelease   = "unspecified"
	AppCommit    = "unspecified"
	AppBuildTime = time.Now().Format("06-01-02_15:04:05")
)

func main() {
	// Парсим конфигурацию приложения
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("can't parse app config: %v", err)
	}

	if cfg.MaxCpu > 0 {
		runtime.GOMAXPROCS(cfg.MaxCpu)
	}

	ApplicationInfo := &app.AppInfo{
		Name:      AppName,
		BuildTime: AppBuildTime,
		Commit:    AppCommit,
		Release:   AppRelease,
	}

	// Инициализируем логгер
	logger, err := initLogger(ApplicationInfo, cfg)
	if err != nil {
		log.Fatalf("can't init logger: %v", err)
		return
	}
	defer func() {
		if e := recover(); e != nil {
			logger.Fatal("panic error", zap.Error(fmt.Errorf("%s", e)))
		}
		if logger, ok := logger.(ditzap.LoggerWithSentry); ok {
			logger.Flush()
		}
	}()
	logger.ReplaceLogger()

	logger.Info(fmt.Sprintf("Application `%s` %s started.", AppName, AppRelease))

	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	// Обработка syscall
	signalHandler(logger, cancelCtx)

	application := app.NewApp(cfg, ApplicationInfo, logger)
	application.Run(ctx)
	logger.Info(fmt.Sprintf("Application `%s` %s is shutdown.", AppName, AppRelease))
}

func initLogger(info *app.AppInfo, cfg *config.Config) (ditzap.Logger, error) {
	var (
		logger ditzap.Logger
		err    error
	)
	if cfg.SentryDSN != "" {
		logger, err = ditzap.NewLoggerWithSentry(&ditzap.LoggerSentryParams{
			LoggerParams: &ditzap.LoggerParams{
				LogLevel:           ditzap.LevelFromString(cfg.LogLevel),
				DevMode:            cfg.DevMode,
				OutputPaths:        []string{"stdout"},
				CallerLevelsToSkip: 1,
				StackTraceLevel:    ditzap.ErrorLevel,
			},
			AppName:     info.Name,
			SentryDSN:   cfg.SentryDSN,
			Release:     info.GetReleaseVersion(),
			Environment: cfg.Environment,
		}, "file_id")
	} else {
		logger, err = ditzap.NewLogger(&ditzap.LoggerParams{
			LogLevel:           ditzap.LevelFromString(cfg.LogLevel),
			DevMode:            cfg.DevMode,
			OutputPaths:        []string{"stdout"},
			CallerLevelsToSkip: 1,
			StackTraceLevel:    ditzap.ErrorLevel,
		})
	}
	return logger, err //nolint:wrapcheck
}

// signalHandler обработчик сигналов системы
func signalHandler(logger ditzap.Logger, cancelFunc context.CancelFunc) {
	osSigCh := make(chan os.Signal, 1)

	signal.Notify(
		osSigCh,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)

	go func() {
		defer signal.Stop(osSigCh)
		s := <-osSigCh
		switch s {
		case syscall.SIGHUP:
			logger.Info("Received signal SIGHUP! Application shutdown")
		case syscall.SIGINT:
			logger.Info("Received signal SIGINT! Application shutdown")
		case syscall.SIGQUIT:
			logger.Info("Received signal SIGQUIT! Application shutdown")
		case syscall.SIGTERM:
			logger.Info("Received signal SIGTERM! Application shutdown")
		}
		cancelFunc()
	}()
}
