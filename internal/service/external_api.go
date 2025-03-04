package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"song-library/internal/model"
	"song-library/pkg/logger"
	"time"
)

// ExternalAPIClient клиент для работы с внешним API
type ExternalAPIClient struct {
	baseURL string
	client  *http.Client
	logger  *logger.Logger
}

// NewExternalAPIClient создает новый клиент внешнего API
func NewExternalAPIClient(baseURL string, logger *logger.Logger) *ExternalAPIClient {
	return &ExternalAPIClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		logger: logger,
	}
}

// GetSongDetails получает детали песни из внешнего API
func (c *ExternalAPIClient) GetSongDetails(ctx context.Context, group, song string) (*model.SongDetail, error) {
	log := c.logger.WithContext(ctx)

	log.Debug("Получение деталей песни из внешнего API", "group", group, "song", song)

	u, err := url.Parse(c.baseURL + "/info")
	if err != nil {
		log.Error("Ошибка при формировании URL", "error", err)
		return nil, fmt.Errorf("ошибка при формировании URL: %w", err)
	}

	q := u.Query()
	q.Set("group", group)
	q.Set("song", song)
	u.RawQuery = q.Encode()

	log.Debug("Отправка запроса к внешнему API", "url", u.String())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		log.Error("Ошибка создания запроса", "error", err)
		return nil, fmt.Errorf("ошибка создания запроса: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		log.Error("Ошибка выполнения запроса", "error", err)
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Error("Внешний API вернул ошибку", "status_code", resp.StatusCode)
		return nil, fmt.Errorf("внешний API вернул код состояния %d", resp.StatusCode)
	}

	var songDetail model.SongDetail
	if err = json.NewDecoder(resp.Body).Decode(&songDetail); err != nil {
		log.Error("Ошибка декодирования ответа", "error", err)
		return nil, fmt.Errorf("ошибка декодирования ответа: %w", err)
	}

	log.Info("Успешно получены детали песни из внешнего API")
	return &songDetail, nil
}
