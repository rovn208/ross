package main

import (
	"context"
	"fmt"
	"github.com/rovn208/ross/pkg/config"
	"github.com/rovn208/ross/pkg/router"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.LoadConfig(".")
	r := router.NewRouter()
	//r.StaticFS("/", http.Dir(cfg.VideoDir))
	//ytClient := youtube.NewYoutubeClient()
	//err := ytClient.DownloadVideo("https://www.youtube.com/watch?v=9os5GBfuvJc")
	//if err != nil {
	//	log.Fatal(err)
	//}
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", cfg.Port),
		Handler: r,
	}
	go func() {
		if err = srv.ListenAndServe(); err != nil {
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}
	log.Println("Server exiting")

}
