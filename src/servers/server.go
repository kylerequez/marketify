package servers

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"

	"github.com/kylerequez/marketify/src/handlers"
	"github.com/kylerequez/marketify/src/shared"
	"github.com/kylerequez/marketify/src/utils"
)

type Server interface {
	Run() error
}

type MarketifyServer struct {
	Config shared.ServerConfig
	App    fiber.App
}

func NewMarketifyServer(config shared.ServerConfig) *MarketifyServer {
	return &MarketifyServer{
		Config: config,
	}
}

func Init() error {
	if err := utils.LoadEnv(); err != nil {
		return err
	}

	config, err := utils.RetrieveServerConfig()
	if err != nil {
		return err
	}

	server := NewMarketifyServer(*config)
	return server.Run()
}

func (server *MarketifyServer) Run() error {
	app := fiber.New(fiber.Config{
		AppName: server.Config.AppName,
	})

	app.Use("/styles", static.New("./src/public/"))

	if err := handlers.Init(app); err != nil {
		return err
	}

	return app.Listen(
		fmt.Sprintf(
			"%s:%s",
			server.Config.Hostname,
			server.Config.Port,
		),
	)
}
