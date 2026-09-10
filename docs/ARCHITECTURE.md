# City of Lies - Architecture Specification

## Overview
The application is structured as a Go modular monolith backend communicating with a Next.js / React Three Fiber frontend via REST HTTP endpoints and a full-duplex WebSocket connection.

```text
[Browser / Next.js Client]
  ├── HTTP /api/v1 (REST actions: interview, inspect, correct, submit)
  └── WS /ws (Real-time events: pulses, movement, belief updates)
          │
          ▼
[Go Backend - Modular Monolith]
  ├── HTTP Handlers & Router (Chi)
  ├── WebSocket Hub (Coder WebSocket)
  ├── Application Services (Orchestration)
  ├── Simulation Engine (Deterministic tick runner)
  ├── Game Rules (Belief, Narrative, Outcome, Evidence)
  ├── Knowledge Guard & Dialogue Queue (OpenAI-compatible / Template)
  └── PostgreSQL Infrastructure (pgx + sqlc)
```

## Security & Privacy Boundary
- Ground Truth classification (`true`/`false`) is strictly server-only.
- NPCs are never omniscient; their dialogue context only receives their subjective memories and observed claims.
- The Knowledge Guard rejects any LLM-generated message attempting to leak server-hidden claims or undiscovered evidence.
- An LLM never alters game outcome, beliefs, or database state directly.
