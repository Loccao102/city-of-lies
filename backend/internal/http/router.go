package http

import (
	"log/slog"

	"github.com/go-chi/chi/v5"

	"city-of-lies/backend/internal/config"
	"city-of-lies/backend/internal/http/handlers"
	"city-of-lies/backend/internal/http/middleware"
)

type RouterParams struct {
	Config        *config.Config
	Logger        *slog.Logger
	HealthHandler *handlers.HealthHandler
	Sessions      *handlers.SessionsHandler
	Agents        *handlers.AgentsHandler
	Investigation *handlers.InvestigationHandler
	WSHandler     *handlers.WSHandler
}

// NewRouter constructs the Chi HTTP mux with standard middlewares and endpoints.
func NewRouter(p RouterParams) *chi.Mux {
	r := chi.NewRouter()

	// Global Middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.Logging(p.Logger))
	r.Use(middleware.Recovery(p.Logger))
	r.Use(middleware.CORS(p.Config.CORS.AllowedOrigins))

	// Health Endpoints
	r.Get("/health/live", p.HealthHandler.Live)
	r.Get("/health/ready", p.HealthHandler.Ready)

	// WebSocket Endpoint
	r.Get("/ws", p.WSHandler.ServeHTTP)

	// REST API v1 Routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/sessions", func(r chi.Router) {
			r.Post("/", p.Sessions.CreateSession)
			r.Get("/{id}", p.Sessions.GetSession)

			r.Get("/{id}/agents", p.Agents.ListAgents)
			r.Get("/{id}/agents/{agentId}", p.Agents.GetAgent)
			r.Post("/{id}/agents/{agentId}/interviews", p.Agents.Interview)

			r.Post("/{id}/locations/{locationId}/inspect", p.Investigation.InspectLocation)
			r.Get("/{id}/notebook", p.Investigation.GetNotebook)
			r.Post("/{id}/corrections", p.Investigation.PublishCorrection)
			r.Post("/{id}/truth-submissions", p.Investigation.SubmitTruth)
		})
	})

	return r
}
