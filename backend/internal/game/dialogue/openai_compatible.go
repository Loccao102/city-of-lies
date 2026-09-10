package dialogue

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"city-of-lies/backend/internal/domain"
)

type OpenAIProvider struct {
	baseURL    string
	apiKey     string
	model      string
	httpClient *http.Client
	guard      *KnowledgeGuard
	fallback   *TemplateProvider
	retryCount int
}

func NewOpenAIProvider(baseURL, apiKey, model string, timeout time.Duration, retryCount int) *OpenAIProvider {
	if timeout <= 0 {
		timeout = 15 * time.Second
	}
	return &OpenAIProvider{
		baseURL: baseURL,
		apiKey:  apiKey,
		model:   model,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		guard:      NewKnowledgeGuard(),
		fallback:   NewTemplateProvider(),
		retryCount: retryCount,
	}
}

type openAIChatRequest struct {
	Model          string              `json:"model"`
	Messages       []openAIChatMessage `json:"messages"`
	ResponseFormat map[string]string   `json:"response_format,omitempty"`
}

type openAIChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func (p *OpenAIProvider) GenerateDialogue(ctx context.Context, dCtx domain.DialogueContext, playerQuery string) (domain.DialogueResponse, error) {
	if p.baseURL == "" || p.model == "" {
		return p.fallback.GenerateDialogue(ctx, dCtx, playerQuery)
	}

	systemPrompt := `You are an NPC in an investigation game set in Riverside district.
You must stay in character and speak natural Vietnamese.
CRITICAL SAFETY RULES:
1. Speak ONLY from the provided known claims and memories. NEVER invent unmentioned facts, casualties, or causes.
2. Output strictly a JSON object with this schema:
{
  "intent": "answer" | "deflect" | "lie" | "refuse" | "ask_question",
  "utterance": "string in Vietnamese",
  "referenced_claim_ids": ["string"],
  "revealed_evidence_ids": ["string"],
  "emotion": "neutral" | "calm" | "uncertain" | "afraid" | "angry" | "excited" | "defensive" | "sad",
  "certainty": float [0.0 - 1.0]
}`

	contextJSON, _ := json.Marshal(dCtx)
	userMessage := fmt.Sprintf("Character Context: %s\n\nPlayer asks: %s", string(contextJSON), playerQuery)

	reqBody := openAIChatRequest{
		Model: p.model,
		Messages: []openAIChatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userMessage},
		},
		ResponseFormat: map[string]string{"type": "json_object"},
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return p.fallback.GenerateDialogue(ctx, dCtx, playerQuery)
	}

	var lastErr error
	for attempt := 0; attempt <= p.retryCount; attempt++ {
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL+"/chat/completions", bytes.NewReader(data))
		if err != nil {
			lastErr = err
			continue
		}
		httpReq.Header.Set("Content-Type", "application/json")
		if p.apiKey != "" {
			httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)
		}

		resp, err := p.httpClient.Do(httpReq)
		if err != nil {
			lastErr = err
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("unexpected status %d", resp.StatusCode)
			continue
		}

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			lastErr = err
			continue
		}

		var chatResp openAIChatResponse
		if err := json.Unmarshal(bodyBytes, &chatResp); err != nil || len(chatResp.Choices) == 0 {
			lastErr = fmt.Errorf("failed to parse choices: %w", err)
			continue
		}

		var dialogueResp domain.DialogueResponse
		rawContent := chatResp.Choices[0].Message.Content
		if err := json.Unmarshal([]byte(rawContent), &dialogueResp); err != nil {
			lastErr = fmt.Errorf("model output is not valid DialogueResponse JSON: %w", err)
			continue
		}

		// Run knowledge guard
		if err := p.guard.ValidateResponse(dialogueResp, dCtx); err != nil {
			lastErr = fmt.Errorf("guard check failed: %w", err)
			continue
		}

		return dialogueResp, nil
	}

	// Fallback to deterministic template on error or guard rejection
	_ = lastErr
	return p.fallback.GenerateDialogue(ctx, dCtx, playerQuery)
}
