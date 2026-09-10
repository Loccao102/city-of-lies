# Test Plan & Verification Strategy

## 1. Domain Unit Tests (Pure Functions)
- `propagation_test.go`: Verifies propagation delta, clamps, relationship trust weights, skepticism reduction, direct observation resistance, and repeated provenance attenuation.
- `narrative_test.go`: Verifies weighted narrative calculation and adoption threshold ($0.65$) logic.
- `strength_test.go`: Verifies evidence strength score (coverage, reliability, diversity).
- `evaluator_test.go`: Verifies game outcome state machine (win conditions, loss threshold $0.75$, terminal state immutability).
- `validator_test.go`: Referential integrity check of scenario JSON.

## 2. Integration & Concurrency Tests
- Session creation and seeding transaction.
- Headless simulation soak test (running deterministic ticks to loss).
- Race detector validation (`go test -race ./...`).

## 3. Frontend Tests & Typecheck
- TypeScript strict compiler check (`tsc --noEmit`).
- Real-time sequence handling and deduplication in Zustand stores.
