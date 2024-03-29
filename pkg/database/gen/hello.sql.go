// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: hello.sql

package gen

import (
	"context"
)

const get = `-- name: Get :many
select user_id, username, password, email from users WHERE username is not null
`

func (q *Queries) Get(ctx context.Context, db DBTX) ([]*User, error) {
	rows, err := db.Query(ctx, get)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*User{}
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.UserID,
			&i.Username,
			&i.Password,
			&i.Email,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
