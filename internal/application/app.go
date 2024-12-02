package application

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/leetcode-golang-classroom/golang-with-mongodb-sample/internal/config"
	"github.com/leetcode-golang-classroom/golang-with-mongodb-sample/internal/logger"
	"github.com/leetcode-golang-classroom/golang-with-mongodb-sample/internal/mongodb"
	"github.com/newrelic/go-agent/v3/newrelic"
	"go.mongodb.org/mongo-driver/mongo"
)

// define app dependency
type App struct {
	Router        *gin.Engine
	config        *config.Config
	mongodbClient *mongo.Client
	nrApp         *newrelic.Application
}

func New(config *config.Config, ctx context.Context) (*App, error) {
	// 設定 mongodb 連線
	mongodbClient, err := mongodb.NewMongoClient(config.MongoDBURL, 10)
	if err != nil {
		return nil, err
	}
	var nrApp *newrelic.Application
	if config.Environment == "PROD" {
		nrApp, err = newrelic.NewApplication(
			newrelic.ConfigAppName(config.AppName),
			newrelic.ConfigLicense(config.NewRelicLicenseKey),
			newrelic.ConfigInfoLogger(os.Stdout),
			newrelic.ConfigAppLogForwardingEnabled(true),
		)
		if err != nil {
			return nil, err
		}
	}
	app := &App{
		config:        config,
		mongodbClient: mongodbClient,
		nrApp:         nrApp,
	}
	app.SetupRoutes(ctx)
	return app, nil
}

func (app *App) Start(ctx context.Context) error {
	logger := logger.FromContext(ctx)
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", app.config.Port),
		Handler: app.Router,
	}
	logger.Info(fmt.Sprintf("Starting server on %s", app.config.Port))
	errCh := make(chan error, 1)
	defer func() {
		logger.Warn(fmt.Sprintf("closing mongodb connection with url %s", app.config.MongoDBURL))
		err := app.mongodbClient.Disconnect(ctx)
		if err != nil {
			logger.Error(fmt.Sprintf("close mongodb failed with url %s", app.config.MongoDBURL))
		}
		logger.Info(fmt.Sprintf("close mongodb connection with url %s", app.config.MongoDBURL))
	}()
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			errCh <- fmt.Errorf("failed to start server: %w", err)
		}
		close(errCh)
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		logger.Info("server cancel")
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		return server.Shutdown(timeout)
	}
}
