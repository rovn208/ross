package api

import (
	"github.com/gin-gonic/gin"
	"github.com/rovn208/ross/pkg/configure"
	db "github.com/rovn208/ross/pkg/db/sqlc"
)

type Server struct {
	config configure.Config
	store  db.Querier
	router *gin.Engine
}

func NewServer(config configure.Config, store db.Store) (*Server, error) {
	server := &Server{
		config: config,
		store:  store,
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	v1 := router.Group("/api/v1")

	usersRouter := v1.Group("/users")
	usersRouter.POST("/", server.createUser)
	usersRouter.POST("/me", server.updateUser)
	usersRouter.GET("/:id", server.getUser)

	videosRouter := v1.Group("/videos")
	videosRouter.POST("/", server.createVideo)
	videosRouter.DELETE("/:id", server.deleteVideo)
	videosRouter.GET("/:id", server.getVideo)
	videosRouter.PUT("/", server.updateVideo)

	followsRouter := v1.Group("/follows")
	followsRouter.GET("/followers", server.getListFollower)
	followsRouter.POST("/followers", server.followUser)
	followsRouter.DELETE("/followers", server.unfollowUser)
	followsRouter.GET("/following", server.getListFollowing)

	server.router = router
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func messageResponse(msg string) gin.H {
	return gin.H{"message": msg}
}
