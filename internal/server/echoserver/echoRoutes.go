package echo

import (
	"movie-booking-app/users-service/internal/middlewares"
	"movie-booking-app/users-service/pkg/user/http/echohttp"
)

func (s *EchoServer) RegisterV1Routes() {
	v1 := s.Echo.Group("/api/v1")

	userHandler := echohttp.NewHttpHandler(*s.Logger, s.Validator, s.Service)

	auth := v1.Group("/auth")
	auth.POST("/register", userHandler.RegisterUser)
	auth.POST("/login", userHandler.LoginUser)

	//api with middleware added
	user := v1.Group("/users", middlewares.WrapWithArgs(s.Auth, s.Repo))
	user.GET("/:id", userHandler.GetUser)
	// user.GET("", echohttp.GetUsers)
	// user.PUT("/:id", echohttp.UpdateUser)
	// user.DELETE("/:id", echohttp.DeleteUser)
}
