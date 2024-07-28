package user

import (
	"context"
	"database/sql"
	"movie-booking-app/users-service/internal/database"
	"movie-booking-app/users-service/pkg/models"

	"github.com/jmoiron/sqlx"
)

//go:generate mockery --name Repository --case underscore
type Repository interface {
	CreateUser(context.Context, *sqlx.Tx, models.User) (int, error)
	CreateUserProfile(context.Context, *sqlx.Tx, models.User) error
	GetUserByID(context.Context, int) (models.User, error)
	GetUserByUsername(context.Context, string) (models.User, error)
	GetUserByEmail(context.Context, string) (models.User, error)
	// GetAllUsers(context.Context, models.GetUsersParams) ([]models.User, error)
	// UpdateUser(context.Context, string, models.UpdateUserReq) (models.User, error)
	// DeleteUser(context.Context, string) error
}

type UserRepo struct {
	dbClient database.Database
}

func NewUserRepo(dbCli database.Database) Repository {
	return &UserRepo{
		dbClient: dbCli,
	}
}

func (ur *UserRepo) CreateUser(ctx context.Context, tx *sqlx.Tx, user models.User) (int, error) {
	query := `INSERT INTO users (user_name, email, password_hash) VALUES (?,?,?)`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return 0, err
	}
	res, err := stmt.ExecContext(ctx, user.UserName, user.Email, user.Password)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (ur *UserRepo) CreateUserProfile(ctx context.Context, tx *sqlx.Tx, user models.User) error {
	query := `INSERT INTO user_profiles (user_id, first_name, last_name) VALUES (?,?,?)`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, user.ID, user.FirstName, user.LastName)
	if err != nil {
		return err
	}
	return nil
}

func (ur *UserRepo) GetUserByID(ctx context.Context, userId int) (models.User, error) {
	var user models.User
	query := "SELECT u.id, user_name, first_name, last_name, email FROM users u JOIN user_profiles up ON u.id=up.user_id WHERE u.id = ?"
	stmt, err := ur.dbClient.GetDB().PrepareContext(ctx, query)
	if err != nil {
		return models.User{}, err
	}
	row := stmt.QueryRowContext(ctx, userId)
	err = row.Scan(&user.ID, &user.UserName, &user.FirstName, &user.LastName, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, nil
		}
		return models.User{}, err
	}
	return user, nil
}

func (ur *UserRepo) GetUserByUsername(ctx context.Context, userName string) (models.User, error) {
	var user models.User
	query := "SELECT id, user_name, password_hash, email from users WHERE user_name = ?"
	stmt, err := ur.dbClient.GetDB().PrepareContext(ctx, query)
	if err != nil {
		return models.User{}, err
	}
	row := stmt.QueryRowContext(ctx, userName)
	err = row.Scan(&user.ID, &user.UserName, &user.Password, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, nil
		}
		return models.User{}, err
	}
	return user, nil
}

func (ur *UserRepo) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User
	query := "SELECT id, user_name, email from users WHERE email = ?"
	stmt, err := ur.dbClient.GetDB().PrepareContext(ctx, query)
	if err != nil {
		return models.User{}, err
	}
	row := stmt.QueryRowContext(ctx, email)
	err = row.Scan(&user.ID, &user.UserName, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, nil
		}
		return models.User{}, err
	}
	return user, nil
}
