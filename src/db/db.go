package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/kylerequez/marketify/src/shared"
)

type DatabaseConnection interface {
	Open(context.Context) error
	Close(context.Context) error
	Ping(context.Context) error
}

type PostgresDatabase struct {
	Config shared.SQLConfig
	Conn   *pgx.Conn
}

func NewPostgresDatabase(config shared.SQLConfig) *PostgresDatabase {
	return &PostgresDatabase{
		Config: config,
	}
}

func (database *PostgresDatabase) Open(ctx context.Context) error {
	if database.Conn != nil {
		return errors.New("there is already an existing database connection")
	}

	conn, err := pgx.Connect(ctx, database.Config.URI)
	if err != nil {
		return err
	}

	database.Conn = conn
	return nil
}

func (database *PostgresDatabase) Close(ctx context.Context) error {
	if database.Conn == nil {
		return errors.New("there is no database connection")
	}

	if err := database.Conn.Close(ctx); err != nil {
		return err
	}

	return nil
}

func (database *PostgresDatabase) Ping(ctx context.Context) error {
	if database.Conn == nil {
		return errors.New("there is no database connection")
	}

	if err := database.Conn.Ping(ctx); err != nil {
		return err
	}

	return nil
}
