package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rovn208/ross/pkg/api"
	"github.com/rovn208/ross/pkg/configure"
	db "github.com/rovn208/ross/pkg/db/sqlc"
	"github.com/rovn208/ross/pkg/token"
)

func main() {
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
				log.Println("Server closed under request")
			} else {
				log.Fatal("Server closed unexpect: ", err)
			}
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server")
}

func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal("cannot create new migrate instance")
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("failed to run migrate up")
	}

	log.Println("db migrated successfully")
}
