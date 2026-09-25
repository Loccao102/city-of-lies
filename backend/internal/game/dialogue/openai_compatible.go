package dialogue

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"city-of-lies/backend/internal/domain"
)

type OpenAIProvider struct {
	provider   string
	baseURL    string
	apiKey     string
	model      string
	httpClient *http.Client
	guard      *KnowledgeGuard
	fallback   *TemplateProvider
	retryCount int
}

func NewOpenAIProvider(baseURL, apiKey, model string, timeout time.Duration, retryCount int) *OpenAIProvider {
	return NewDynamicLLMProvider("openai", baseURL, apiKey, model, timeout, retryCount)
}

// NewDynamicLLMProvider auto-configures defaults for OpenAI, Gemini, Groq, Ollama, OpenRouter.
func NewDynamicLLMProvider(provider, baseURL, apiKey, model string, timeout time.Duration, retryCount int) *OpenAIProvider {
	if timeout <= 0 {
		timeout = 15 * time.Second
	}
	if retryCount < 0 {
		retryCount = 1
	}

	normProvider := strings.ToLower(strings.TrimSpace(provider))
	normBaseURL := strings.TrimRight(strings.TrimSpace(baseURL), "/")

	// Auto-fill default endpoints and recommended fast models if omitted
	if normBaseURL == "" {
		switch normProvider {
		case "gemini":
			normBaseURL = "https://generativelanguage.googleapis.com/v1beta/openai"
			if model == "" {
				model = "gemini-2.0-flash"
			}
		case "groq":
			normBaseURL = "https://api.groq.com/openai/v1"
			if model == "" {
				model = "llama-3.3-70b-versatile"
			}
		case "ollama":
			normBaseURL = "http://localhost:11434/v1"
			if model == "" {
				model = "llama3.2"
			}
		case "openrouter":
			normBaseURL = "https://openrouter.ai/api/v1"
			if model == "" {
				model = "openai/gpt-4o-mini"
			}
		default: // "openai" or custom
			normBaseURL = "https://api.openai.com/v1"
			if model == "" {
				model = "gpt-4o-mini"
			}
		}
	}

	return &OpenAIProvider{
		provider: normProvider,
		baseURL:  normBaseURL,
		apiKey:   strings.TrimSpace(apiKey),
		model:    model,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		guard:      NewKnowledgeGuard(),
		fallback:   NewTemplateProvider(),
		retryCount: retryCount,
	}
}

func (p *OpenAIProvider) SetHTTPClient(client *http.Client) {
	if client != nil {
		p.httpClient = client
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

	systemPrompt := fmt.Sprintf(`You are roleplaying as NPC %q (Role: %q) in a city investigation game.
You MUST speak natural, authentic Vietnamese and stay strictly in character.

CRITICAL TRUTH & SAFETY RULES (NEVER VIOLATE):
1. Speak ONLY using facts, observations, and claims provided in your character context.
2. NEVER hallucinate or invent new casualties, dates, numbers, suspect names, causes, or locations not present in your known information.
3. If the player asks about something outside your known claims or memories, you must either truthfully state in Vietnamese that you do not know, or deflect naturally according to your role and personality.
4. "referenced_claim_ids" MUST contain ONLY claim IDs that are explicitly listed in your KnownClaims context.
5. "revealed_evidence_ids" MUST contain ONLY evidence IDs that are explicitly listed in your AllowedRevealEvidenceIDs context. If none are allowed, return an empty array [].
6. Output strictly a single valid JSON object adhering to this schema:
{
  "intent": "answer" | "deflect" | "lie" | "refuse" | "ask_question",
  "utterance": "string in natural Vietnamese",
  "referenced_claim_ids": ["string"],
  "revealed_evidence_ids": ["string"],
  "emotion": "neutral" | "calm" | "uncertain" | "afraid" | "angry" | "excited" | "defensive" | "sad",
  "certainty": float [0.0 - 1.0]
}`, dCtx.Identity.Name, dCtx.Identity.Role)

	contextJSON, _ := json.Marshal(dCtx)
	userMessage := fmt.Sprintf("Character Context:\n%s\n\nPlayer asks: %s", string(contextJSON), playerQuery)

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

		bodyBytes, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = err
			continue
		}

		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(bodyBytes))
			continue
		}

		var chatResp openAIChatResponse
		if err := json.Unmarshal(bodyBytes, &chatResp); err != nil || len(chatResp.Choices) == 0 {
			lastErr = fmt.Errorf("failed to parse choices: %w", err)
			continue
		}

		rawContent := strings.TrimSpace(chatResp.Choices[0].Message.Content)
		if strings.HasPrefix(rawContent, "```json") {
			rawContent = strings.TrimPrefix(rawContent, "```json")
			rawContent = strings.TrimSuffix(rawContent, "```")
			rawContent = strings.TrimSpace(rawContent)
		} else if strings.HasPrefix(rawContent, "```") {
			rawContent = strings.TrimPrefix(rawContent, "```")
			rawContent = strings.TrimSuffix(rawContent, "```")
			rawContent = strings.TrimSpace(rawContent)
		}

		var dialogueResp domain.DialogueResponse
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
