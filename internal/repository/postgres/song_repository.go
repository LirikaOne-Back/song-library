package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"song-library/internal/model"
	"song-library/pkg/logger"
	"strings"
	"time"
)

// SongRepository представляет репозиторий для работы с песнями в PostgreSQL
type SongRepository struct {
	db     *sqlx.DB
	logger *logger.Logger
}

// NewSongRepository создает новый репозиторий песен
func NewSongRepository(db *sqlx.DB, logger *logger.Logger) *SongRepository {
	return &SongRepository{
		db:     db,
		logger: logger,
	}
}

// NewPostgresDB устанавливает соединение с базой данных PostgreSQL
func NewPostgresDB(host, port, user, password, dbname string, logger *logger.Logger) (*sqlx.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	logger.Debug("Подключение к базе данных", "connection_string", connStr)
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к базе данных: %w", err)
	}

	logger.Info("Успешное подключение к базе данных")

	return db, nil
}

// CreateSong создает новую песню в базе данных
func (r *SongRepository) CreateSong(ctx context.Context, song *model.Song) (int64, error) {
	log := r.logger.WithContext(ctx)

	query := `INSERT INTO songs (group_name, song_name, release_date, text, link, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`

	log.Debug("Создание новой песни", "group", song.Group, "song", song.Song)

	now := time.Now()
	song.CreatedAt = now
	song.UpdatedAt = now

	var id int64
	err := r.db.QueryRowContext(
		ctx,
		query,
		song.Group,
		song.Song,
		song.ReleaseDate,
		song.Text,
		song.Link,
		song.CreatedAt,
		song.UpdatedAt,
	).Scan(&id)
	if err != nil {
		log.Error("Ошибка создания песни", "error", err)
		return 0, fmt.Errorf("ошибка создания песни: %w", err)
	}

	log.Info("Песня успешно создана", "id", id)
	return id, nil
}

// GetSongs получает список песен с фильтрацией и пагинацией
func (r *SongRepository) GetSongs(ctx context.Context, filter model.SongFilter) ([]*model.Song, error) {
	log := r.logger.WithContext(ctx)

	log.Debug("Получение списка песен с фильтром",
		"group", filter.Group,
		"song", filter.SongName,
		"page", filter.Page,
		"pageSize", filter.PageSize)

	query := `SELECT id, group_name, song_name, release_date, text, link, created_at, updated_at 
		FROM songs WHERE 1=1`
	params := []interface{}{}
	paramCount := 1

	if filter.Group != "" {
		query += fmt.Sprintf(" AND group_name ILIKE $%d", paramCount)
		params = append(params, "%"+filter.Group+"%")
		paramCount++
	}

	if filter.SongName != "" {
		query += fmt.Sprintf(" AND song_name ILIKE $%d", paramCount)
		params = append(params, "%"+filter.SongName+"%")
		paramCount++
	}

	offset := (filter.Page - 1) * filter.PageSize
	query += fmt.Sprintf(" ORDER BY id DESC LIMIT $%d OFFSET $%d", paramCount, paramCount+1)
	params = append(params, filter.PageSize, offset)

	log.Debug("Выполнение запроса", "query", query, "params", params)

	rows, err := r.db.QueryxContext(ctx, query, params...)
	if err != nil {
		log.Error("Ошибка получения списка песен", "error", err)
		return nil, fmt.Errorf("ошибка получения списка песен: %w", err)
	}
	defer rows.Close()

	var songs []*model.Song
	for rows.Next() {
		var song model.Song
		if err = rows.StructScan(&song); err != nil {
			log.Error("Ошибка сканирования песни", "error", err)
			return nil, fmt.Errorf("ошибка сканирования песни: %w", err)
		}
		songs = append(songs, &song)
	}

	log.Info("Успешно получен список песен", "count", len(songs))
	return songs, nil
}

// GetSongByID получает песню по идентификатору
func (r *SongRepository) GetSongByID(ctx context.Context, id int64) (*model.Song, error) {
	log := r.logger.WithContext(ctx)

	log.Debug("Получение песни по ID", "id", id)

	query := `SELECT id, group_name, song_name, release_date, text, link, created_at, updated_at FROM songs WHERE id = $1`

	var song model.Song
	err := r.db.GetContext(ctx, &song, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Info("Песня не найдена", "id", id)
			return nil, nil
		}
		log.Error("Ошибка получения песни", "error", err)
		return nil, fmt.Errorf("ошибка получения песни: %w", err)
	}

	log.Info("Песня успешно получена", "id", id)
	return &song, nil
}

// UpdateSong обновляет данные песни
func (r *SongRepository) UpdateSong(ctx context.Context, song *model.Song) error {
	log := r.logger.WithContext(ctx)

	log.Debug("Обновление песни", "id", song.ID)

	query := `UPDATE songs SET group_name = $1, song_name = $2, release_date = $3, text = $4, link = $5, updated_at = $6 WHERE id = $7`

	song.UpdatedAt = time.Now()
	result, err := r.db.ExecContext(
		ctx,
		query,
		song.Group,
		song.Song,
		song.ReleaseDate,
		song.Text,
		song.Link,
		song.UpdatedAt,
		song.ID,
	)

	if err != nil {
		log.Error("Ошибка обновления песни", "error", err)
		return fmt.Errorf("ошибка обновления песни: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error("Ошибка получения количества затронутых строк", "error", err)
		return fmt.Errorf("ошибка получения количества затронутых строк: %w", err)
	}

	if rowsAffected == 0 {
		log.Info("Песня для обновления не найдена", "id", song.ID)
		return fmt.Errorf("песня с id %d не найдена", song.ID)
	}

	log.Info("Песня успешно обновлена", "id", song.ID)
	return nil
}

// DeleteSong удаляет песню из базы данных
func (r *SongRepository) DeleteSong(ctx context.Context, id int64) error {
	log := r.logger.WithContext(ctx)

	log.Debug("Удаление песни", "id", id)

	query := `DELETE FROM songs WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		log.Error("Ошибка удаления песни", "error", err)
		return fmt.Errorf("ошибка удаления песни: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error("Ошибка получения количества затронутых строк", "error", err)
		return fmt.Errorf("ошибка получения количества затронутых строк: %w", err)
	}
	if rowsAffected == 0 {
		log.Info("Песня для удаления не найдена", "id", id)
		return fmt.Errorf("песня с id %d не найдена", id)
	}

	log.Info("Песня успешно удалена", "id", id)
	return nil
}

// GetSongVerses получает куплеты песни с пагинацией
func (r *SongRepository) GetSongVerses(ctx context.Context, id int64, pagination model.VersesPagination) ([]string, error) {
	log := r.logger.WithContext(ctx)

	log.Debug("Получение куплетов песни", "id", id, "page", pagination.Page, "pageSize", pagination.PageSize)

	song, err := r.GetSongByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if song == nil {
		log.Info("Песня не найдена", "id", id)
		return nil, fmt.Errorf("песня с id %d не найдена", id)
	}

	verses := strings.Split(song.Text, "\n\n")
	start := (pagination.Page - 1) * pagination.PageSize
	end := start + pagination.PageSize
	if start >= len(verses) {
		log.Info("Пагинация выходит за пределы", "verses_count", len(verses), "start", start)
		return []string{}, nil
	}

	if end > len(verses) {
		end = len(verses)
	}

	log.Info("Успешно получены куплеты песни", "verses_count", len(verses[start:end]))
	return verses[start:end], nil
}
