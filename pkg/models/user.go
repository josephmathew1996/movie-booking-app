package models

// User represents a user in the system
type User struct {
	ID        int    `json:"id"`
	UserName  string `json:"user_name" validate:"required,gt=1,lte=20"`
	FirstName string `json:"first_name" validate:"required,gt=1,lte=30"`
	LastName  string `json:"last_name"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password,omitempty" validate:"required"`
}

// LoginUserRequest represents a user in the system
type LoginUserRequest struct {
	UserName string `json:"user_name"`
	Password string `json:"password,omitempty" validate:"required"`
}

// LoginUserResponse represents a user in the system
type LoginUserResponse struct {
	AccessToken string `json:"access_token"`
}

type GetUsersParams struct {
	PageDetail
	UserName string `json:"user_name"`
}

type UpdateUserReq struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
