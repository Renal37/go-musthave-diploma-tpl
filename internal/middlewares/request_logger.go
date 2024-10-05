package middlewares

import (
	"net/http"
	"time"

	"go.uber.org/zap"
	"github.com/Renal37/go-musthave-diploma-tpl/internal/logger" // Импортируем пакет logger
)

// RequestLogger является middleware, которое логирует информацию о каждом HTTP-запросе.
// Логируются URI, метод запроса, длительность обработки и код статуса ответа.
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		wrappedWriter := newResponseWriter(w)

		// Обработка запроса следующим обработчиком.
		next.ServeHTTP(wrappedWriter, r)

		duration := time.Since(startTime)

		// Логирование информации о запросе.
		logger.Log.Info("Запрос обработан",
			zap.String("URI", r.RequestURI),
			zap.String("метод", r.Method),
			zap.Duration("длительность", duration),
			zap.Int("статус", wrappedWriter.statusCode),
		)
	})
}

// responseWriter оборачивает http.ResponseWriter и сохраняет код статуса ответа.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// newResponseWriter создает новый экземпляр responseWriter с кодом статуса по умолчанию (200 OK).
func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

// WriteHeader сохраняет код статуса ответа и вызывает метод WriteHeader у встроенного ResponseWriter.
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
