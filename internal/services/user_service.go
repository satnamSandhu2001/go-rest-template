package services

import (
	"context"
	"database/sql"
	"errors"
	"go-rest-template/internal/dto"
	"go-rest-template/internal/models"
	"go-rest-template/pkg"

	"github.com/jmoiron/sqlx"
)

type UserService struct {
	db *sqlx.DB
}

func NewUserService(db *sqlx.DB) *UserService {
	return &UserService{
		db: db,
	}
}

func (s *UserService) CreateUser(ctx context.Context, u *dto.User_RegisterRequest) error {
	hash, err := pkg.GenerateHash(u.Password)
	if err != nil {
		return err
	}
	u.Password = hash
	if u.Role == "" {
		u.Role = "superadmin"
	}
	query := `INSERT INTO users (email, password, role) VALUES (?, ?, ?)`
	res, err := s.db.ExecContext(ctx, query, u.Email, u.Password, u.Role)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	u.ID = id
	u.Password = ""
	return nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	err := s.db.GetContext(ctx, &u, "SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}
func (s *UserService) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	var u models.User
	err := s.db.GetContext(ctx, &u, "SELECT * FROM users WHERE id = ?", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (s *UserService) ListUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User
	err := s.db.SelectContext(ctx, &users, "SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *UserService) Authenticate(ctx context.Context, email string, password string) (*models.User, error) {
	var u models.User
	err := s.db.GetContext(ctx, &u, "SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.New("user not found")
	}

	err = pkg.CompareHashAndPassword(u.Password, password)
	if err != nil {
		return nil, errors.New("invalid password")
	}

	u.Password = ""
	return &u, nil
}
