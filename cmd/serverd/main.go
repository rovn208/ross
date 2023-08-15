package main

import (
	"context"
	"github.com/rovn208/ross/pkg/util"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rovn208/ross/pkg/api"
	"github.com/rovn208/ross/pkg/configure"
	db "github.com/rovn208/ross/pkg/db/sqlc"
	"github.com/rovn208/ross/pkg/token"
)

// @title ROSS API
// @version 0.0.1
// @description Streaming service YouTube alike
//
// @contact.name	Ro Ngoc Vo
// @contact.url	github.com/rovn208
// @contact.email	ngocro208@gmail.com

// @host localhost:8080
// @BasePath /api/v1
func main() {
	util.Logger.Info("Starting Server")
	config, err := configure.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}
	connPool, _ := pgxpool.New(context.Background(), config.DBUrl)
	err = connPool.Ping(context.Background())
	if err != nil {
		log.Fatal("cannot when initialize db", err)
	}
	store := db.NewStore(connPool)
	tokenMaker, err := token.NewJWTMaker(config.TokenSecretKey)
	if err != nil {
		log.Fatal("error when initializing token maker", err)
	}
	server, err := api.NewServer(config, store, tokenMaker)
	if err != nil {
		log.Fatal("error when initializing server")
	}
	//runDBMigration(config.MigrationURL, config.DBUrl)

	go func() {
		if err = server.Start(config.HTTPServerAddress); err != nil {
			if err == http.ErrServerClosed {
				log.Fatal("Server closed under request")
			} else {
				log.Fatal("Server closed unexpect: ", err)
			}
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	util.Logger.Info("Shutting down server")
}
