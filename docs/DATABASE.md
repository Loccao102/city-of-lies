# Database Schema Specification

PostgreSQL 18 is the authoritative persistent store for City of Lies.

## Schema Tables
1. `game_sessions`: Core session aggregate and metrics.
2. `session_agents`: 20 persistent agent actors per session.
3. `session_claims`: Cloned scenario claims (classification kept server-safe).
4. `agent_beliefs`: Directed belief pairs `(agent_id, claim_id)` with confidence and provenance tracking.
5. `agent_memories`: Ephemeral and direct memories with root source tags.
6. `relationships`: Directed trust graph `(from_agent, to_agent, trust)`.
7. `evidence_items`: Discovered and undiscovered clues.
8. `evidence_claim_links`: Relational mappings between evidence and supported/contradicted claims.
9. `conversations`: Header for player-agent interview sessions.
10. `conversation_messages`: Chronological dialogue log.
11. `publications`: Corrections submitted by the player with calculated strength.
12. `publication_evidence`: Relational links between corrections and attached evidence.
13. `world_events`: Monotonically ordered, immutable sequence of simulation and player events.
