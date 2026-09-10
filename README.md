# City of Lies

> **City of Lies is a browser-based real-time multi-agent simulation in which autonomous characters form and propagate beliefs from partial observations, social trust and memory. Players investigate the same event and must establish the evidence-backed Ground Truth before a false narrative reaches social dominance.**

---

## 🔍 Core Concept

At approximately 17:58, an electrical fault in distribution panel DB-4 triggers a localized fire in the storage facility of Riverside Factory. There are no fatalities and no hazardous chemical release. However, as 20 distinct autonomous characters converse, misinterpret fragments, and share sensational rumors across the district, a dangerous false narrative emerges:

> *"A chemical explosion happened at Riverside Factory. Several workers died, and management is hiding the real casualty count."*

As the player, an independent investigator, you must interview witnesses, inspect locations, gather verified documentary and physical evidence, publish targeted corrections to curb misinformation, and submit the atomic Ground Truth before the primary false narrative reaches the **75% social dominance threshold**.

---

## 🛠️ Architecture & Technology Stack

- **Backend**: Go (Go 1.27) modular monolith, `net/http` + `go-chi/chi/v5`, `github.com/coder/websocket` for real-time WebSocket transport.
- **Database**: PostgreSQL 18 with `pgx/v5` connection pool and type-safe `sqlc` query generation.
- **Simulation**: Authoritative deterministic simulation engine with seeded RNG, provenance-tracked belief propagation, and weighted narrative metrics.
- **Frontend**: Next.js 16 (React 19, TypeScript strict), Three.js / React Three Fiber 9.x for 3D isometric district visualization, Zustand state management, Tailwind CSS.
- **LLM Safety Layer**: Zero LLM hallucination in game state. LLM controls only dialogue expression, bound by a strict Knowledge Guard. The simulation is 100% deterministic and fully playable even with LLM disabled.

---

## 🚀 Quick Start

### Prerequisites
- Go 1.27+
- Node.js 20+ / npm 10+
- Docker & Docker Compose (for PostgreSQL)

### 1. Run with Docker Compose
```bash
docker compose up --build
```
Open `http://localhost:3000` to start investigating.

### 2. Local Development
```bash
# Start PostgreSQL
docker compose -f docker-compose.dev.yml up -d

# Run Backend
cd backend
go run ./cmd/api

# Run Frontend
cd frontend
npm install
npm run dev
```

### 3. Headless Simulation Engine
Run the simulation engine without web UI to simulate rumor propagation and test scenarios:
```bash
cd backend
go run ./cmd/sim
```

### 4. Scenario Validation
Verify scenario integrity and referential constraints:
```bash
cd backend
go run ./cmd/scenario-check ../data/scenarios/riverside-factory
```

---

## 📖 Documentation

- [Master Game Bible](docs/GAME_BIBLE.md)
- [Riverside Scenario Specification](docs/SCENARIO_RIVERSIDE.md)
- [System Architecture](docs/ARCHITECTURE.md)
- [Protocol Specification](docs/PROTOCOL.md)
- [Database Schema & Migrations](docs/DATABASE.md)
- [Test Plan](docs/TEST_PLAN.md)
- [Implementation Status](docs/IMPLEMENTATION_STATUS.md)

---

## 📜 License
MIT License. See [LICENSE](LICENSE) for details.
