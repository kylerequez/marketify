package repositories

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/kylerequez/marketify/src/models"
	"github.com/kylerequez/marketify/src/shared"
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
			age,
			gender,
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
            $7,
			$8,
			$9
        );
        `, ur.Table)

	_, err := ur.Conn.Prepare(ctx, sql, sql)
	if err != nil {
		return err
	}

	res, err := ur.Conn.Exec(ctx, sql,
		user.Firstname,
		user.Middlename,
		user.Lastname,
		user.Age,
		user.Gender,
		user.Email,
		user.Password,
		user.Authorities,
		user.Status,
	)
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
            updated_at
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

func (ur *UserRepository) GetUserById(ctx context.Context, id uuid.UUID) (*models.User, error) {
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
            updated_at
        FROM
            %s
        WHERE
            id = $1
        LIMIT
            1;
        `, ur.Table)

	_, err := ur.Conn.Prepare(ctx, sql, sql)
	if err != nil {
		return nil, err
	}

	var user models.User
	row := ur.Conn.QueryRow(ctx, sql,
		id,
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

func (ur *UserRepository) GetAllUsers(ctx context.Context) (*[]models.User, error) {
	sql := fmt.Sprintf(`
        SELECT
            id,
            firstname,
            middlename,
            lastname,
            email,
            authorities,
            status,
            created_at,
            updated_at
        FROM
            %s
		LIMIT
			%d;
        `, ur.Table, shared.MAX_USERS_PER_PAGE)

	_, err := ur.Conn.Prepare(ctx, sql, sql)
	if err != nil {
		return nil, err
	}

	var users []models.User
	rows, err := ur.Conn.Query(ctx, sql)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var user models.User

		if err := rows.Scan(
			&user.ID,
			&user.Firstname,
			&user.Middlename,
			&user.Lastname,
			&user.Email,
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

		users = append(users, user)
	}

	return &users, nil
}

func (ur *UserRepository) IsRole(ctx context.Context, id uuid.UUID, role string) (bool, error) {
	sql := fmt.Sprintf(`
		SELECT
			id,
			authorities
		FROM
			%s
		WHERE
			id = $1
		LIMIT
			1;
		`, ur.Table)

	_, err := ur.Conn.Prepare(ctx, sql, sql)
	if err != nil {
		return false, err
	}

	var user models.User
	row := ur.Conn.QueryRow(ctx, sql,
		id,
	)
	if err := row.Scan(
		&user.ID,
		&user.Authorities,
	); err != nil {
		return false, nil
	}

	return slices.Contains(user.Authorities, role), nil
}
