// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package gen

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	UserID   pgtype.UUID `db:"user_id" json:"user_id"`
	Username string      `db:"username" json:"username"`
	Password string      `db:"password" json:"password"`
	Email    string      `db:"email" json:"email"`
}
