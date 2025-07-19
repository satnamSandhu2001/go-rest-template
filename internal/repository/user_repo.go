package repository

import (
	"context"
	"database/sql"
	"errors"
	"go-rest-template/internal/db"
)

type UserRepository struct {
	*db.Queries
}

func NewUserRepository(q *db.Queries) *UserRepository {
	return &UserRepository{
		Queries: q,
	}
}

func (r *UserRepository) CreateNew(ctx context.Context, user db.CreateUserParams) (int32, error) {
	return r.Queries.CreateUser(ctx, user)
}

func (r *UserRepository) FindById(ctx context.Context, id int32) (db.User, error) {
	user, err := r.Queries.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.User{}, nil
		}
		return db.User{}, err
	}
	return user, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (db.User, error) {
	user, err := r.Queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.User{}, nil
		}
		return db.User{}, err
	}
	return user, nil
}
