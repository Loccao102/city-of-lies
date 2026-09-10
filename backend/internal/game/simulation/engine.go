package simulation

import (
	"encoding/json"
	"sync"
	"time"

	"city-of-lies/backend/internal/config"
	"city-of-lies/backend/internal/domain"
	"city-of-lies/backend/internal/game/belief"
	"city-of-lies/backend/internal/game/clock"
	"city-of-lies/backend/internal/game/movement"
	"city-of-lies/backend/internal/game/outcome"
	"city-of-lies/backend/internal/game/scenario"
	"city-of-lies/backend/internal/game/scheduler"
)

type SimulationEngine struct {
	mu                  sync.RWMutex
	Session             domain.GameSession
	Agents              map[string]*domain.Agent       // scenario_agent_id -> Agent
	Claims              map[string]domain.Claim        // scenario_claim_id -> Claim
	Beliefs             map[string]map[string]*domain.AgentBelief // agentID -> claimID -> Belief
	Memories            map[string][]domain.AgentMemory // agentID -> memories
	Relationships       map[string]map[string]float64  // from -> to -> trust
	Evidence            map[string]*domain.EvidenceItem // scenario_evidence_id -> Evidence
	Waypoints           *movement.WaypointGraph
	Scheduler           *scheduler.Scheduler
	RNG                 *scheduler.SeededRNG
	Clock               clock.Clock
	Config              *config.Config
	LatestSequence      int64
	EventLog            []domain.WorldEvent
	OnEventEmitted      func(event domain.WorldEvent)
	ScenarioBundle      *scenario.ScenarioBundle
}

// NewSimulationEngine initializes an in-memory simulation engine from a scenario bundle and seed.
func NewSimulationEngine(bundle *scenario.ScenarioBundle, seed int64, cfg *config.Config) *SimulationEngine {
	rng := scheduler.NewSeededRNG(seed)
	clk := clock.NewRealClock(0)
	sched := scheduler.NewScheduler(rng, cfg.Simulation.PairCooldownSeconds, cfg.Simulation.MaxNewConversationsPerTick)

	sessionID := domain.NewUUID()
	eng := &SimulationEngine{
		Session: domain.GameSession{
			ID:                  sessionID,
			ScenarioID:          bundle.Metadata.ID,
			Status:              domain.SessionStatusRunning,
			SimulationSeed:      seed,
			GameSecond:          0,
			PlayerCredibility:   cfg.Game.StartingPlayerCredibility,
			FalseNarrativeRatio: 0.0,
			EvidenceStrength:    0.0,
			Revision:            1,
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
		},
		Agents:         make(map[string]*domain.Agent),
		Claims:         make(map[string]domain.Claim),
		Beliefs:        make(map[string]map[string]*domain.AgentBelief),
		Memories:       make(map[string][]domain.AgentMemory),
		Relationships:  make(map[string]map[string]float64),
		Evidence:       make(map[string]*domain.EvidenceItem),
		Waypoints:      movement.NewWaypointGraph(bundle.Locations),
		Scheduler:      sched,
		RNG:            rng,
		Clock:          clk,
		Config:         cfg,
		ScenarioBundle: bundle,
	}

	// 1. Seed Agents
	for _, aSeed := range bundle.Agents {
		agentID := domain.NewUUID()
		agent := &domain.Agent{
			ID:                agentID,
			SessionID:         sessionID,
			ScenarioAgentID:   aSeed.ID,
			Name:              aSeed.Name,
			Role:              aSeed.Role,
			CurrentLocationID: aSeed.LocationID,
			Influence:         aSeed.Influence,
			Credibility:       aSeed.Credibility,
			Skepticism:        aSeed.Skepticism,
			Sociality:         aSeed.Sociality,
			DeceptionTendency: aSeed.DeceptionTendency,
			Receptiveness:     aSeed.Receptiveness,
			Revision:          1,
		}
		eng.Agents[aSeed.ID] = agent
		eng.Beliefs[agent.ID] = make(map[string]*domain.AgentBelief)
		eng.Memories[agent.ID] = make([]domain.AgentMemory, 0)
	}

	// 2. Seed Claims
	for _, c := range bundle.Claims {
		c.SessionID = sessionID
		c.ID = domain.NewUUID()
		eng.Claims[c.ScenarioClaimID] = c
	}

	// 3. Seed Observations / Memories / Initial Beliefs
	for _, obs := range bundle.Observations {
		agent, okA := eng.Agents[obs.AgentID]
		claim, okC := eng.Claims[obs.ClaimID]
		if okA && okC {
			eng.Beliefs[agent.ID][claim.ScenarioClaimID] = &domain.AgentBelief{
				SessionID:           sessionID,
				AgentID:             agent.ID,
				ClaimID:             claim.ScenarioClaimID,
				Confidence:          obs.Confidence,
				UpdatedAtGameSecond: 0,
				Revision:            1,
			}
			eng.Memories[agent.ID] = append(eng.Memories[agent.ID], domain.AgentMemory{
				ID:                  domain.NewUUID(),
				SessionID:           sessionID,
				AgentID:             agent.ID,
				MemoryType:          domain.MemoryTypeDirectObservation,
				Content:             claim.DisplayText,
				Reliability:         0.90,
				Importance:          0.80,
				RelatedClaimIDs:     []string{claim.ScenarioClaimID},
				CreatedAtGameSecond: 0,
			})
		}
	}

	// 4. Seed Relationships
	for _, rel := range bundle.Relationships {
		fromAgent, okF := eng.Agents[rel.FromAgentID]
		toAgent, okT := eng.Agents[rel.ToAgentID]
		if okF && okT {
			if _, exists := eng.Relationships[fromAgent.ID]; !exists {
				eng.Relationships[fromAgent.ID] = make(map[string]float64)
			}
			eng.Relationships[fromAgent.ID][toAgent.ID] = rel.Trust
		}
	}

	// 5. Seed Evidence
	for _, evSeed := range bundle.Evidence {
		eng.Evidence[evSeed.ID] = &domain.EvidenceItem{
			ID:                 domain.NewUUID(),
			SessionID:          sessionID,
			ScenarioEvidenceID: evSeed.ID,
			Name:               evSeed.Name,
			Description:        evSeed.Description,
			LocationID:         evSeed.LocationID,
			Reliability:        evSeed.Reliability,
			SourceCategory:     evSeed.SourceCategory,
			Discovered:         false,
		}
	}

	// Calculate initial ratio
	eng.recalculateNarrativeRatio()

	return eng
}

func (e *SimulationEngine) recalculateNarrativeRatio() {
	var weights []belief.NarrativeClaimWeight
	for _, c := range e.ScenarioBundle.Metadata.PrimaryFalseNarrative.Claims {
		weights = append(weights, belief.NarrativeClaimWeight{
			ClaimID: c.ClaimID,
			Weight:  c.Weight,
		})
	}

	adoptedCount := 0
	totalAgents := len(e.Agents)

	for _, agent := range e.Agents {
		agentBeliefs := make(map[string]float64)
		for claimID, b := range e.Beliefs[agent.ID] {
			agentBeliefs[claimID] = b.Confidence
		}
		score := belief.CalculateAgentNarrativeScore(weights, agentBeliefs)
		if belief.IsNarrativeAdopted(score, e.Config.Game.BeliefAdoptionThreshold) {
			adoptedCount++
		}
	}

	e.Session.FalseNarrativeRatio = belief.CalculateFalseNarrativeRatio(adoptedCount, totalAgents)

	// Check immediate loss
	if outcome.CheckLossCondition(e.Session.FalseNarrativeRatio, e.Config.Game.DefeatFalseBeliefRatio) {
		e.Session.Status = domain.SessionStatusLostFalseBelief
		now := time.Now()
		e.Session.EndedAt = &now
	}
}

// AppendEvent stores and broadcasts an immutable simulation event.
func (e *SimulationEngine) AppendEvent(eventType domain.WorldEventType, actorID, targetID *string, payload interface{}) domain.WorldEvent {
	e.LatestSequence++
	payloadBytes, _ := json.Marshal(payload)
	evt := domain.WorldEvent{
		ID:            domain.NewUUID(),
		SessionID:     e.Session.ID,
		Sequence:      e.LatestSequence,
		EventType:     eventType,
		GameSecond:    e.Session.GameSecond,
		ActorAgentID:  actorID,
		TargetAgentID: targetID,
		Payload:       payloadBytes,
		OccurredAt:    time.Now(),
	}
	e.EventLog = append(e.EventLog, evt)
	if e.OnEventEmitted != nil {
		e.OnEventEmitted(evt)
	}
	return evt
}
