# ADR 0003: LLM Never Owns Truth or Game State

### Status
Accepted

### Context
LLMs are prone to hallucinations, unbounded latency, and non-deterministic behavior. If an LLM directly decided what is true, what characters believe, or whether the player won or lost, the game would become unbalanceable and exploitable.

### Decision
The simulation engine owns reality. The LLM is only an expressive dialogue surface:
1. Ground Truth is never included in LLM prompts.
2. The LLM receives only the character's subjective memories, known claims, and allowed evidence IDs.
3. A strict server-side Knowledge Guard rejects any LLM response referencing unknown claims or undiscovered evidence.
4. If an LLM call fails or times out, a deterministic template fallback responds immediately.
5. The game remains 100% playable with LLM integration disabled.

### Consequences
- Zero LLM hallucination leak into game state.
- Predictable game loops and sub-millisecond tick evaluations.
- Operational resilience against external API outages.
