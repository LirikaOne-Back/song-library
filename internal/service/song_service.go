package service

import (
	"context"
	"fmt"
	"song-library/internal/model"
	"song-library/pkg/logger"
)

// SongRepository интерфейс репозитория песен
type SongRepository interface {
	CreateSong(ctx context.Context, song *model.Song) (int64, error)
	GetSongs(ctx context.Context, filter model.SongFilter) ([]*model.Song, error)
	GetSongByID(ctx context.Context, id int64) (*model.Song, error)
	UpdateSong(ctx context.Context, song *model.Song) error
	DeleteSong(ctx context.Context, id int64) error
	GetSongVerses(ctx context.Context, id int64, pagination model.VersesPagination) ([]string, error)
}

// SongService сервис для работы с песнями
type SongService struct {
	repo      SongRepository
	apiClient *ExternalAPIClient
	logger    *logger.Logger
}

// NewSongService создает новый сервис для работы с песнями
func NewSongService(repo SongRepository, apiClient *ExternalAPIClient, logger *logger.Logger) *SongService {
	return &SongService{repo: repo, apiClient: apiClient, logger: logger}
}

// CreateSong создает новую песню
func (s *SongService) CreateSong(ctx context.Context, input model.SongInput) (int64, error) {
	log := s.logger.WithContext(ctx)

	log.Debug("Создание песни", "group", input.Group, "song", input.Song)

	details, err := s.apiClient.GetSongDetails(ctx, input.Group, input.Song)
	if err != nil {
		log.Error("Ошибка получения данных из внешнего API", "error", err)
		return 0, fmt.Errorf("ошибка получения данных песни: %w", err)
	}

	song := &model.Song{
		Group:       input.Group,
		Song:        input.Song,
		ReleaseDate: details.ReleaseDate,
		Text:        details.Text,
		Link:        details.Link,
	}

	id, err := s.repo.CreateSong(ctx, song)
	if err != nil {
		log.Error("Ошибка создания песни в репозитории", "error", err)
		return 0, fmt.Errorf("ошибка создания песни: %w", err)
	}

	log.Info("Песня успешно создана", "id", id)
	return id, nil
}

// GetSongs получает список песен с фильтрами
func (s *SongService) GetSongs(ctx context.Context, filter model.SongFilter) ([]*model.Song, error) {
	log := s.logger.WithContext(ctx)

	log.Debug("Получение списка песен с фильтром",
		"group", filter.Group,
		"song", filter.SongName,
		"page", filter.Page,
		"pageSize", filter.PageSize)

	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 10
	}

	songs, err := s.repo.GetSongs(ctx, filter)
	if err != nil {
		log.Error("Ошибка получения списка песен из репозитория", "error", err)
		return nil, fmt.Errorf("ошибка получения списка песен: %w", err)
	}

	log.Info("Список песен успешно получен", "count", len(songs))
	return songs, nil
}

// GetSongByID получает песню по идентификатору
func (s *SongService) GetSongByID(ctx context.Context, id int64) (*model.Song, error) {
	log := s.logger.WithContext(ctx)

	log.Debug("Получение песни по ID", "id", id)
	song, err := s.repo.GetSongByID(ctx, id)
	if err != nil {
		log.Error("Ошибка получения песни из репозитория", "error", err)
		return nil, fmt.Errorf("ошибка получения песни: %w", err)
	}

	if song == nil {
		log.Info("Песня не найдена", "id", id)
		return nil, fmt.Errorf("песня с id %d не найдена", id)
	}

	log.Info("Песня успешно получена", "id", id)
	return song, nil
}

// UpdateSong обновляет данные песни
func (s *SongService) UpdateSong(ctx context.Context, song *model.Song) error {
	log := s.logger.WithContext(ctx)

	log.Debug("Обновление песни", "id", song.ID)

	err := s.repo.UpdateSong(ctx, song)
	if err != nil {
		log.Error("Ошибка обновления песни в репозитории", "error", err)
		return fmt.Errorf("ошибка обновления песни: %w", err)
	}

	log.Info("Песня успешно обновлена", "id", song.ID)
	return nil
}

// DeleteSong удаляет песню
func (s *SongService) DeleteSong(ctx context.Context, id int64) error {
	log := s.logger.WithContext(ctx)

	log.Debug("Удаление песни", "id", id)

	err := s.repo.DeleteSong(ctx, id)
	if err != nil {
		log.Error("Ошибка удаления песни из репозитория", "error", err)
		return fmt.Errorf("ошибка удаления песни: %w", err)
	}

	log.Info("Песня успешно удалена", "id", id)
	return nil
}

// GetSongVerses получает куплеты песни с пагинацией
func (s *SongService) GetSongVerses(ctx context.Context, id int64, pagination model.VersesPagination) ([]string, error) {
	log := s.logger.WithContext(ctx)

	log.Debug("Получение куплетов песни", "id", id, "page", pagination.Page, "pageSize", pagination.PageSize)
	if pagination.Page <= 0 {
		pagination.Page = 1
	}
	if pagination.PageSize <= 0 {
		pagination.PageSize = 5
	}

	verses, err := s.repo.GetSongVerses(ctx, id, pagination)
	if err != nil {
		log.Error("Ошибка получения куплетов песни из репозитория", "error", err)
		return nil, fmt.Errorf("ошибка получения куплетов песни: %w", err)
	}

	log.Info("Куплеты песни успешно получены", "count", len(verses))
	return verses, nil
}
