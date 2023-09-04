package server

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ethertun/agent-server/server/endpoints"
	"github.com/ethertun/agent-server/server/errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Callbacks struct {
	RunTask endpoints.RunTaskCallback
}

type AgentServer struct {
	router chi.Router
}

var (
	// HTTP Authentication header
	authentication = http.CanonicalHeaderKey("Authentication")
)

func setDefaultResponder() {
	slog.Info("setting default responder")
	render.Respond = func(w http.ResponseWriter, r *http.Request, v interface{}) {
		if _, ok := v.(*errors.ErrResponse); ok {
			resp := v.(*errors.ErrResponse)

			// log the error
			slog.Error("error occured during handler", "error", resp.ErrorText)
		}

		render.DefaultResponder(w, r, v)
	}
}

// Checks for the "Authentication" http header to be set to a bearer token
//
// Example:
// `Authentication: Bearer <token>`
func BearerAuth(key string) func(http.Handler) http.Handler {
	token := fmt.Sprintf("Bearer %s", key)
	f := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if token != r.Header.Get(authentication) {
                slog.Error("authentication failed: bad token")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}

	return f
}

func (a *AgentServer) Start(port int16) {
	slog.Info("starting agent server", "port", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), a.router)
}

func NewServer(authToken string, callbacks Callbacks) *AgentServer {
	setDefaultResponder()
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(BearerAuth(authToken))

	// all routes from here down require authentication
	r.Get("/healthz", endpoints.Healthz)
    r.Get("/capabilities", endpoints.Capabilities)
	r.Post("/task/run", endpoints.RunTask(callbacks.RunTask))

	return &AgentServer{router: r}
}
