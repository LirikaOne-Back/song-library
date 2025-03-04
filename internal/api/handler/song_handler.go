package handler

import (
	"context"
	"song-library/internal/model"
	"song-library/pkg/logger"
)

// SongService интерфейс сервиса песен
type SongService interface {
	CreateSong(ctx context.Context, input model.SongInput) (int64, error)
	GetSongs(ctx context.Context, filter model.SongFilter) ([]*model.Song, error)
	GetSongByID(ctx context.Context, id int64) (*model.Song, error)
	UpdateSong(ctx context.Context, song *model.Song) error
	DeleteSong(ctx context.Context, id int64) error
	GetSongVerses(ctx context.Context, id int64, pagination model.VersesPagination) ([]string, error)
}

// SongHandler обработчик HTTP запросов для работы с песнями
type SongHandler struct {
	service SongService
	logger  *logger.Logger
}

// NewSongHandler создает новый обработчик песен
func NewSongHandler(service SongService, logger *logger.Logger) *SongHandler {
	return &SongHandler{
		service: service,
		logger:  logger,
	}
}
