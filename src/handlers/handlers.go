package handlers

import (
	"context"

	"github.com/gofiber/fiber/v3"

	"github.com/kylerequez/marketify/src/db"
	"github.com/kylerequez/marketify/src/middlewares"
	"github.com/kylerequez/marketify/src/repositories"
	"github.com/kylerequez/marketify/src/services"
	"github.com/kylerequez/marketify/src/shared"
	"github.com/kylerequez/marketify/src/utils"
)

func Init(app *fiber.App) error {
	config, err := utils.RetrieveSQLConfig()
	if err != nil {
		return err
	}

	database := db.NewPostgresDatabase(*config)
	ctx := context.TODO()
	if err := database.Open(ctx); err != nil {
		return err
	}

	if err := database.Ping(ctx); err != nil {
		return err
	}

	store := db.NewPostgresStorage("user_sessions", *config)
	if err := store.Init(); err != nil {
		return err
	}

	ur := repositories.NewUserRepository(database.Conn, shared.TABLES["USERS"])
	mh := middlewares.NewMiddlewareHandler(ur, store)
	us := services.NewUserService(ur, store)
	uh := NewUserHandler(app, us, mh)
	if err := uh.Init(); err != nil {
		return err
	}

	ah := NewAuthHandler(app, us)
	if err := ah.Init(); err != nil {
		return err
	}

	return nil
}
