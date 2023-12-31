package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rovn208/ross/pkg/configure"
	db "github.com/rovn208/ross/pkg/db/sqlc"
	"github.com/rovn208/ross/pkg/docs"
	"github.com/rovn208/ross/pkg/token"
	"github.com/rovn208/ross/pkg/youtube"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var v1RouterPrefix = "/api/v1"

type Server struct {
	config     configure.Config
	store      db.Querier
	router     *gin.Engine
	tokenMaker token.Maker
	ytClient   *youtube.Client
}

// NewServer returns a new server with the given configuration
func NewServer(config configure.Config, store db.Store, tokenMaker token.Maker) (*Server, error) {
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
		ytClient:   youtube.NewYoutubeClient(config),
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	v1 := router.Group(v1RouterPrefix)
	v1.Use(corsMiddleware())
	v1.StaticFS("/sources/", http.Dir(server.config.VideoDir))

	v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	docs.SwaggerInfo.BasePath = v1RouterPrefix

	usersRouter := v1.Group("/users")
	usersRouter.POST("/login", server.login)
	usersRouter.POST("/", server.createUser)
	usersRouter.Use(authMiddleware(server.tokenMaker))
	usersRouter.POST("/me", server.updateUser)
	usersRouter.GET("/me", server.getUser)
	usersRouter.GET("/:id", server.getUserByID)

	videosRouter := v1.Group("/videos")
	videosRouter.GET("/:id", server.getVideo)
	videosRouter.GET("/", server.getListVideo)
	videosRouter.Use(authMiddleware(server.tokenMaker))
	videosRouter.POST("/", server.createVideo)
	videosRouter.DELETE("/:id", server.deleteVideo)
	videosRouter.PUT("/:id", server.updateVideo)
	videosRouter.POST("/youtube", server.addYoutubeVideo)
	videosRouter.POST("/upload", server.uploadVideo)

	followsRouter := v1.Group("/follows")
	followsRouter.Use(authMiddleware(server.tokenMaker))
	followsRouter.GET("/followers", server.getListFollower)
	followsRouter.POST("/followers", server.followUser)
	followsRouter.DELETE("/followers", server.unfollowUser)
	followsRouter.GET("/following", server.getListFollowing)

	server.router = router
}

// Start starts the HTTP server
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func messageResponse(msg string) gin.H {
	return gin.H{"message": msg}
}
