package simulation

import (
	"city-of-lies/backend/internal/domain"
	"city-of-lies/backend/internal/game/belief"
	"city-of-lies/backend/internal/game/scheduler"
)

type RumorVisualizedPayload struct {
	SpeakerAgentID   string `json:"speaker_agent_id"`
	ListenerAgentID  string `json:"listener_agent_id"`
	SpeakerName      string `json:"speaker_name"`
	ListenerName     string `json:"listener_name"`
	LocationID       string `json:"location_id"`
	Category         string `json:"category"`
}

type BeliefStatsPayload struct {
	PrimaryFalseNarrativeRatio float64 `json:"primary_false_narrative_ratio"`
	EvidenceStrength           float64 `json:"evidence_strength"`
	PlayerCredibility          float64 `json:"player_credibility"`
}

// Tick executes one atomic step of the simulation.
func (e *SimulationEngine) Tick(deltaSeconds int64) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if !e.Session.IsRunning() {
		return
	}

	e.Session.GameSecond += deltaSeconds
	currentSecond := e.Session.GameSecond
	e.Clock.AdvanceSeconds(deltaSeconds)

	// 1. Process Master Event Timeline anchors
	for _, pe := range e.ScenarioBundle.PublicEvents {
		if pe.GameSecond == currentSecond {
			e.AppendEvent(
				domain.WorldEventTypeRumorShared,
				nil, nil,
				map[string]interface{}{
					"headline":    pe.Headline,
					"description": pe.Description,
					"location_id": pe.LocationID,
				},
			)

			// 1a. Process data-driven BeliefSeeds from PublicEvent (works for any scenario)
			for _, bs := range pe.BeliefSeeds {
				e.seedAgentBelief(bs.AgentID, bs.ClaimID, bs.Confidence, bs.Reasoning)
			}

			// 1b. Legacy hardcoded fallback for riverside-factory when no belief seeds are specified in json
			if len(pe.BeliefSeeds) == 0 && (e.ScenarioBundle == nil || e.ScenarioBundle.Metadata.ID == "riverside-factory") {
				switch pe.GameSecond {
				case 420: // Lan -> Hùng
					e.seedAgentBelief("agent_factory_worker_b", "claim_chemical_explosion", 0.45, "Lan nói bên trong có nổ lớn sau khi tủ điện chập.")
				case 480: // Hùng -> Thảo
					e.seedAgentBelief("agent_shop_owner", "claim_chemical_explosion", 0.50, "Hùng kể trong nhà máy xảy ra nổ hóa chất và xe cấp cứu đã vào.")
				case 720: // Blogger Khoa posts initial rumor
					e.seedAgentBelief("agent_blogger", "claim_chemical_explosion", 0.60, "Nhận tin báo và ảnh khói nghi nổ tại nhà kho.")
				case 840: // Khoa updates headline: "Nghi vấn nổ hóa chất"
					e.seedAgentBelief("agent_blogger", "claim_chemical_explosion", 0.72, "Đăng bài nghi vấn nổ hóa chất tại Riverside.")
					e.seedAgentBelief("agent_influencer", "claim_chemical_explosion", 0.68, "Thấy bài đăng nghi nổ hóa chất của Khoa Bùi.")
				case 900: // Influencer Vy shares post
					e.seedAgentBelief("agent_influencer", "claim_chemical_explosion", 0.80, "Chia sẻ bài viết nổ hóa chất và cảnh báo cư dân.")
				case 1080: // Sơn gives evasive statement -> seeds coverup
					e.seedAgentBelief("agent_blogger", "claim_company_coverup", 0.65, "Ban quản lý nhà máy né tránh trả lời thương vong, có dấu hiệu giấu tin.")
					e.seedAgentBelief("agent_influencer", "claim_company_coverup", 0.70, "Quản lý ấp úng khi phóng viên hỏi, nghi vấn bưng bít thông tin.")
					e.seedAgentBelief("agent_resident", "claim_company_coverup", 0.55, "Cư dân bàn tán công ty đang tìm cách ém nhẹm vụ việc.")
					e.seedAgentBelief("agent_shop_owner", "claim_company_coverup", 0.50, "Xem tivi thấy quản lý trả lời quanh co.")
					e.seedAgentBelief("agent_vendor", "claim_company_coverup", 0.50, "Người dân xem tin tức bàn tán công ty giấu tin người chết.")
					e.seedAgentBelief("agent_taxi_driver", "claim_company_coverup", 0.50, "Nghe đài radio thấy công ty từ chối công bố số thương vong.")
					e.seedAgentBelief("agent_delivery_driver", "claim_company_coverup", 0.50, "Thấy xe cấp cứu mà công ty bảo chưa có gì nghiêm trọng.")
					e.seedAgentBelief("agent_student", "claim_company_coverup", 0.52, "Đọc tin thấy ban quản lý né tránh câu hỏi của báo chí.")
					e.seedAgentBelief("agent_factory_worker_a", "claim_company_coverup", 0.50, "Bản thân ở kho thấy cháy mà quản lý không nói rõ sự thật.")
				case 1200: // Vy livestreams claiming casualties
					e.seedAgentBelief("agent_influencer", "claim_multiple_deaths", 0.75, "Phát sóng livestream nói có nguồn tin khẳng định có người chết.")
					e.seedAgentBelief("agent_blogger", "claim_multiple_deaths", 0.68, "Nghe livestream của Vy Hoàng nói đã có người tử vong.")
					e.seedAgentBelief("agent_shop_owner", "claim_multiple_deaths", 0.55, "Khách vào quán bàn tán xôn xao về thông tin có công nhân chết.")
				}
			}
		}
	}

	// 1c. Additional timeline anchor at 1320 (18:22): Khoa headlines casualties (riverside fallback)
	if currentSecond == 1320 && (e.ScenarioBundle == nil || e.ScenarioBundle.Metadata.ID == "riverside-factory") {
		e.seedAgentBelief("agent_blogger", "claim_multiple_deaths", 0.78, "Cập nhật tiêu đề bài viết: nguồn tin nói có thương vong tử vong.")
	}

	// 2. Prepare Agents for Scheduling
	agentViews := make([]scheduler.AgentSchedulingView, 0, len(e.Agents))
	agentLookup := make(map[string]*domain.Agent)
	for _, a := range e.Agents {
		agentLookup[a.ID] = a
		beliefsMap := make(map[string]float64)
		for cID, b := range e.Beliefs[a.ID] {
			beliefsMap[cID] = b.Confidence
		}
		agentViews = append(agentViews, scheduler.AgentSchedulingView{
			ID:              a.ID,
			ScenarioAgentID: a.ScenarioAgentID,
			LocationID:      a.CurrentLocationID,
			Influence:       a.Influence,
			Sociality:       a.Sociality,
			BusyUntilSecond: a.BusyUntilGameSecond,
			Beliefs:         beliefsMap,
		})
	}

	// 3. Schedule autonomous NPC conversations
	interactions := e.Scheduler.PlanTickInteractions(agentViews, e.Relationships, currentSecond)

	for _, inter := range interactions {
		speaker := agentLookup[inter.SpeakerAgentID]
		listener := agentLookup[inter.ListenerAgentID]
		if speaker == nil || listener == nil {
			continue
		}

		speakerBelief, ok := e.Beliefs[speaker.ID][inter.ClaimID]
		if !ok || speakerBelief.Confidence < 0.20 {
			continue
		}

		// Trust
		trust := e.Config.Simulation.DefaultRelationshipTrust
		if rels, has := e.Relationships[speaker.ID]; has {
			if t, okT := rels[listener.ID]; okT {
				trust = t
			}
		}

		// Check listener direct contradiction or prior provenance
		listenerBelief, hasBelief := e.Beliefs[listener.ID][inter.ClaimID]
		oldConf := 0.10
		if hasBelief {
			oldConf = listenerBelief.Confidence
		}

		hasContradiction := false
		hasSameRoot := false
		if hasBelief && listenerBelief.RootSourceAgentID != nil && speakerBelief.RootSourceAgentID != nil {
			if *listenerBelief.RootSourceAgentID == *speakerBelief.RootSourceAgentID {
				hasSameRoot = true
			}
		}

		// Calculate updated belief
		newConf := belief.CalculatePropagation(belief.PropagationInput{
			SpeakerConfidence:        speakerBelief.Confidence,
			RelationshipTrust:        trust,
			SpeakerInfluence:         speaker.Influence,
			SpeakerCredibility:       speaker.Credibility,
			ListenerReceptiveness:     listener.Receptiveness,
			ListenerSkepticism:       listener.Skepticism,
			HasDirectContradiction:   hasContradiction,
			HasSameRootSource:        hasSameRoot,
			OldConfidence:            oldConf,
			BeliefUpdateRate:         e.Config.Simulation.BeliefUpdateRate,
			DirectEvidenceMultiplier: e.Config.Simulation.DirectEvidenceMultiplier,
			SameRootSourceMultiplier: e.Config.Simulation.SameRootSourceMultiplier,
		})

		rootSource := speaker.ID
		if speakerBelief.RootSourceAgentID != nil {
			rootSource = *speakerBelief.RootSourceAgentID
		}

		e.Beliefs[listener.ID][inter.ClaimID] = &domain.AgentBelief{
			SessionID:           e.Session.ID,
			AgentID:             listener.ID,
			ClaimID:             inter.ClaimID,
			Confidence:          newConf,
			FirstSourceAgentID:  &speaker.ID,
			LastSourceAgentID:   &speaker.ID,
			RootSourceAgentID:   &rootSource,
			UpdatedAtGameSecond: currentSecond,
			Revision:            1,
		}

		// Append memory to listener
		claimObj := e.Claims[inter.ClaimID]
		e.Memories[listener.ID] = append(e.Memories[listener.ID], domain.AgentMemory{
			ID:                  domain.NewUUID(),
			SessionID:           e.Session.ID,
			AgentID:             listener.ID,
			MemoryType:          domain.MemoryTypeHearsay,
			Content:             speaker.Name + " nói: " + claimObj.DisplayText,
			Reliability:         speaker.Credibility * trust,
			Importance:          inter.Salience,
			SourceAgentID:       &speaker.ID,
			RootSourceAgentID:   &rootSource,
			RelatedClaimIDs:     []string{inter.ClaimID},
			CreatedAtGameSecond: currentSecond,
		})

		// Emit rumor.visualized for 3D map pulse
		e.AppendEvent(
			domain.WorldEventTypeRumorVisualized,
			&speaker.ID, &listener.ID,
			RumorVisualizedPayload{
				SpeakerAgentID:  speaker.ID,
				ListenerAgentID: listener.ID,
				SpeakerName:     speaker.Name,
				ListenerName:    listener.Name,
				LocationID:      speaker.CurrentLocationID,
				Category:        string(claimObj.NarrativeRole),
			},
		)
	}

	// 4. Update overall false narrative ratio & check loss
	e.recalculateNarrativeRatio()

	e.AppendEvent(
		domain.WorldEventTypeBeliefStatsUpdated,
		nil, nil,
		BeliefStatsPayload{
			PrimaryFalseNarrativeRatio: e.Session.FalseNarrativeRatio,
			EvidenceStrength:           e.Session.EvidenceStrength,
			PlayerCredibility:          e.Session.PlayerCredibility,
		},
	)

	// If game was lost during this tick
	if e.Session.Status == domain.SessionStatusLostFalseBelief {
		e.AppendEvent(
			domain.WorldEventTypeGameLost,
			nil, nil,
			map[string]interface{}{
				"reason":           "Primary false narrative achieved social dominance (>= 75%)",
				"false_ratio":      e.Session.FalseNarrativeRatio,
				"evidence_strength": e.Session.EvidenceStrength,
			},
		)
	}
}

func (e *SimulationEngine) seedAgentBelief(scenarioAgentID string, claimID string, confidence float64, memoryContent string) {
	agent, okA := e.Agents[scenarioAgentID]
	claim, okC := e.Claims[claimID]
	if !okA || !okC {
		return
	}

	b, exists := e.Beliefs[agent.ID][claim.ScenarioClaimID]
	if !exists {
		e.Beliefs[agent.ID][claim.ScenarioClaimID] = &domain.AgentBelief{
			SessionID:           e.Session.ID,
			AgentID:             agent.ID,
			ClaimID:             claim.ScenarioClaimID,
			Confidence:          confidence,
			UpdatedAtGameSecond: e.Session.GameSecond,
			Revision:            1,
		}
	} else if b.Confidence < confidence {
		b.Confidence = confidence
		b.UpdatedAtGameSecond = e.Session.GameSecond
		b.Revision++
	}

	e.Memories[agent.ID] = append(e.Memories[agent.ID], domain.AgentMemory{
		ID:                  domain.NewUUID(),
		SessionID:           e.Session.ID,
		AgentID:             agent.ID,
		MemoryType:          domain.MemoryTypeHearsay,
		Content:             memoryContent,
		Reliability:         0.70,
		Importance:          0.85,
		RelatedClaimIDs:     []string{claim.ScenarioClaimID},
		CreatedAtGameSecond: e.Session.GameSecond,
	})
}

