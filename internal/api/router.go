package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"song-library/internal/api/handler"
	"song-library/pkg/logger"
)

// Router структура для маршрутизации API
type Router struct {
	engine      *gin.Engine
	songHandler *handler.SongHandler
	logger      *logger.Logger
}

// NewRouter создает и настраивает новый маршрутизатор
func NewRouter(songHandler *handler.SongHandler, log *logger.Logger, environment string) *Router {
	if environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(gin.Recovery())

	engine.Use(func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		ctx := context.WithValue(c.Request.Context(), "requestID", requestID)
		c.Request = c.Request.WithContext(ctx)

		c.Header("X-Request-ID", requestID)

		log.Info("HTTP запрос", "method", c.Request.Method, "path", c.Request.URL.Path, "requestID", requestID)
		c.Next()
	})

	return &Router{
		engine:      engine,
		songHandler: songHandler,
		logger:      log,
	}
}

// SetupRoutes настраивает все маршруты API
func (r *Router) SetupRoutes() {
	api := r.engine.Group("/api/v1")
	{
		songs := api.Group("/songs")
		{
			songs.GET("", r.songHandler.GetSongs)
			songs.POST("", r.songHandler.CreateSong)
			songs.GET("/:id", r.songHandler.GetSongByID)
			songs.PUT("/:id", r.songHandler.UpdateSong)
			songs.DELETE("/:id", r.songHandler.DeleteSong)
			songs.GET("/:id/verses", r.songHandler.GetSongVerses)
		}
	}

	r.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// GetEngine возвращает настроенный экземпляр gin.Engine
func (r *Router) GetEngine() *gin.Engine {
	return r.engine
}
