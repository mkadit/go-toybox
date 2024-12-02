package pgxsql

import (
	"github.com/jackc/pgx/v5"
	"github.com/mkadit/go-toybox/internal/models"
)

func (ad Adapter) CreateUser(user models.User) (err error) {
	query := `
	INSERT INTO users(
		email, username, password_hash 
	) VALUES (
		@email, @username, @password_hash
	)
	`

	_, err = ad.db.Exec(ad.ctx, query, pgx.NamedArgs{
		"email":         user.Email,
		"username":      user.Username,
		"password_hash": user.PasswordHash,
	},
	)
	if err != nil {
		return err
	}
	return nil
}

func (ad Adapter) GetUserByEmail(email string) (result models.User, err error) {
	query := `
	SELECT * FROM users WHERE email = @email
	`
	row, err := ad.db.Query(ad.ctx, query,
		pgx.NamedArgs{
			"email": email,
		},
	)
	result, err = pgx.CollectOneRow(row, pgx.RowToStructByName[models.User])
	if err != nil {
		return result, err
	}
	return result, nil
}
