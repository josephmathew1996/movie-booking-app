package middlewares

import (
	"movie-booking-app/users-service/pkg/user"
	"net/http"

	"github.com/labstack/echo/v4"
)

// func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		return next(c)
// 	}
// }

func WrapWithArgs(auth IAuth, repo user.Repository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get(echo.HeaderAuthorization)
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Authorization header is required")
			}
			claims, err := auth.Authenticate(authHeader)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
			}
			c.Set("userID", claims.UserID)
			//we could make a db call to fetch user detail by id parsed from token
			// and then compare it with userid present in the request to verify user level or admin level authorisation
			return next(c)
		}
	}
}
