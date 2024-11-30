package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/leetcode-golang-classroom/golang-with-mongodb-sample/internal/application"
	"github.com/leetcode-golang-classroom/golang-with-mongodb-sample/internal/config"
	mlog "github.com/leetcode-golang-classroom/golang-with-mongodb-sample/internal/logger"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(
		os.Stdout, &slog.HandlerOptions{
			AddSource: true,
		},
	))
	rootContext := context.WithValue(context.Background(), mlog.CtxKey{}, logger)
	app, err := application.New(config.AppConfig, rootContext)
	if err != nil {
		logger.Error("failed to build app", "error", err)
		return
	}
	ctx, cancel := signal.NotifyContext(rootContext, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()
	err = app.Start(ctx)
	if err != nil {
		logger.Error("failed to start app", "error", err)
	}
}
