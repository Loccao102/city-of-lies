package dto

type AgentSummaryDTO struct {
	ID              string `json:"id"`
	ScenarioAgentID string `json:"scenario_agent_id"`
	Name            string `json:"name"`
	Role            string `json:"role"`
	LocationID      string `json:"location_id"`
	PublicStatus    string `json:"public_status"`
	Interviewed     bool   `json:"interviewed"`
}

type AgentDetailDTO struct {
	AgentSummaryDTO
	Personality string `json:"personality"`
}
