package db

import (
	"errors"
	"time"

	"github.com/gofiber/storage/postgres/v3"
	"github.com/google/uuid"

	"github.com/kylerequez/marketify/src/models"
	"github.com/kylerequez/marketify/src/shared"
)

type PostgresStorage struct {
	Name    string
	Config  shared.SQLConfig
	Storage *postgres.Storage
}

func NewPostgresStorage(name string, config shared.SQLConfig) *PostgresStorage {
	return &PostgresStorage{
		Name:   name,
		Config: config,
	}
}

func (storage *PostgresStorage) Init() error {
	if storage.Name == "" {
		return errors.New("name must not be empty")
	}

	if storage.Config.URI == "" {
		return errors.New("URI must not be empty")
	}

	storage.Storage = postgres.New(postgres.Config{
		Table:         storage.Name,
		ConnectionURI: storage.Config.URI,
	})

	return nil
}

func (storage *PostgresStorage) CreateNewSession(user models.User) (*shared.UserSession, error) {
	sessionId, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	session := shared.UserSession{
		ID:         sessionId.String(),
		Value:      []byte(user.ID.String()),
		Expiration: time.Now().Add(shared.USER_SESSION_LIFESPAN),
	}

	if err := storage.Storage.Set(
		session.ID,
		[]byte(session.Value),
		shared.USER_SESSION_LIFESPAN,
	); err != nil {
		return nil, err
	}

	return &session, nil
}

func (storage *PostgresStorage) GetSession(id uuid.UUID) (*shared.UserSession, error) {
	isAlive, err := storage.IsSessionAlive(id)
	if err != nil {
		return nil, err
	}

	if !isAlive {
		return nil, errors.New("session does not exists")
	}

	value, err := storage.Storage.Get(id.String())
	if err != nil {
		return nil, err
	}

	session := shared.UserSession{
		Value: value,
	}

	return &session, nil
}

func (storage *PostgresStorage) IsSessionAlive(id uuid.UUID) (bool, error) {
	_, err := storage.Storage.Get(id.String())
	if err != nil {
		return false, err
	}

	return true, nil
}

func (storage *PostgresStorage) RemoveSession(userId uuid.UUID) error {
	isAlive, err := storage.IsSessionAlive(userId)
	if err != nil {
		return err
	}

	if !isAlive {
		return errors.New("session does not exist")
	}

	if err := storage.Storage.Delete(userId.String()); err != nil {
		return err
	}

	return nil
}
