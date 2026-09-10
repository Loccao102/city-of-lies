# ADR 0004: WebSocket Over Polling for Real-Time State

### Status
Accepted

### Context
Characters converse, move, and spread rumors every second. HTTP polling would introduce unacceptable latency, redundant query load, and race conditions between ticks and player views.

### Decision
Use full-duplex WebSockets powered by `github.com/coder/websocket` for streaming game events:
- Monotonically increasing `sequence` IDs per session.
- Bounded per-client send channels (256 messages) with drop-and-disconnect policy for slow consumers.
- Automatic snapshot resynchronization on reconnection if a sequence gap is detected.

### Consequences
- Fluid 3D rumor pulses and real-time belief meters.
- Low server CPU and database polling overhead.
- Clean client reconciliation on network reconnects.
