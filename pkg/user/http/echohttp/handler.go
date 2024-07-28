package echohttp

import (
	"fmt"
	"movie-booking-app/users-service/internal/errors"
	"movie-booking-app/users-service/pkg/models"
	"movie-booking-app/users-service/pkg/user"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type HttpHandler struct {
	logger    zap.SugaredLogger
	validator *validator.Validate
	service   user.Service
}

func NewHttpHandler(
	logger zap.SugaredLogger,
	validator *validator.Validate,
	service user.Service,
) *HttpHandler {
	handler := &HttpHandler{
		logger:    logger,
		validator: validator,
		service:   service,
	}
	return handler
}

// Handler function to register a new user
func (h *HttpHandler) RegisterUser(c echo.Context) error {
	var (
		req models.User
	)

	// Bind and validate the request payload
	if err := c.Bind(&req); err != nil {
		h.logger.Errorln("Error parsing the user request payload: ", err.Error())
		// Create an HTTP error and set the internal error

		// in built error handler
		// return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload").SetInternal(err)

		//custom error handler
		return errors.New(errors.ErrBadRequest, "Invalid request payload", err)
	}
	if err := h.validator.Struct(&req); err != nil {
		h.logger.Errorln("Error validating the user request payload: ", err.Error())
		fields := []string{}
		if val, ok := err.(validator.ValidationErrors); ok {
			for _, err := range val {
				fields = append(fields, strings.ToLower(err.Field()))
			}
		}
		if len(fields) > 0 {
			// Create an HTTP error and set the internal error
			return echo.NewHTTPError(http.StatusBadRequest, "Please provide valid value(s) for user registration: "+strings.Join(fields, ", ")).SetInternal(err)
		}
		// Create an HTTP error and set the internal error
		return echo.NewHTTPError(http.StatusBadRequest, "Validation failed").SetInternal(err)
	}

	// Call the service to create the user
	res, err := h.service.CreateUser(c.Request().Context(), req)
	if err != nil {
		h.logger.Errorln("Error creating the user: ", err.Error())
		// Create an HTTP error and set the internal error
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create user").SetInternal(err)
	}

	return c.JSON(http.StatusCreated, res)
}

// Handler function to login a new user
func (h *HttpHandler) LoginUser(c echo.Context) error {
	// Logic to login a user
	var (
		req           models.LoginUserRequest
		res           interface{}
		err           error
		errMessageStr string
		resStatusCode int
	)

	defer func() {
		if err != nil {
			fields := []string{}
			errMessageStr = err.Error()
			h.logger.Error("Error logging in a user: ", err)

			if val, ok := err.(validator.ValidationErrors); ok {
				for _, err := range val {
					fields = append(fields, strings.ToLower(err.Field()))
				}
			}

			if len(fields) > 0 {
				errMessageStr = "Please provide valid value(s) for user login: " + strings.Join(fields, ", ")
			}

			res = models.ErrorResponse{
				Code:    resStatusCode,
				Message: errMessageStr,
			}

			c.JSON(resStatusCode, res)
			return
		}
	}()
	// request validation
	err = c.Bind(&req)
	if err != nil {
		h.logger.Errorln("Error parsing the user login request payload: ", err.Error())
		resStatusCode = 400
		return err
	}
	err = h.validator.Struct(&req)
	if err != nil {
		h.logger.Errorln("Error validating the user login request payload: ", err.Error())
		resStatusCode = 400
		return err
	}
	res, err = h.service.LoginUser(c.Request().Context(), req)
	if err != nil {
		h.logger.Errorln("Error logging in the user: ", err.Error())
		resStatusCode = 500
		return err
	}
	return c.JSON(http.StatusOK, res)
}

// // Handler function to get all users
// func GetUsers(c echo.Context) error {
// 	// Logic to get all users
// 	return c.JSON(http.StatusOK, map[string]interface{}{"items": []string{"User1"}})
// }

// Handler function to get a user by ID
func (h *HttpHandler) GetUser(c echo.Context) error {
	// Logic to get a user by ID
	var (
		res           interface{}
		err           error
		errMessageStr string
		resStatusCode int
	)

	defer func() {
		if err != nil {
			errMessageStr = err.Error()
			h.logger.Error("Error fectching user by id: ", err)

			res = models.ErrorResponse{
				Code:    resStatusCode,
				Message: errMessageStr,
			}

			c.JSON(resStatusCode, res)
			return
		}
	}()

	userid := c.Get("userID")
	fmt.Println("userid form context", userid)
	userIdStr := c.Param("id")
	fmt.Println("user id from request", userIdStr)
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		resStatusCode = 400
		return err
	}
	res, err = h.service.GetUser(c.Request().Context(), userId)
	if err != nil {
		resStatusCode = 500
		return err
	}
	return c.JSON(http.StatusOK, res)
}

// // Handler function to update a user by ID
// func UpdateUser(c echo.Context) error {
// 	// Logic to update a user by ID
// 	return c.String(http.StatusOK, "User updated")
// }

// // Handler function to delete a user by ID
// func DeleteUser(c echo.Context) error {
// 	// Logic to delete a user by ID
// 	return c.String(http.StatusOK, "User deleted")
// }
