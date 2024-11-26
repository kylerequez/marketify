package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/kylerequez/marketify/src/models"
)

type UserRepository struct {
	Conn  *pgx.Conn
	Table string
}

func NewUserRepository(conn *pgx.Conn, table string) *UserRepository {
	return &UserRepository{
		Conn:  conn,
		Table: table,
	}
}

func (ur *UserRepository) CreateUser(ctx context.Context, user models.User) error {
	sql := fmt.Sprintf(`
        INSERT INTO
            %s (
            firstname,
            middlename,
            lastname,
            email,
            password,
            authorities,
            status
        ) VALUES (
            $1,
            $2,
            $3,
            $4,
            $5,
            $6,
            $7
        );
        `, ur.Table)

	_, err := ur.Conn.Prepare(ctx, sql, sql)
	if err != nil {
		return err
	}

	res, err := ur.Conn.Exec(ctx, sql)
	if err != nil {
		return err
	}

	if count := res.RowsAffected(); count <= 0 {
		return errors.New("user was not created")
	}

	return nil
}

func (ur *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	sql := fmt.Sprintf(`
        SELECT
            id,
            firstname,
            middlename,
            lastname,
            email,
            password,
            authorities,
            status,
            created_at,
            updaated_at
        FROM
            %s
        WHERE
            email = $1
        LIMIT
            1;
        `, ur.Table)

	_, err := ur.Conn.Prepare(ctx, sql, sql)
	if err != nil {
		return nil, err
	}

	var user models.User
	row := ur.Conn.QueryRow(ctx, sql,
		email,
	)
	if err := row.Scan(
		&user.ID,
		&user.Firstname,
		&user.Middlename,
		&user.Lastname,
		&user.Email,
		&user.Password,
		&user.Authorities,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		if errors.Is(pgx.ErrNoRows, err) {
			return nil, errors.New("user does not exist")
		}
		return nil, err
	}

	return &user, nil
}
