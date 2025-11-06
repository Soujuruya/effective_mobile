package v1

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"effective_mobile/internal/handlers"
	"effective_mobile/internal/service"
	middleware "effective_mobile/internal/transport"

	_ "effective_mobile/internal/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

var (
	defaultHeaderTimeout = time.Second * 5
)

type Server struct {
	srv  *http.Server
	Subs *handlers.SubscriptionHandler
}

func NewServer(port int, subsService *service.SubscriptionService) *Server {
	srv := http.Server{
		Addr:              ":" + strconv.Itoa(port),
		Handler:           nil,
		ReadHeaderTimeout: defaultHeaderTimeout,
	}
	return &Server{
		srv:  &srv,
		Subs: &handlers.SubscriptionHandler{Service: subsService},
	}
}

func (s *Server) RegisterHandlers() {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/subscriptions/create", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		s.Subs.Create(r.Context(), w, r)
	})

	mux.HandleFunc("/api/v1/subscriptions/list", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		s.Subs.List(w, r)
	})

	mux.HandleFunc("/api/v1/subscriptions/get", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		s.Subs.Get(r.Context(), w, r)
	})

	mux.HandleFunc("/api/v1/subscriptions/update", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		s.Subs.Update(r.Context(), w, r)
	})

	mux.HandleFunc("/api/v1/subscriptions/delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		s.Subs.Delete(r.Context(), w, r)
	})

	mux.HandleFunc("/api/v1/subscriptions/sum", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		s.Subs.SumPrice(r.Context(), w, r)
	})

	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	s.srv.Handler = middleware.LoggingMiddleware(mux)
}

func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
