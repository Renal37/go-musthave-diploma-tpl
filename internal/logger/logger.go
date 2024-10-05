package logger

import (
	"fmt"
	"go.uber.org/zap"
)

// Log глобальный логгер, инициализируется функцией Initialize.
// По умолчанию используется заглушка zap.NewNop(), которая не выводит никаких логов.
var Log *zap.Logger = zap.NewNop()

// Initialize инициализирует логгер с заданным уровнем логирования и средой выполнения.
// Параметры:
// - level: уровень логирования (например, "debug", "info", "warn", "error").
// - env: среда выполнения ("development" или "production").
func Initialize(level, env string) error {
	// Парсинг уровня логирования.
	logLevel, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return fmt.Errorf("ошибка парсинга уровня логирования: %w", err)
	}

	var config zap.Config

	// Выбор конфигурации логгера в зависимости от среды выполнения.
	if env == "development" {
		config = zap.NewDevelopmentConfig()
	} else {
		config = zap.NewProductionConfig()
	}

	// Установка уровня логирования.
	config.Level = logLevel

	// Построение логгера на основе конфигурации.
	logger, err := config.Build()
	if err != nil {
		return fmt.Errorf("ошибка построения логгера: %w", err)
	}

	// Присваиваем глобальной переменной Log инициализированный логгер.
	Log = logger

	return nil
}
