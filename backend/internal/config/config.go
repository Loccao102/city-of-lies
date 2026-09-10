package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type HTTPConfig struct {
	Addr string
}

type DatabaseConfig struct {
	URL string
}

type GameConfig struct {
	DefeatFalseBeliefRatio    float64
	BeliefAdoptionThreshold   float64
	TruthEvidenceThreshold     float64
	StartingPlayerCredibility float64
}

type SimulationConfig struct {
	TickDuration               time.Duration
	BeliefUpdateRate           float64
	DirectEvidenceMultiplier   float64
	SameRootSourceMultiplier   float64
	DefaultRelationshipTrust   float64
	MaxNewConversationsPerTick int
	PairCooldownSeconds        int64
	AgentShareCooldownSeconds  int64
}

type RealtimeConfig struct {
	WSSendQueue    int
	WSWriteTimeout time.Duration
	WSPingInterval time.Duration
}

type LLMConfig struct {
	Enabled       bool
	Provider      string
	BaseURL       string
	APIKey        string
	Model         string
	Workers       int
	QueueCapacity int
	Timeout       time.Duration
	RetryCount    int
}

type CORSConfig struct {
	AllowedOrigins []string
}

type Config struct {
	HTTP        HTTPConfig
	Database    DatabaseConfig
	ScenarioDir string
	Game        GameConfig
	Simulation  SimulationConfig
	Realtime    RealtimeConfig
	LLM         LLMConfig
	CORS        CORSConfig
}

// Load reads settings from environment variables with safe defaults.
func Load() (*Config, error) {
	cfg := &Config{
		HTTP: HTTPConfig{
			Addr: getEnv("CITYOFLIES_HTTP_ADDR", DefaultHTTPAddr),
		},
		Database: DatabaseConfig{
			URL: getEnv("CITYOFLIES_DATABASE_URL", ""),
		},
		ScenarioDir: getEnv("CITYOFLIES_SCENARIO_DIR", DefaultScenarioDir),
		Game: GameConfig{
			DefeatFalseBeliefRatio:    getEnvFloat("CITYOFLIES_DEFEAT_FALSE_BELIEF_RATIO", DefaultDefeatFalseBeliefRatio),
			BeliefAdoptionThreshold:   getEnvFloat("CITYOFLIES_BELIEF_ADOPTION_THRESHOLD", DefaultBeliefAdoptionThreshold),
			TruthEvidenceThreshold:     getEnvFloat("CITYOFLIES_TRUTH_EVIDENCE_THRESHOLD", DefaultTruthEvidenceThreshold),
			StartingPlayerCredibility: getEnvFloat("CITYOFLIES_STARTING_PLAYER_CREDIBILITY", DefaultStartingPlayerCredibility),
		},
		Simulation: SimulationConfig{
			TickDuration:               time.Duration(getEnvInt("CITYOFLIES_TICK_MS", 1000)) * time.Millisecond,
			BeliefUpdateRate:           getEnvFloat("CITYOFLIES_BELIEF_UPDATE_RATE", DefaultBeliefUpdateRate),
			DirectEvidenceMultiplier:   getEnvFloat("CITYOFLIES_DIRECT_EVIDENCE_MULTIPLIER", DefaultDirectEvidenceMultiplier),
			SameRootSourceMultiplier:   getEnvFloat("CITYOFLIES_SAME_ROOT_SOURCE_MULTIPLIER", DefaultSameRootSourceMultiplier),
			DefaultRelationshipTrust:   getEnvFloat("CITYOFLIES_DEFAULT_RELATIONSHIP_TRUST", DefaultRelationshipTrust),
			MaxNewConversationsPerTick: getEnvInt("CITYOFLIES_MAX_NEW_CONVERSATIONS_PER_TICK", DefaultMaxNewConversationsPerTick),
			PairCooldownSeconds:        int64(getEnvInt("CITYOFLIES_PAIR_COOLDOWN_SECONDS", DefaultPairCooldownSeconds)),
			AgentShareCooldownSeconds:  int64(getEnvInt("CITYOFLIES_AGENT_SHARE_COOLDOWN_SECONDS", DefaultAgentShareCooldownSeconds)),
		},
		Realtime: RealtimeConfig{
			WSSendQueue:    getEnvInt("CITYOFLIES_WS_SEND_QUEUE", DefaultWSSendQueue),
			WSWriteTimeout: time.Duration(getEnvInt("CITYOFLIES_WS_WRITE_TIMEOUT_MS", 5000)) * time.Millisecond,
			WSPingInterval: time.Duration(getEnvInt("CITYOFLIES_WS_PING_INTERVAL_MS", 20000)) * time.Millisecond,
		},
		LLM: LLMConfig{
			Enabled:       getEnvBool("CITYOFLIES_LLM_ENABLED", DefaultLLMEnabled),
			Provider:      getEnv("CITYOFLIES_LLM_PROVIDER", DefaultLLMProvider),
			BaseURL:       getEnv("CITYOFLIES_LLM_BASE_URL", ""),
			APIKey:        getEnv("CITYOFLIES_LLM_API_KEY", ""),
			Model:         getEnv("CITYOFLIES_LLM_MODEL", ""),
			Workers:       getEnvInt("CITYOFLIES_LLM_WORKERS", DefaultLLMWorkers),
			QueueCapacity: getEnvInt("CITYOFLIES_LLM_QUEUE_CAPACITY", DefaultLLMQueueCapacity),
			Timeout:       time.Duration(getEnvInt("CITYOFLIES_LLM_TIMEOUT_MS", 15000)) * time.Millisecond,
			RetryCount:    getEnvInt("CITYOFLIES_LLM_RETRY_COUNT", DefaultLLMRetryCount),
		},
		CORS: CORSConfig{
			AllowedOrigins: []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

// Validate checks sanity of all settings before application boot.
func (c *Config) Validate() error {
	if c.Game.DefeatFalseBeliefRatio <= 0.0 || c.Game.DefeatFalseBeliefRatio > 1.0 {
		return fmt.Errorf("invalid DefeatFalseBeliefRatio: %f", c.Game.DefeatFalseBeliefRatio)
	}
	if c.Game.BeliefAdoptionThreshold <= 0.0 || c.Game.BeliefAdoptionThreshold > 1.0 {
		return fmt.Errorf("invalid BeliefAdoptionThreshold: %f", c.Game.BeliefAdoptionThreshold)
	}
	if c.Game.TruthEvidenceThreshold <= 0.0 || c.Game.TruthEvidenceThreshold > 1.0 {
		return fmt.Errorf("invalid TruthEvidenceThreshold: %f", c.Game.TruthEvidenceThreshold)
	}
	if c.Game.StartingPlayerCredibility < 0.0 || c.Game.StartingPlayerCredibility > 1.0 {
		return fmt.Errorf("invalid StartingPlayerCredibility: %f", c.Game.StartingPlayerCredibility)
	}
	if c.Simulation.TickDuration <= 0 {
		return fmt.Errorf("invalid TickDuration: %v", c.Simulation.TickDuration)
	}
	if c.Realtime.WSSendQueue <= 0 {
		return fmt.Errorf("invalid WSSendQueue: %d", c.Realtime.WSSendQueue)
	}
	if c.LLM.Enabled && c.LLM.Provider != "template" && c.LLM.BaseURL == "" {
		return fmt.Errorf("LLM enabled with provider %s but BaseURL is empty", c.LLM.Provider)
	}
	return nil
}

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if valStr, ok := os.LookupEnv(key); ok && valStr != "" {
		if val, err := strconv.Atoi(valStr); err == nil {
			return val
		}
	}
	return defaultVal
}

func getEnvFloat(key string, defaultVal float64) float64 {
	if valStr, ok := os.LookupEnv(key); ok && valStr != "" {
		if val, err := strconv.ParseFloat(valStr, 64); err == nil {
			return val
		}
	}
	return defaultVal
}

func getEnvBool(key string, defaultVal bool) bool {
	if valStr, ok := os.LookupEnv(key); ok && valStr != "" {
		if val, err := strconv.ParseBool(valStr); err == nil {
			return val
		}
	}
	return defaultVal
}
