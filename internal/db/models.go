// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package db

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	ID           int32            `json:"id"`
	Email        string           `json:"email"`
	PasswordHash string           `json:"-"`
	CreatedAt    pgtype.Timestamp `json:"created_at"`
	UpdatedAt    pgtype.Timestamp `json:"updated_at"`
}
