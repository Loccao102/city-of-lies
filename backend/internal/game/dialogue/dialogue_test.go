package dialogue

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"city-of-lies/backend/internal/domain"
)

func TestDialogueGuardAndTemplate(t *testing.T) {
	guard := NewKnowledgeGuard()
	template := NewTemplateProvider()

	var dCtx domain.DialogueContext
	dCtx.Identity.Name = "Minh Trần"
	dCtx.Identity.Role = "Bảo vệ"
	dCtx.KnownClaims = []domain.AgentKnownClaim{
		{ClaimID: "claim_fire_occurred", Confidence: 0.90, Basis: "direct"},
	}
	dCtx.AllowedRevealEvidenceIDs = []string{"ev_cctv_gate"}

	resp, err := template.GenerateDialogue(context.Background(), dCtx, "Có chuyện gì vậy?")
	if err != nil {
		t.Fatalf("expected template generation to succeed, got %v", err)
	}

	if err := guard.ValidateResponse(resp, dCtx); err != nil {
		t.Errorf("expected guard validation to pass for template response, got %v", err)
	}

	// Test guard rejecting unknown claim
	badResp := resp
	badResp.ReferencedClaimIDs = []string{"claim_chemical_explosion"}
	if err := guard.ValidateResponse(badResp, dCtx); err == nil {
		t.Errorf("expected guard to reject response referencing unknown claim")
	}

	// Test guard rejecting disallowed evidence
	badEvResp := resp
	badEvResp.RevealedEvidenceIDs = []string{"ev_hospital_summary"}
	if err := guard.ValidateResponse(badEvResp, dCtx); err == nil {
		t.Errorf("expected guard to reject response revealing unpermitted evidence")
	}
}

func TestDynamicLLMProviderAutoConfig(t *testing.T) {
	providers := []struct {
		name        string
		expectedURL string
		model       string
	}{
		{"gemini", "https://generativelanguage.googleapis.com/v1beta/openai", "gemini-2.0-flash"},
		{"groq", "https://api.groq.com/openai/v1", "llama-3.3-70b-versatile"},
		{"ollama", "http://localhost:11434/v1", "llama3.2"},
		{"openrouter", "https://openrouter.ai/api/v1", "openai/gpt-4o-mini"},
		{"openai", "https://api.openai.com/v1", "gpt-4o-mini"},
	}

	for _, tc := range providers {
		p := NewDynamicLLMProvider(tc.name, "", "dummy-key", "", 5*time.Second, 1)
		if p.baseURL != tc.expectedURL {
			t.Errorf("provider %s: expected baseURL %s, got %s", tc.name, tc.expectedURL, p.baseURL)
		}
		if p.model != tc.model {
			t.Errorf("provider %s: expected model %s, got %s", tc.name, tc.model, p.model)
		}
	}
}

func TestDynamicLLMProviderSuccessAndCodeFences(t *testing.T) {
	// Mock server that returns markdown-wrapped JSON response
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-api-key" {
			t.Errorf("expected Bearer test-api-key, got %s", r.Header.Get("Authorization"))
		}

		mockResp := map[string]interface{}{
			"choices": []map[string]interface{}{
				{
					"message": map[string]interface{}{
						"content": "```json\n{\n  \"intent\": \"answer\",\n  \"utterance\": \"Tôi thấy khói bốc lên từ cổng bảo vệ lúc 8 giờ.\",\n  \"referenced_claim_ids\": [\"claim_fire_occurred\"],\n  \"revealed_evidence_ids\": [\"ev_cctv_gate\"],\n  \"emotion\": \"calm\",\n  \"certainty\": 0.95\n}\n```",
					},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(mockResp)
	}))
	defer mockServer.Close()

	provider := NewDynamicLLMProvider("openai", mockServer.URL, "test-api-key", "test-model", 2*time.Second, 1)

	var dCtx domain.DialogueContext
	dCtx.Identity.Name = "Bảo vệ Nam"
	dCtx.Identity.Role = "Bảo vệ"
	dCtx.KnownClaims = []domain.AgentKnownClaim{
		{ClaimID: "claim_fire_occurred", Confidence: 0.95, Basis: "direct"},
	}
	dCtx.AllowedRevealEvidenceIDs = []string{"ev_cctv_gate"}

	resp, err := provider.GenerateDialogue(context.Background(), dCtx, "Anh có thấy khói không?")
	if err != nil {
		t.Fatalf("expected successful dialogue generation, got %v", err)
	}

	if resp.Intent != domain.DialogueIntentAnswer {
		t.Errorf("expected intent answer, got %s", resp.Intent)
	}
	if !strings.Contains(resp.Utterance, "Tôi thấy khói bốc lên") {
		t.Errorf("unexpected utterance: %s", resp.Utterance)
	}
	if len(resp.ReferencedClaimIDs) != 1 || resp.ReferencedClaimIDs[0] != "claim_fire_occurred" {
		t.Errorf("expected claim_fire_occurred, got %v", resp.ReferencedClaimIDs)
	}
	if len(resp.RevealedEvidenceIDs) != 1 || resp.RevealedEvidenceIDs[0] != "ev_cctv_gate" {
		t.Errorf("expected ev_cctv_gate, got %v", resp.RevealedEvidenceIDs)
	}
}

func TestDynamicLLMProviderGuardRejectionFallback(t *testing.T) {
	// Mock server that returns a response with hallucinated claim
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mockResp := map[string]interface{}{
			"choices": []map[string]interface{}{
				{
					"message": map[string]interface{}{
						"content": `{"intent": "answer", "utterance": "Người ngoài hành tinh tấn công!", "referenced_claim_ids": ["claim_alien_attack"], "revealed_evidence_ids": [], "emotion": "afraid", "certainty": 0.99}`,
					},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(mockResp)
	}))
	defer mockServer.Close()

	provider := NewDynamicLLMProvider("openai", mockServer.URL, "test-api-key", "test-model", 2*time.Second, 1)

	var dCtx domain.DialogueContext
	dCtx.Identity.Name = "Bảo vệ Nam"
	dCtx.Identity.Role = "Bảo vệ"
	dCtx.KnownClaims = []domain.AgentKnownClaim{
		{ClaimID: "claim_fire_occurred", Confidence: 0.95, Basis: "direct"},
	}

	resp, err := provider.GenerateDialogue(context.Background(), dCtx, "Chuyện gì vậy?")
	if err != nil {
		t.Fatalf("expected fallback generation without error, got %v", err)
	}

	// Should have fallen back to template provider
	for _, claim := range resp.ReferencedClaimIDs {
		if claim == "claim_alien_attack" {
			t.Fatalf("guard failed to reject hallucinated claim!")
		}
	}
}

func TestDynamicLLMProviderNetworkFailureFallback(t *testing.T) {
	// Mock server that returns HTTP 500
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}))
	defer mockServer.Close()

	provider := NewDynamicLLMProvider("openai", mockServer.URL, "test-api-key", "test-model", 1*time.Second, 0)

	var dCtx domain.DialogueContext
	dCtx.Identity.Name = "Bác sĩ Nghiêm"
	dCtx.Identity.Role = "Bác sĩ"
	dCtx.KnownClaims = []domain.AgentKnownClaim{
		{ClaimID: "claim_zero_fatalities", Confidence: 0.90, Basis: "direct"},
	}

	resp, err := provider.GenerateDialogue(context.Background(), dCtx, "Có ai chết không?")
	if err != nil {
		t.Fatalf("expected fallback to succeed even with HTTP 500, got err: %v", err)
	}

	if !strings.Contains(resp.Utterance, "không có ca tử vong nào") {
		t.Errorf("expected template fallback response about fatalities, got %s", resp.Utterance)
	}
}
