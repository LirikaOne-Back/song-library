package migration

import (
	"database/sql"
	"fmt"
	"song-library/pkg/logger"
)

// Миграционные SQL-запросы
var migrations = []string{
	`CREATE TABLE IF NOT EXISTS songs (
		id SERIAL PRIMARY KEY,
		group_name VARCHAR(255) NOT NULL,
		song_name VARCHAR(255) NOT NULL,
		release_date VARCHAR(50) NOT NULL,
		text TEXT NOT NULL,
		link VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL,
		CONSTRAINT unique_group_song UNIQUE (group_name, song_name)
	);`,
}

// RunMigrations выполняет все миграции базы данных
func RunMigrations(db *sql.DB, logger *logger.Logger) error {
	logger.Info("Запуск миграций базы данных")

	for i, migration := range migrations {
		logger.Debug("Выполнение миграции", "index", i)

		_, err := db.Exec(migration)
		if err != nil {
			logger.Error("Ошибка выполнения миграции", "index", i, "error", err)
			return fmt.Errorf("ошибка выполнения миграции %d: %w", i, err)
		}

		logger.Debug("Миграция успешно выполнена", "index", i)
	}
	
	logger.Info("Все миграции успешно выполнены")
	return nil
}
