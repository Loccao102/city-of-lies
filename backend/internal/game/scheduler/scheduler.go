package scheduler



type AgentSchedulingView struct {
	ID                string
	ScenarioAgentID   string
	LocationID        string
	Influence         float64
	Sociality         float64
	BusyUntilSecond   int64
	Beliefs           map[string]float64 // claim_id -> confidence
}

type ScheduledInteraction struct {
	SpeakerAgentID string
	ListenerAgentID string
	ClaimID        string
	Salience       float64
}

type Scheduler struct {
	rng             *SeededRNG
	cooldowns       *CooldownTracker
	pairCooldownSec int64
	maxPerTick      int
	claimRoles      map[string]string
	claimWeights    map[string]float64
}

func NewScheduler(rng *SeededRNG, pairCooldownSec int64, maxPerTick int) *Scheduler {
	if pairCooldownSec <= 0 {
		pairCooldownSec = 90
	}
	if maxPerTick <= 0 {
		maxPerTick = 2
	}
	return &Scheduler{
		rng:             rng,
		cooldowns:       NewCooldownTracker(),
		pairCooldownSec: pairCooldownSec,
		maxPerTick:      maxPerTick,
		claimRoles:      make(map[string]string),
		claimWeights:    make(map[string]float64),
	}
}

func (s *Scheduler) SetClaimMetadata(roles map[string]string, weights map[string]float64) {
	s.claimRoles = roles
	s.claimWeights = weights
}

// PlanTickInteractions identifies valid communication pairs and salient claims for the current tick.
func (s *Scheduler) PlanTickInteractions(
	agents []AgentSchedulingView,
	relationships map[string]map[string]float64, // from -> to -> trust
	currentSecond int64,
) []ScheduledInteraction {
	if len(agents) < 2 {
		return nil
	}

	var scheduled []ScheduledInteraction
	agentCount := len(agents)
	perm := s.rng.Perm(agentCount)

	for _, idxA := range perm {
		if len(scheduled) >= s.maxPerTick {
			break
		}

		speaker := agents[idxA]
		if speaker.BusyUntilSecond > currentSecond {
			continue
		}

		// Find speaker's most salient claim
		bestClaimID := ""
		bestSalience := -1.0

		for claimID, conf := range speaker.Beliefs {
			if conf < 0.20 {
				continue // Don't actively spread claims with very low conviction
			}
			novelty, emotionalWeight, trending := s.getClaimEmotionalWeights(claimID)

			salience := CalculateSalience(SalienceInput{
				Confidence:      conf,
				Novelty:         novelty,
				EmotionalWeight: emotionalWeight,
				Sociality:       speaker.Sociality,
				TrendingScore:   trending,
			})
			if salience > bestSalience {
				bestSalience = salience
				bestClaimID = claimID
			}
		}

		if bestClaimID == "" {
			continue
		}

		// Look for a suitable listener
		for _, idxB := range s.rng.Perm(agentCount) {
			if idxA == idxB {
				continue
			}

			listener := agents[idxB]
			if listener.BusyUntilSecond > currentSecond {
				continue
			}

			// Check if on cooldown
			if s.cooldowns.IsPairOnCooldown(speaker.ID, listener.ID, currentSecond) {
				continue
			}

			// Check social affinity or colocation
			isColocated := speaker.LocationID == listener.LocationID
			trust := 0.50
			if rels, exists := relationships[speaker.ID]; exists {
				if t, ok := rels[listener.ID]; ok {
					trust = t
				}
			}

			// Must be in same location OR have established trust relationship OR speaker has broadcast influence (e.g. Blogger, Influencer)
			canReach := isColocated || trust >= 0.50 || (speaker.Influence >= 0.60 && s.rng.Float64() <= speaker.Influence)
			if !canReach {
				continue
			}

			// Probabilistic roll modulated by sociality or broadcast reach
			prob := (speaker.Sociality + trust) / 2.0
			if speaker.Influence >= 0.60 && prob < speaker.Influence*0.75 {
				prob = speaker.Influence * 0.75
			}
			if s.rng.Float64() <= prob {
				scheduled = append(scheduled, ScheduledInteraction{
					SpeakerAgentID:  speaker.ID,
					ListenerAgentID: listener.ID,
					ClaimID:         bestClaimID,
					Salience:        bestSalience,
				})
				s.cooldowns.SetPairCooldown(speaker.ID, listener.ID, currentSecond, s.pairCooldownSec)
				break
			}
		}
	}

	return scheduled
}

func (s *Scheduler) getClaimEmotionalWeights(claimID string) (novelty, emotionalWeight, trending float64) {
	// 1. Dynamic lookup by claim role & narrative weight if metadata exists
	if s.claimRoles != nil {
		role, hasRole := s.claimRoles[claimID]
		weight := s.claimWeights[claimID]
		if hasRole {
			switch role {
			case "primary_false":
				if weight >= 0.38 { // Highest shock value: casualties / deaths / bank run / collapse
					return 0.95, 0.98, 0.95
				} else if weight <= 0.28 { // Intermediate escalation: coverup / conspiracy / corruption
					return 0.85, 0.88, 0.85
				} else { // Initial sensational rumor: explosion / virus / poison
					return 0.70, 0.80, 0.75
				}
			case "secondary_false":
				return 0.70, 0.75, 0.70
			case "fact":
				return 0.30, 0.35, 0.30
			}
		}
	}

	// 2. Legacy fallback for riverside hardcoded claim IDs
	switch claimID {
	case "claim_multiple_deaths":
		return 0.95, 0.98, 0.95
	case "claim_company_coverup":
		return 0.85, 0.88, 0.85
	case "claim_chemical_explosion":
		return 0.70, 0.80, 0.75
	case "claim_toxic_cloud":
		return 0.85, 0.90, 0.85
	default:
		return 0.30, 0.35, 0.30
	}
}
