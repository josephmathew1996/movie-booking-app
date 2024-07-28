package echo

import (
	"movie-booking-app/users-service/internal/config"
	"movie-booking-app/users-service/internal/database"
	"movie-booking-app/users-service/internal/database/mysql"
	"movie-booking-app/users-service/internal/errors"
	"movie-booking-app/users-service/internal/logger"
	"movie-booking-app/users-service/internal/middlewares"
	"movie-booking-app/users-service/internal/ratelimiter"
	"movie-booking-app/users-service/pkg/user"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

type EchoServer struct {
	//add all dependencies
	Config      config.Config
	Logger      *zap.SugaredLogger
	Validator   *validator.Validate
	RateLimiter ratelimiter.Limiter
	DB          database.Database
	Echo        *echo.Echo
	Service     user.Service
	Repo        user.Repository
	Auth        middlewares.IAuth
}

func NewEchoServer() (*EchoServer, error) {
	// Initialize configuration
	cfg := config.InitializeConfig()

	// Initialize logger
	log, err := logger.InitializeLogger(cfg.App.LogLevel)
	if err != nil {
		return nil, err
	}

	db, err := mysql.NewMySQLDatabase(cfg.App.DBConnectionString)
	if err != nil {
		return nil, err
	}

	// Rate limiter per client
	rateLimiter := ratelimiter.NewFixedWindowClientsRateLimiter(
		cfg.RateLimiter.RequestsPerTimeFrame,
		cfg.RateLimiter.TimeFrame,
	)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middlewares.RateLimiterMiddleware(rateLimiter, cfg.RateLimiter.Enabled, log))

	e.HTTPErrorHandler = CustomHTTPErrorHandler

	// add all dependencies to user domain here
	repo := user.NewUserRepo(db)
	service := user.NewUserService(repo, db)
	auth := middlewares.NewAuth()

	return &EchoServer{
		Config:      cfg,
		Logger:      log,
		Validator:   validator.New(),
		RateLimiter: rateLimiter,
		DB:          db,
		Echo:        e,
		Service:     service,
		Repo:        repo,
		Auth:        auth,
	}, nil
}

func CustomHTTPErrorHandler(err error, c echo.Context) {
	var (
		code    = http.StatusInternalServerError
		message interface{}
	)

	if customErr, ok := err.(*errors.CustomError); ok {
		switch customErr.Code {
		case errors.ErrBadRequest:
			code = http.StatusBadRequest
		case errors.ErrValidationFailed:
			code = http.StatusUnprocessableEntity
		case errors.ErrUserAlreadyExists:
			code = http.StatusConflict
		case errors.ErrInternalServer:
			code = http.StatusInternalServerError
		}
		message = echo.Map{"code": customErr.Code, "message": customErr.Message}
	} else {
		message = echo.Map{"message": err.Error()}
	}

	if !c.Response().Committed {
		c.JSON(code, message)
	}
}

func (s *EchoServer) Start() {
	port := s.Config.App.Port
	portStr := strconv.Itoa(port)
	s.Logger.Info("Starting server on port " + portStr)
	if err := s.Echo.Start(":" + portStr); err != nil {
		s.Logger.Fatal("Failed to start echo server", zap.Error(err))
	}
}
