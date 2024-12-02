package application_test

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/leetcode-golang-classroom/golang-with-mongodb-sample/internal/application"
	"github.com/leetcode-golang-classroom/golang-with-mongodb-sample/internal/config"
	"github.com/leetcode-golang-classroom/golang-with-mongodb-sample/internal/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

type defaultRouteSuite struct {
	ctx       context.Context
	container testcontainers.Container
	app       *application.App
	suite.Suite
}

func TestDefaultRouteSuite(t *testing.T) {
	suite.Run(t, new(defaultRouteSuite))
}

func (s *defaultRouteSuite) SetupSuite() {
	ctx := context.WithValue(context.Background(), logger.CtxKey{}, slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	})))
	s.ctx = ctx
	c, err := mongodb.Run(ctx, "mongo:latest",
		mongodb.WithUsername("eddie"),
		mongodb.WithPassword("test"),
	)
	require.NoError(s.T(), err)
	s.container = c
	url, err := c.ConnectionString(ctx)
	require.NoError(s.T(), err)
	config.AppConfig.MongoDBURL = url
	config.AppConfig.Environment = "TEST"
	app, err := application.New(config.AppConfig, ctx)
	app.SetupRoutes(ctx)
	require.NoError(s.T(), err)
	s.app = app
}

func (s *defaultRouteSuite) TearDownSuite() {
	err := s.container.Terminate(s.ctx)
	require.NoError(s.T(), err)
}
func (s *defaultRouteSuite) Test_DefaultRoute() {
	testCases := []struct {
		name            string
		expectedStatus  int
		expectedMessage string
	}{
		{
			name:            "default route",
			expectedStatus:  http.StatusOK,
			expectedMessage: "ok",
		},
	}
	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodGet, "/", nil)
			// Act
			s.app.Router.ServeHTTP(w, r)
			type response struct {
				Message string `json:"message"`
			}
			var resp response
			err := json.NewDecoder(w.Body).Decode(&resp)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.Equal(t, tc.expectedMessage, resp.Message)
		})
	}
}
