package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"song-library/internal/model"
	"song-library/pkg/logger"
	"strconv"
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

// @Summary Получение списка песен
// @Description Получение списка песен с фильтрацией и пагинацией
// @Tags songs
// @Accept json
// @Produce json
// @Param group query string false "Фильтр по группе"
// @Param song query string false "Фильтр по названию песни"
// @Param page query int false "Номер страницы" default(1)
// @Param page_size query int false "Размер страницы" default(10)
// @Success 200 {array} model.Song
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /songs [get]
func (h *SongHandler) GetSongs(c *gin.Context) {
	log := h.logger.WithContext(c.Request.Context())

	log.Debug("Получение списка песен")

	filter := model.SongFilter{
		Group:    c.Query("group"),
		SongName: c.Query("song"),
		Page:     1,
		PageSize: 10,
	}

	if page, err := strconv.Atoi(c.Query("page")); err == nil && page > 0 {
		filter.Page = page
	}

	if pageSize, err := strconv.Atoi(c.Query("page_size")); err == nil && pageSize > 0 {
		filter.PageSize = pageSize
	}

	songs, err := h.service.GetSongs(c.Request.Context(), filter)
	if err != nil {
		log.Error("Ошибка получения списка песен", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Ошибка получения списка песен"})
		return
	}

	c.JSON(http.StatusOK, songs)
}

// @Summary Получение песни по ID
// @Description Получение данных конкретной песни по ID
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Success 200 {object} model.Song
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /songs/{id} [get]
func (h *SongHandler) GetSongByID(c *gin.Context) {
	log := h.logger.WithContext(c.Request.Context())
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Error("Неверный формат ID", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Неверный формат ID"})
		return
	}

	song, err := h.service.GetSongByID(c.Request.Context(), id)
	if err != nil {
		log.Error("Ошибка получения песни", "error", err, "id", id)
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Песня не найдена"})
		return
	}

	c.JSON(http.StatusOK, song)
}

// @Summary Создание новой песни
// @Description Добавление новой песни в библиотеку
// @Tags songs
// @Accept json
// @Produce json
// @Param input body model.SongInput true "Данные песни"
// @Success 201 {object} IdResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /songs [post]
func (h *SongHandler) CreateSong(c *gin.Context) {
	log := h.logger.WithContext(c.Request.Context())
	var input model.SongInput
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Error("Ошибка декодирования JSON", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Неверный формат данных"})
		return
	}

	id, err := h.service.CreateSong(c.Request.Context(), input)
	if err != nil {
		log.Error("Ошибка создания песни", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Ошибка создания песни"})
		return
	}

	c.JSON(http.StatusCreated, IdResponse{ID: id})
}

// @Summary Обновление песни
// @Description Обновление данных существующей песни
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Param input body model.Song true "Обновленные данные песни"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /songs/{id} [put]
func (h *SongHandler) UpdateSong(c *gin.Context) {
	log := h.logger.WithContext(c.Request.Context())
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Error("Неверный формат ID", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Неверный формат ID"})
		return
	}

	var song model.Song
	if err = c.ShouldBindJSON(&song); err != nil {
		log.Error("Ошибка декодирования JSON", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Неверный формат данных"})
		return
	}

	song.ID = id
	if err = h.service.UpdateSong(c.Request.Context(), &song); err != nil {
		log.Error("Ошибка обновления песни", "error", err, "id", id)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Ошибка обновления песни"})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Message: "Песня успешно обновлена"})
}

// @Summary Удаление песни
// @Description Удаление песни из библиотеки
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /songs/{id} [delete]
func (h *SongHandler) DeleteSong(c *gin.Context) {
	log := h.logger.WithContext(c.Request.Context())
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Error("Неверный формат ID", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Неверный формат ID"})
		return
	}

	if err = h.service.DeleteSong(c.Request.Context(), id); err != nil {
		log.Error("Ошибка удаления песни", "error", err, "id", id)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Ошибка удаления песни"})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Message: "Песня успешно удалена"})
}

// @Summary Получение текста песни по куплетам
// @Description Получение текста песни с пагинацией по куплетам
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Param page query int false "Номер страницы" default(1)
// @Param page_size query int false "Размер страницы" default(5)
// @Success 200 {object} VersesResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /songs/{id}/verses [get]
func (h *SongHandler) GetSongVerses(c *gin.Context) {
	log := h.logger.WithContext(c.Request.Context())
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Error("Неверный формат ID", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Неверный формат ID"})
		return
	}

	pagination := model.VersesPagination{
		Page:     1,
		PageSize: 5,
	}

	if page, err := strconv.Atoi(c.Query("page")); err == nil && page > 0 {
		pagination.Page = page
	}

	if pageSize, err := strconv.Atoi(c.Query("page_size")); err == nil && pageSize > 0 {
		pagination.PageSize = pageSize
	}

	verses, err := h.service.GetSongVerses(c.Request.Context(), id, pagination)
	if err != nil {
		log.Error("Ошибка получения куплетов песни", "error", err, "id", id)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Ошибка получения куплетов песни"})
		return
	}

	c.JSON(http.StatusOK, VersesResponse{Verses: verses})
}

// IdResponse ответ с идентификатором
type IdResponse struct {
	ID int64 `json:"id"`
}

// SuccessResponse ответ с сообщением об успехе
type SuccessResponse struct {
	Message string `json:"message"`
}

// ErrorResponse ответ с сообщением об ошибке
type ErrorResponse struct {
	Error string `json:"error"`
}

// VersesResponse ответ с куплетами песни
type VersesResponse struct {
	Verses []string `json:"verses"`
}
