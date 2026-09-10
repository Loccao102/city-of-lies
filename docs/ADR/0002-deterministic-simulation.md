# ADR 0002: Deterministic Simulation Engine

### Status
Accepted

### Context
Scientific testing, balance tuning, and regression test suites require identical outcomes given identical scenario inputs and player actions. Relying on unseeded random sources or real-time clocks inside simulation math leads to flakiness.

### Decision
All simulation state transitions, agent conversation scheduling, and rumor propagation probabilities are derived deterministically using a per-session seed and integer `game_second` ticks. No global `rand` is permitted in the simulation domain.

### Consequences
- Any gameplay session can be perfectly replayed by supplying the initial seed.
- Headless simulation benchmark suites run at thousands of ticks per second.
- Reproducible unit and soak testing.
