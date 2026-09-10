package dialogue

import (
	"context"
	"testing"

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
