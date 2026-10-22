package server

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/amanuel-tk/E-commerce/services/core/internal/httpx"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Deps struct {
	Logger *slog.Logger
	DB     *pgxpool.Pool
	RDB    *redis.Client
}

func NewRouter(d Deps) http.Handler {
	r := chi.NewRouter()

	r.Use(
		middleware.RequestID,
		middleware.RealIP,
		requestLogger(d.Logger),
		middleware.Recoverer,
		middleware.Timeout(30*time.Second),
	)

	r.Get("/healthz", healthHandler(d))

	return r
}

func requestLogger(l *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)

			status := ww.Status()
			if status == 0 {
				status = http.StatusOK
			}

			l.Info("request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", status,
				"duration", time.Since(start),
				"request_id", middleware.GetReqID(r.Context()),
			)
		})
	}
}

func healthHandler(d Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		status := map[string]string{"postgres": "ok", "redis": "ok"}
		code := http.StatusOK

		if err := d.DB.Ping(ctx); err != nil {
			status["postgres"], code = "down", http.StatusServiceUnavailable
		}

		if err := d.RDB.Ping(ctx).Err(); err != nil {
			status["redis"], code = "down", http.StatusServiceUnavailable
		}

		httpx.WriteJSON(w, code, status)
	}
}
