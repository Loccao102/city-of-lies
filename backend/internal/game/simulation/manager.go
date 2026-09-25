package simulation

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"city-of-lies/backend/internal/config"
	"city-of-lies/backend/internal/domain"
	"city-of-lies/backend/internal/game/scenario"
)

type SessionManager struct {
	mu      sync.RWMutex
	engines map[string]*SimulationEngine
	cancels map[string]context.CancelFunc
	bundle  *scenario.ScenarioBundle
	bundles map[string]*scenario.ScenarioBundle
	cfg     *config.Config
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewSessionManager(bundle *scenario.ScenarioBundle, cfg *config.Config) *SessionManager {
	ctx, cancel := context.WithCancel(context.Background())
	bundles := make(map[string]*scenario.ScenarioBundle)
	if bundle != nil {
		bundles[bundle.Metadata.ID] = bundle
	}
	return &SessionManager{
		engines: make(map[string]*SimulationEngine),
		cancels: make(map[string]context.CancelFunc),
		bundle:  bundle,
		bundles: bundles,
		cfg:     cfg,
		ctx:     ctx,
		cancel:  cancel,
	}
}

func NewMultiScenarioSessionManager(bundles map[string]*scenario.ScenarioBundle, defaultBundle *scenario.ScenarioBundle, cfg *config.Config) *SessionManager {
	ctx, cancel := context.WithCancel(context.Background())
	return &SessionManager{
		engines: make(map[string]*SimulationEngine),
		cancels: make(map[string]context.CancelFunc),
		bundle:  defaultBundle,
		bundles: bundles,
		cfg:     cfg,
		ctx:     ctx,
		cancel:  cancel,
	}
}

// GetBundles returns all registered scenario bundles.
func (m *SessionManager) GetBundles() map[string]*scenario.ScenarioBundle {
	m.mu.RLock()
	defer m.mu.RUnlock()
	res := make(map[string]*scenario.ScenarioBundle, len(m.bundles))
	for k, v := range m.bundles {
		res[k] = v
	}
	return res
}

// CreateSession initializes and registers a new simulation session with the requested or random scenario.
func (m *SessionManager) CreateSession(scenarioID string, seed int64) (*SimulationEngine, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if seed == 0 {
		seed = time.Now().UnixNano()
	}

	targetBundle := m.bundle
	if (scenarioID == "" || scenarioID == "random") && len(m.bundles) > 0 {
		keys := make([]string, 0, len(m.bundles))
		for k := range m.bundles {
			keys = append(keys, k)
		}
		r := rand.New(rand.NewSource(seed))
		chosenKey := keys[r.Intn(len(keys))]
		targetBundle = m.bundles[chosenKey]
	} else if scenarioID != "" {
		if b, ok := m.bundles[scenarioID]; ok {
			targetBundle = b
		} else {
			return nil, fmt.Errorf("scenario %q not found", scenarioID)
		}
	}

	if targetBundle == nil {
		return nil, fmt.Errorf("no scenario bundle available")
	}

	engine := NewSimulationEngine(targetBundle, seed, m.cfg)
	sessionID := engine.Session.ID

	sessionCtx, sessionCancel := context.WithCancel(m.ctx)
	m.engines[sessionID] = engine
	m.cancels[sessionID] = sessionCancel

	// Launch session ticker
	go m.runSessionLoop(sessionCtx, engine)

	return engine, nil
}


func (m *SessionManager) runSessionLoop(ctx context.Context, engine *SimulationEngine) {
	ticker := time.NewTicker(m.cfg.Simulation.TickDuration)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if !engine.Session.IsRunning() {
				return
			}
			engine.Tick(1)
		}
	}
}

// GetEngine returns an active session engine by ID.
func (m *SessionManager) GetEngine(sessionID string) (*SimulationEngine, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	engine, ok := m.engines[sessionID]
	if !ok {
		return nil, domain.ErrSessionNotFound
	}
	return engine, nil
}

// StopSession halts a running simulation session.
func (m *SessionManager) StopSession(sessionID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if cancel, ok := m.cancels[sessionID]; ok {
		cancel()
		delete(m.cancels, sessionID)
	}
	delete(m.engines, sessionID)
}

// Shutdown gracefully cancels all running simulation loops.
func (m *SessionManager) Shutdown() {
	m.cancel()
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, cancel := range m.cancels {
		cancel()
	}
	m.engines = make(map[string]*SimulationEngine)
	m.cancels = make(map[string]context.CancelFunc)
}
