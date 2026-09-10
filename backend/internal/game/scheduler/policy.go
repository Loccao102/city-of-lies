package scheduler

type SalienceInput struct {
	Confidence      float64
	Novelty         float64
	EmotionalWeight float64
	Sociality       float64
	TrendingScore   float64
}

// CalculateSalience determines how eagerly an agent wants to share a specific claim.
func CalculateSalience(in SalienceInput) float64 {
	return (in.Confidence * 0.35) +
		(in.Novelty * 0.20) +
		(in.EmotionalWeight * 0.15) +
		(in.Sociality * 0.15) +
		(in.TrendingScore * 0.15)
}

type CooldownTracker struct {
	pairCooldowns map[string]int64 // "agentA:agentB" -> expiresAtSecond
}

func NewCooldownTracker() *CooldownTracker {
	return &CooldownTracker{
		pairCooldowns: make(map[string]int64),
	}
}

func pairKey(a, b string) string {
	if a < b {
		return a + ":" + b
	}
	return b + ":" + a
}

func (c *CooldownTracker) IsPairOnCooldown(a, b string, currentSecond int64) bool {
	key := pairKey(a, b)
	if expiresAt, exists := c.pairCooldowns[key]; exists {
		return currentSecond < expiresAt
	}
	return false
}

func (c *CooldownTracker) SetPairCooldown(a, b string, currentSecond int64, cooldownDuration int64) {
	key := pairKey(a, b)
	c.pairCooldowns[key] = currentSecond + cooldownDuration
}
