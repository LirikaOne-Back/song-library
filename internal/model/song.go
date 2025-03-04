package model

import "time"

// Song представляет песню в библиотеке
type Song struct {
	ID          int64     `json:"id" db:"id"`
	Group       string    `json:"group" db:"group_name"`
	Song        string    `json:"song" db:"song_name"`
	ReleaseDate string    `json:"releaseDate" db:"release_date"`
	Text        string    `json:"text" db:"text"`
	Link        string    `json:"link" db:"link"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

// SongInput модель для добавления новой песни
type SongInput struct {
	Group string `json:"group" binding:"required"`
	Song  string `json:"song" binding:"required"`
}

// SongDetail ответ от внешнего API
type SongDetail struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

// SongFilter параметры фильтрации для списка песен
type SongFilter struct {
	Group    string
	SongName string
	Page     int
	PageSize int
}

// VersesPagination параметры пагинации для куплетов
type VersesPagination struct {
	Page     int
	PageSize int
}
