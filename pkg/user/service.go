package user

import (
	"context"
	"fmt"
	"movie-booking-app/users-service/internal/database"
	"movie-booking-app/users-service/pkg/models"
	"movie-booking-app/users-service/pkg/utils"

	"github.com/jmoiron/sqlx"
)

//go:generate mockery --name Service --case underscore
type Service interface {
	CreateUser(context.Context, models.User) (models.User, error)
	LoginUser(context.Context, models.LoginUserRequest) (models.LoginUserResponse, error)
	// GetAllUsers(context.Context, models.GetUsersParams) ([]models.User, error)
	GetUser(context.Context, int) (models.User, error)
	// UpdateUser(context.Context, string, models.UpdateUserReq) (models.User, error)
	// DeleteUser(context.Context, string) error
}

type UserService struct {
	repo     Repository
	dbClient database.Database
}

func NewUserService(repo Repository, dbClient database.Database) Service {
	return &UserService{
		repo:     repo,
		dbClient: dbClient,
	}
}

func (us *UserService) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	currentUser, err := us.repo.GetUserByUsername(ctx, user.UserName)
	if err != nil || currentUser.ID != 0 {
		return models.User{}, fmt.Errorf("user name already exists")
	}
	currentUser, err = us.repo.GetUserByEmail(ctx, user.Email)
	if err != nil || currentUser.ID != 0 {
		return models.User{}, fmt.Errorf("user email already exists")
	}
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return models.User{}, err
	}
	user.Password = hashedPassword
	err = us.dbClient.WrapWithTransaction(ctx, func(tx *sqlx.Tx) error {
		id, err := us.repo.CreateUser(ctx, tx, user)
		if err != nil {
			return err
		}
		user.ID = id
		err = us.repo.CreateUserProfile(ctx, tx, user)
		if err != nil {
			return err
		}
		user.Password = ""
		return nil
	})
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (us *UserService) LoginUser(ctx context.Context, user models.LoginUserRequest) (models.LoginUserResponse, error) {
	currentUser, err := us.repo.GetUserByUsername(ctx, user.UserName)
	if err != nil {
		return models.LoginUserResponse{}, err
	}
	if currentUser.ID == 0 {
		return models.LoginUserResponse{}, fmt.Errorf("user does not exists")
	}

	if err := utils.CompareHashAndPassword(currentUser.Password, user.Password); err != nil {
		return models.LoginUserResponse{}, fmt.Errorf("username or password does not match")
	}

	token, err := utils.GenerateJWT(currentUser.ID)
	if err != nil {
		return models.LoginUserResponse{}, err
	}
	return models.LoginUserResponse{
		AccessToken: token,
	}, nil
}

func (us *UserService) GetUser(ctx context.Context, userId int) (models.User, error) {
	currentUser, err := us.repo.GetUserByID(ctx, userId)
	if err != nil {
		return models.User{}, err
	}
	if currentUser.ID == 0 {
		return models.User{}, fmt.Errorf("user does not exists")
	}
	return currentUser, nil
}

// func (us *UserService) CreateUser(ctx context.Context, user models.User) (string, error) {
// 	return "", nil
// }
