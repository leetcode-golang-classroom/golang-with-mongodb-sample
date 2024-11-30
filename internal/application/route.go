package application

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leetcode-golang-classroom/golang-with-mongodb-sample/internal/logger"
	sloggin "github.com/samber/slog-gin"
)

// define route
func (app *App) SetupRoutes(ctx context.Context) {
	gin.SetMode(app.config.GinMode)
	router := gin.New()
	// recovery middleward
	router.Use(sloggin.New(logger.FromContext(ctx)))
	router.Use(gin.Recovery())
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, map[string]string{"message": "ok"})
	})
	app.Router = router
}
