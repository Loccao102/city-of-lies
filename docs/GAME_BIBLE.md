# City of Lies - Game Bible

## 1. Product Thesis
City of Lies is a real-time investigation game where the player reconstructs the truth while misinformation spreads autonomously through a society of 20 AI-driven characters with incomplete knowledge, individual beliefs, relationships, and motives.

## 2. Core Decisions & Rules
- **Winning Condition**: Submit all required Ground Truth fields correctly with Evidence Strength $\ge 0.70$ and Primary False Narrative Ratio $< 0.75$.
- **Losing Condition**: Primary False Narrative Ratio reaches $\ge 0.75$.
- **Belief Adoption Threshold**: An NPC adopts the false narrative when their narrative score reaches $\ge 0.65$.
- **Player Actions**:
  1. Interview a person (costs 30s game time).
  2. Inspect a location (costs 30–60s game time).
  3. Acquire evidence (costs 10–30s game time).
  4. Review notebook (0s game time).
  5. Publish Correction (costs 45s game time).
  6. Submit Ground Truth (costs 15s game time).

## 3. Propagation & Belief Mechanics
- **Propagation Strength**:
  $$S = 0.30 \times \text{speaker\_conf} + 0.25 \times \text{trust} + 0.15 \times \text{influence} + 0.15 \times \text{speaker\_cred} + 0.15 \times \text{receptiveness}$$
- **Skepticism Multiplier**:
  $$M_{\text{skep}} = 1 - (\text{listener\_skepticism} \times 0.55)$$
- **Contradictory Observation Modifier**:
  $$M_{\text{direct}} = 0.35$$
- **Repeated Root Source Multiplier**:
  $$M_{\text{provenance}} = 0.45$$
- **Update Delta**:
  $$\Delta = (1 - \text{old\_conf}) \times S \times 0.38 \times M_{\text{skep}} \times M_{\text{direct}} \times M_{\text{provenance}}$$
  $$\text{new\_conf} = \min(1.0, \max(0.0, \text{old\_conf} + \Delta))$$

## 4. Evidence Strength
- **Formula**:
  $$\text{Evidence Strength} = 0.50 \times \text{Coverage} + 0.30 \times \text{Reliability} + 0.20 \times \text{Diversity}$$
- Covers 4 mandatory fields: `event_type`, `cause_category`, `major_explosion`, `fatalities`.
