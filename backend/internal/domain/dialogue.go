package domain

type DialogueIntent string

const (
	DialogueIntentAnswer      DialogueIntent = "answer"
	DialogueIntentDeflect     DialogueIntent = "deflect"
	DialogueIntentLie         DialogueIntent = "lie"
	DialogueIntentRefuse      DialogueIntent = "refuse"
	DialogueIntentAskQuestion DialogueIntent = "ask_question"
)

type DialogueEmotion string

const (
	DialogueEmotionNeutral   DialogueEmotion = "neutral"
	DialogueEmotionCalm      DialogueEmotion = "calm"
	DialogueEmotionUncertain DialogueEmotion = "uncertain"
	DialogueEmotionAfraid    DialogueEmotion = "afraid"
	DialogueEmotionAngry     DialogueEmotion = "angry"
	DialogueEmotionExcited   DialogueEmotion = "excited"
	DialogueEmotionDefensive DialogueEmotion = "defensive"
	DialogueEmotionSad       DialogueEmotion = "sad"
)

// DialogueResponse is the strictly enforced structured response from a dialogue provider.
type DialogueResponse struct {
	Intent             DialogueIntent  `json:"intent"`
	Utterance          string          `json:"utterance"`
	ReferencedClaimIDs []string        `json:"referenced_claim_ids"`
	RevealedEvidenceIDs []string       `json:"revealed_evidence_ids"`
	Emotion            DialogueEmotion `json:"emotion"`
	Certainty          float64         `json:"certainty"`
}

type AgentKnownClaim struct {
	ClaimID    string  `json:"claim_id"`
	Confidence float64 `json:"confidence"`
	Basis      string  `json:"basis"`
}

type MemoryBrief struct {
	Type    string  `json:"type"`
	Content string  `json:"content"`
}

type ChatMessage struct {
	Speaker string `json:"speaker"`
	Text    string `json:"text"`
}

// DialogueContext represents the privacy-safe, knowledge-guarded context passed to dialogue generation.
type DialogueContext struct {
	Identity struct {
		Name string `json:"name"`
		Role string `json:"role"`
	} `json:"identity"`
	Traits struct {
		Skepticism        float64 `json:"skepticism"`
		DeceptionTendency float64 `json:"deception_tendency"`
	} `json:"traits"`
	CurrentLocation          string            `json:"current_location"`
	KnownClaims              []AgentKnownClaim `json:"known_claims"`
	RelevantMemories         []MemoryBrief     `json:"relevant_memories"`
	RelationshipToPlayer     float64           `json:"relationship_to_player"`
	ConversationHistory      []ChatMessage     `json:"conversation_history"`
	AllowedRevealEvidenceIDs []string          `json:"allowed_reveal_evidence_ids"`
}
