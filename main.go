package main

import (
	"eth-fetcher/app"
	"eth-fetcher/auth"
	"eth-fetcher/config"
	"eth-fetcher/database"
	"eth-fetcher/handlers"
	node "eth-fetcher/nodeconnect"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	err := godotenv.Load(".env")
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sLog := logger.Sugar()
	if err != nil {
		sLog.Fatal("error loading .env file")
	}

	cfg := config.LoadConfig()

	db, err := database.NewClient(cfg.DBConnectionURL)
	if err != nil {
		sLog.Fatalf("error creating db client: %v", err)
	}

	tg := node.NewNode(cfg.ETHNodeURL, sLog)

	app := app.NewApp(db, tg, sLog)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			app.Shutdown()
			os.Exit(1)
		}
	}()
	jwtAuth := auth.NewJWTAuth(cfg.JWT.Secret, cfg.JWT.Duration)
	handler := handlers.NewHTTP(app, cfg.APIPort, jwtAuth, sLog)
	handler.InitRoutes()
	handler.Run()
}
