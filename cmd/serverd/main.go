package main

import (
	"context"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rovn208/ross/pkg/api"
	"github.com/rovn208/ross/pkg/configure"
	db "github.com/rovn208/ross/pkg/db/sqlc"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config, err := configure.LoadConfig(".")
	connPool, err := pgxpool.New(context.Background(), config.DBUrl)
	if err != nil {
		log.Fatal("cannot connect to db")
	}
	store := db.NewStore(connPool)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("error when initializing server")
	}
	runDBMigration(config.MigrationURL, config.DBUrl)

	//r.StaticFS("/", http.Dir(cfg.VideoDir))
	//ytClient := youtube.NewYoutubeClient()
	//err := ytClient.DownloadVideo("https://www.youtube.com/watch?v=9os5GBfuvJc")
	//if err != nil {
	//	log.Fatal(err)
	//}

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
