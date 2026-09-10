# Wire Protocol Specification

## REST API (`/api/v1`)

### Sessions
- `POST /api/v1/sessions`: Create new session. Optional `{"scenario_id": "riverside-factory", "seed": 12345}`. Returns session snapshot.
- `GET /api/v1/sessions/{id}`: Get player-safe session snapshot.
- `GET /api/v1/sessions/{id}/agents`: List all 20 agents with public status.
- `GET /api/v1/sessions/{id}/agents/{agentId}`: Get detailed public profile and interview history.
- `POST /api/v1/sessions/{id}/agents/{agentId}/interviews`: Submit a question. Returns `200` with answer or `202` pending async worker response.
- `POST /api/v1/sessions/{id}/locations/{locationId}/inspect`: Inspect location to discover evidence items.
- `GET /api/v1/sessions/{id}/notebook`: Get player's current notebook (discovered evidence, known claims, timeline).
- `POST /api/v1/sessions/{id}/corrections`: Publish a fact-checking correction. Body: `{"challenged_claim_id": "...", "evidence_ids": [...], "message": "..."}`.
- `POST /api/v1/sessions/{id}/truth-submissions`: Submit final Ground Truth report. Body: `{"event_type": "fire", "cause_category": "electrical_fault", "major_explosion": false, "fatalities": 0}`.
- `GET /health/live`, `GET /health/ready`: Service health endpoints.

## WebSocket Protocol (`/ws?session_id={id}`)

### Messages to Client
- `sequence`: Monotonically increasing integer.
- `type`: Event type string.
- `game_second`: Current simulation time.
- `payload`: Event-specific JSON payload.

### Event Types:
- `session.snapshot`: Full player-safe state on connect.
- `session.status_changed`: State transition (`running`, `won`, `lost_false_belief`).
- `agent.moved`: Agent path between locations.
- `conversation.agent_message`: Agent dialogue utterance in response to player interview.
- `rumor.visualized`: NPC-to-NPC conversation pulse for 3D map.
- `belief_stats.updated`: Updated false narrative ratio.
- `evidence.discovered`: Notification of new evidence added to notebook.
- `correction.published`: Public correction result and credibility update.
- `game.won`, `game.lost`: Game outcome terminal event.
