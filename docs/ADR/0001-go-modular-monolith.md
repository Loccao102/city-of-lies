# Architecture Decision Records

## ADR 0001: Go Modular Monolith Architecture

### Status
Accepted

### Context
City of Lies requires real-time simulation with deterministic multi-agent state updates, concurrent player sessions, and WebSocket event distribution. Microservices or distributed message brokers (Kafka, RabbitMQ, Redis) would introduce unnecessary operational complexity, network latency, and serialization overhead for V1.

### Decision
Build the backend as a clean Go modular monolith structured with clear domain boundaries:
- `domain`: Pure business entities and invariants (Agent, Belief, Claim, Evidence, Memory).
- `game`: Simulation, propagation formulas, narrative calculations, and deterministic scheduler.
- `application`: Use-case services and ports.
- `infrastructure`: PostgreSQL repositories, WebSocket hub, and dialogue adapters.
- `http`: Routing, middleware, and request handlers.

### Consequences
- Single binary deployment with high throughput and low memory footprint.
- Zero distributed transaction complexity; PostgreSQL transactions ensure complete state consistency.
- Simple testing and local developer workflows.
