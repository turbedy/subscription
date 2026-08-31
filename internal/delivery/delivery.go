package delivery

import (
	"errors"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/turbedy/subscription/internal/health"
)

type Server struct {
	*health.Handler
}

func New(
	log *slog.Logger,
	health *health.Handler,
) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.ClientIPFromRemoteAddr)
	router.Use(Logger(log))
	router.Use(Recoverer(log))
	HandlerFromMux(&Server{health}, router)
	return router
}

func Logger(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			st := time.Now()
			next.ServeHTTP(ww, r)

			log.Info(
				"request completed",
				"id", middleware.GetReqID(r.Context()),
				"client", middleware.GetClientIP(r.Context()),
				"method", r.Method,
				"path", r.URL.Path,
				"status", ww.Status(),
				"bytes", ww.BytesWritten(),
				"duration", time.Since(st),
			)
		})
	}
}

func Recoverer(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				rvr := recover()
				if rvr == nil {
					return
				}

				if err, ok := rvr.(error); ok && errors.Is(err, http.ErrAbortHandler) {
					panic(rvr)
				}

				log.Error(
					"panic recovered",
					"id", middleware.GetReqID(r.Context()),
					"panic", rvr,
					"stack", string(debug.Stack()),
				)
			}()

			next.ServeHTTP(w, r)
		})
	}
}
