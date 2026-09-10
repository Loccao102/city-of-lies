package config

import "time"

const (
	DefaultHTTPAddr = ":8080"
	DefaultScenarioDir = "./data/scenarios"

	// Game Defaults
	DefaultDefeatFalseBeliefRatio    = 0.75
	DefaultBeliefAdoptionThreshold   = 0.65
	DefaultTruthEvidenceThreshold     = 0.70
	DefaultStartingPlayerCredibility = 0.50

	// Simulation Defaults
	DefaultTickDuration               = 1000 * time.Millisecond
	DefaultBeliefUpdateRate           = 0.38
	DefaultDirectEvidenceMultiplier   = 0.35
	DefaultSameRootSourceMultiplier   = 0.45
	DefaultRelationshipTrust         = 0.50
	DefaultMaxNewConversationsPerTick = 2
	DefaultPairCooldownSeconds        = 90
	DefaultAgentShareCooldownSeconds  = 40

	// Realtime Defaults
	DefaultWSSendQueue        = 256
	DefaultWSWriteTimeout     = 5 * time.Second
	DefaultWSPingInterval     = 20 * time.Second

	// LLM Defaults
	DefaultLLMEnabled       = false
	DefaultLLMProvider      = "template"
	DefaultLLMWorkers       = 3
	DefaultLLMQueueCapacity = 50
	DefaultLLMTimeout       = 15 * time.Second
	DefaultLLMRetryCount    = 1
	DefaultLLMMaxOutputChars = 1200
)
