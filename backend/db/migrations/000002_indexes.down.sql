-- 000002_indexes.down.sql
drop index if exists idx_evidence_session_discovered;
drop index if exists idx_beliefs_session_claim_confidence;
drop index if exists idx_session_agents_session_location;
drop index if exists idx_agent_memories_agent_time;
drop index if exists idx_world_events_session_sequence;
