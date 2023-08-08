package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rovn208/ross/pkg/youtube"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	const musicDir = "videos"
	const port = 8080
	router := gin.Default()
	router.StaticFS("/", http.Dir(musicDir))
	ytClient := youtube.NewYoutubeClient()
	err := ytClient.DownloadVideo("https://www.youtube.com/watch?v=9os5GBfuvJc")
	if err != nil {
		log.Fatal(err)
	}
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: router,
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

func fileServerHandler(dir string) http.HandlerFunc {
	musicDir := http.Dir(dir)

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		http.FileServer(musicDir).ServeHTTP(w, r)
	}
}
