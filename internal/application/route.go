package application

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leetcode-golang-classroom/golang-with-mongodb-sample/internal/logger"
	"github.com/leetcode-golang-classroom/golang-with-mongodb-sample/internal/services/movie"
	nrgin "github.com/newrelic/go-agent/v3/integrations/nrgin"
	sloggin "github.com/samber/slog-gin"
)

// define route
func (app *App) SetupRoutes(ctx context.Context) {
	gin.SetMode(app.config.GinMode)
	router := gin.New()
	// recovery middleward
	router.Use(sloggin.New(logger.FromContext(ctx)))
	router.Use(gin.Recovery())
	// setup router for new relic
	if app.config.Environment == "PROD" {
		router.Use(nrgin.Middleware(app.nrApp))
	}
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, map[string]string{"message": "ok"})
	})
	app.Router = router
	app.SetupNewsRoutes()
}

func (app *App) SetupNewsRoutes() {
	movieGroups := app.Router.Group("/movies")
	store := movie.NewStore(app.mongodbClient, app.config)
	handler := movie.NewHandler(store)
	handler.RegisterRoute(movieGroups)
}
