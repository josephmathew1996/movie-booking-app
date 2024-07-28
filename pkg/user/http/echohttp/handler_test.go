package echohttp_test

import (
	"movie-booking-app/users-service/internal/errors"
	"movie-booking-app/users-service/pkg/models"
	"movie-booking-app/users-service/pkg/user/http/echohttp"
	"movie-booking-app/users-service/pkg/user/http/echohttp/fixtures"
	"movie-booking-app/users-service/pkg/user/mocks"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// func TestRegisterUser(t *testing.T) {
// 	tests := []struct {
// 		name                string
// 		reqBody             models.User
// 		mockedCreateUserRes interface{}
// 		expectedRes         interface{}
// 		expectedErr         error
// 		expectedStatusCode  int
// 	}{
// 		{"success", fixtures.User1, fixtures.User1, fixtures.User1, nil, 201},
// 		{"error", fixtures.User2, nil, nil, fmt.Errorf(""), 400},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			serviceMock := &mocks.Service{}
// 			repoMock := &mocks.Repository{}
// 			testLogger := zap.NewExample().Sugar()
// 			echoServer := &echoServer.EchoServer{
// 				Logger:    testLogger,
// 				Validator: validator.New(),
// 				Echo:      echo.New(),
// 				Service:   serviceMock,
// 				Repo:      repoMock,
// 			}
// 			echoServer.RegisterV1Routes()

// 			serviceMock.On(
// 				"CreateUser",
// 				mock.Anything,
// 				mock.Anything,
// 			).Return(
// 				tt.mockedCreateUserRes,
// 			)

// 			w := httptest.NewRecorder()
// 			reqBody, err := json.Marshal(tt.reqBody)
// 			// t.Error(err, "Invalid request body")
// 			req, err := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(reqBody))
// 			if err != nil {
// 				// t.Errorf(err, "Failed to create enrollmentcreate api request, err: %s")
// 			}
// 			req.Header.Add("Content-Type", "application/json")
// 			req.Header.Add("Authorization", "Bearer 123")
// 			// echoServer.Echo.Router.ServeHTTP(w, req)

// 			// bs, err := io.ReadAll(w.Body)
// 			// t.NoError(err, "Invalid response")

// 			// expectedBS, err := json.Marshal(tc.expectedRes)
// 			// t.NoError(err, "Invalid expected response")
// 			// t.Equal(tc.expectedStatusCode, w.Code)
// 			// t.Equal(string(expectedBS), string(bs))
// 		})
// 	}
// }

func TestRegisterUser(t *testing.T) {
	e := echo.New()
	testLogger := zap.NewExample().Sugar()
	mockService := &mocks.Service{}
	handler := echohttp.NewHttpHandler(*testLogger, validator.New(), mockService)

	tests := []struct {
		name               string
		reqBody            string
		mockedServiceRes   models.User
		mockedServiceErr   error
		expectedStatusCode int
		expectedResponse   string
		expectedErrorCode  interface{}
		expectedErrorMsg   string
	}{
		{
			name:               "InvalidRequestPayload",
			reqBody:            `invalid json`,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  errors.ErrBadRequest,
			expectedErrorMsg:   "Invalid request payload",
		},
		{
			name:               "Success",
			reqBody:            `{"user_name":"john_doe","first_name":"John","email":"john.doe@example.com","password":"secured_password"}`,
			mockedServiceRes:   fixtures.User1,
			expectedStatusCode: http.StatusCreated,
			expectedResponse:   `{"id":1,"user_name":"john_doe","first_name":"John","last_name":"","email":"john.doe@example.com"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Success" {
				// Set up the mock service to return the user and no error
				mockService.On("CreateUser", mock.Anything, mock.Anything).Return(tt.mockedServiceRes, tt.mockedServiceErr)
			}
			req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(tt.reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := handler.RegisterUser(c)
			if tt.expectedErrorCode != nil {
				if assert.Error(t, err) {
					he, ok := err.(*errors.CustomError)
					if assert.True(t, ok) {
						assert.Equal(t, tt.expectedErrorCode, he.Code)
						assert.Equal(t, tt.expectedErrorMsg, he.Message)
					}
				}
			} else {
				if assert.NoError(t, err) {
					assert.Equal(t, tt.expectedStatusCode, rec.Code)
					assert.JSONEq(t, tt.expectedResponse, rec.Body.String())
				}
			}

			// Assert that the mock method was called with the expected parameters
			if tt.name == "Success" {
				mockService.AssertExpectations(t)
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	e := echo.New()
	testLogger := zap.NewExample().Sugar()
	mockService := &mocks.Service{}
	handler := echohttp.NewHttpHandler(*testLogger, validator.New(), mockService)

	tests := []struct {
		name               string
		reqParam           string
		mockedServiceRes   models.User
		mockedServiceErr   error
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name:               "InvalidRequestParam",
			reqParam:           "a",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   "",
		},
		{
			name:               "Success",
			reqParam:           "1",
			mockedServiceRes:   fixtures.User1,
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"id":1,"user_name":"john_doe","first_name":"John","last_name":"","email":"john.doe@example.com"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Success" {
				// Set up the mock service to return the user and no error
				mockService.On("GetUser", mock.Anything, mock.Anything).Return(tt.mockedServiceRes, tt.mockedServiceErr)
			}
			req := httptest.NewRequest(http.MethodGet, "/users/"+tt.reqParam, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.reqParam)

			handler.GetUser(c)
			assert.Equal(t, tt.expectedStatusCode, rec.Code)
			if tt.expectedResponse != "" {
				assert.JSONEq(t, tt.expectedResponse, rec.Body.String())
			}

			// Assert that the mock method was called with the expected parameters
			if tt.name == "Success" {
				mockService.AssertExpectations(t)
			}
		})
	}
}
