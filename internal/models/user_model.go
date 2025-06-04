package models

type Role string

const (
	RoleSuperAdmin Role = "superadmin"
)

type User struct {
	ID       int64  `db:"id" json:"id"`
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"-"`
	Role     Role   `db:"role" json:"role"`
}
