package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"city-of-lies/backend/internal/application/services"
	"city-of-lies/backend/internal/config"
	"city-of-lies/backend/internal/game/dialogue"
	"city-of-lies/backend/internal/game/scenario"
	"city-of-lies/backend/internal/game/simulation"
	appHTTP "city-of-lies/backend/internal/http"
	"city-of-lies/backend/internal/http/handlers"
	"city-of-lies/backend/internal/infrastructure/realtime"
)

func main() {
	// 1. Setup Structured Logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	logger.Info("starting City of Lies API server...")

	// 2. Load Configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Error("configuration error", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// 3. Load All Scenarios
	scenariosDir := cfg.ScenarioDir
	if _, err := os.Stat(filepath.Join(scenariosDir, "riverside-factory")); os.IsNotExist(err) {
		scenariosDir = filepath.Join("..", "data", "scenarios")
	}
	if _, err := os.Stat(filepath.Join(scenariosDir, "riverside-factory")); os.IsNotExist(err) {
		scenariosDir = filepath.Join("data", "scenarios")
	}

	logger.Info("loading all scenario bundles", slog.String("dir", scenariosDir))
	bundles, err := scenario.LoadAllScenarios(scenariosDir)
	if err != nil {
		logger.Error("failed to load scenarios", slog.String("error", err.Error()))
		os.Exit(1)
	}

	defaultBundle := bundles["riverside-factory"]
	if defaultBundle == nil {
		for _, b := range bundles {
			defaultBundle = b
			break
		}
	}
	logger.Info("loaded scenario bundles successfully", slog.Int("count", len(bundles)))

	// 4. Initialize Simulation Session Manager
	sessionManager := simulation.NewMultiScenarioSessionManager(bundles, defaultBundle, cfg)

	// 5. Initialize Dialogue Provider
	var dialProvider *dialogue.TemplateProvider
	dialProvider = dialogue.NewTemplateProvider()

	// 6. Initialize Realtime WebSocket Hub
	hub := realtime.NewHub()

	// 7. Initialize Application Services
	startSessionSvc := services.NewStartSessionService(sessionManager)
	listScenariosSvc := services.NewListScenariosService(sessionManager)
	getStateSvc := services.NewGetStateService(sessionManager)
	interviewSvc := services.NewInterviewService(sessionManager, dialProvider)
	inspectSvc := services.NewInspectService(sessionManager)
	notebookSvc := services.NewNotebookService(sessionManager)
	correctionSvc := services.NewPublishCorrectionService(sessionManager)
	truthSvc := services.NewSubmitTruthService(sessionManager)

	// 8. Initialize HTTP Handlers
	healthHandler := handlers.NewHealthHandler(func() bool { return true })
	sessionsHandler := handlers.NewSessionsHandler(startSessionSvc, getStateSvc)
	scenariosHandler := handlers.NewScenariosHandler(listScenariosSvc)
	agentsHandler := handlers.NewAgentsHandler(sessionManager, interviewSvc)
	investigationHandler := handlers.NewInvestigationHandler(inspectSvc, notebookSvc, correctionSvc, truthSvc)
	wsHandler := handlers.NewWSHandler(hub, sessionManager)

	// 9. Build Router
	router := appHTTP.NewRouter(appHTTP.RouterParams{
		Config:        cfg,
		Logger:        logger,
		HealthHandler: healthHandler,
		Sessions:      sessionsHandler,
		Scenarios:     scenariosHandler,
		Agents:        agentsHandler,
		Investigation: investigationHandler,
		WSHandler:     wsHandler,
	})

	// 10. Start HTTP Server with Graceful Shutdown
	srv := &http.Server{
		Addr:         cfg.HTTP.Addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info("HTTP server listening", slog.String("addr", cfg.HTTP.Addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server failure", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	// Wait for termination signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	logger.Info("shutting down server...", slog.String("signal", sig.String()))

	sessionManager.Shutdown()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("server forced to shutdown", slog.String("error", err.Error()))
	}

	logger.Info("server exited cleanly")
}
