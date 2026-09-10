package dto

type SessionResponse struct {
	ID                  string  `json:"id"`
	ScenarioID          string  `json:"scenario_id"`
	Status              string  `json:"status"`
	GameSecond          int64   `json:"game_second"`
	FalseNarrativeRatio float64 `json:"primary_false_narrative_ratio"`
	EvidenceStrength    float64 `json:"evidence_strength"`
	PlayerCredibility   float64 `json:"player_credibility"`
	SimulationSeed      int64   `json:"simulation_seed"`
	Revision            int64   `json:"revision"`
	Sequence            int64   `json:"sequence"`
}
