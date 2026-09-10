export type SessionStatus =
  | "created"
  | "running"
  | "won"
  | "lost_false_belief"
  | "lost_timeout"
  | "aborted";

export interface SessionResponse {
  id: string;
  scenario_id: string;
  status: SessionStatus;
  game_second: number;
  primary_false_narrative_ratio: number;
  evidence_strength: number;
  player_credibility: number;
  simulation_seed: number;
  revision: number;
  sequence: number;
}

export interface AgentSummary {
  id: string;
  scenario_agent_id: string;
  name: string;
  role: string;
  location_id: string;
  public_status: string;
  interviewed: boolean;
}

export interface AgentDetail extends AgentSummary {
  personality: string;
}

export interface EvidenceItem {
  id: string;
  scenario_evidence_id: string;
  name: string;
  description: string;
  location_id: string;
  reliability: number;
  source_category: string;
  supports?: string[];
  contradicts?: string[];
  discovered: boolean;
}

export interface ClaimSummary {
  id: string;
  display_text: string;
  status: "Unverified" | "Supported" | "Contradicted";
}

export interface TimelineEvent {
  game_second: number;
  headline: string;
  description: string;
}

export interface NotebookData {
  discovered_evidence: EvidenceItem[];
  known_claims: ClaimSummary[];
  interviewed_agents: string[];
  timeline_events: TimelineEvent[];
}

export interface DialogueResponse {
  intent: string;
  utterance: string;
  referenced_claim_ids: string[];
  revealed_evidence_ids: string[];
  emotion: string;
  certainty: number;
}

export interface InspectResponse {
  location_id: string;
  discovered_evidence: EvidenceItem[];
  message: string;
}

export interface Publication {
  id: string;
  session_id: string;
  challenged_claim_id: string;
  message: string;
  strength: number;
  player_credibility_before: number;
  player_credibility_after: number;
  published_at_game_second: number;
  evidence_ids: string[];
}

export interface SubmitTruthRequest {
  event_type: string;
  cause_category: string;
  major_explosion: boolean;
  fatalities: number;
  chemical_release?: boolean;
  management_issue?: string;
}

export interface TruthSubmissionResult {
  won: boolean;
  terminal: boolean;
  evidence_strength: number;
  required_evidence: number;
  false_belief_ratio: number;
  field_status: Record<string, boolean>;
  feedback_message: string;
}

export interface WorldEvent {
  id: string;
  session_id: string;
  sequence: number;
  event_type: string;
  game_second: number;
  actor_agent_id?: string;
  target_agent_id?: string;
  payload: Record<string, any>;
  occurred_at: string;
}

export interface LocationInfo {
  id: string;
  name: string;
  description: string;
  x: number;
  z: number;
  category: "factory" | "public" | "official" | "residential";
}
