-- 000002_indexes.up.sql
create index if not exists idx_world_events_session_sequence on world_events(session_id, sequence);
create index if not exists idx_agent_memories_agent_time on agent_memories(agent_id, created_at_game_second desc);
create index if not exists idx_session_agents_session_location on session_agents(session_id, current_location_id);
create index if not exists idx_beliefs_session_claim_confidence on agent_beliefs(session_id, claim_id, confidence);
create index if not exists idx_evidence_session_discovered on evidence_items(session_id, discovered);
