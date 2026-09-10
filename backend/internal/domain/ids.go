package domain

import (
	"crypto/rand"
	"fmt"
	"strings"
)

// AgentID represents a unique agent identifier (UUID string).
type AgentID string

// ClaimID represents a unique claim identifier (UUID string).
type ClaimID string

// EvidenceID represents a unique evidence identifier (UUID string).
type EvidenceID string

// SessionID represents a unique session identifier (UUID string).
type SessionID string

// LocationID represents a scenario location string key.
type LocationID string

// NewUUID generates a new standard UUID v4 string using crypto/rand.
func NewUUID() string {
	var u [16]byte
	_, _ = rand.Read(u[:])
	u[6] = (u[6] & 0x0f) | 0x40 // Version 4
	u[8] = (u[8] & 0x3f) | 0x80 // RFC 4122 variant
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
}

// ValidateUUID checks if a string matches basic UUID format.
func ValidateUUID(id string) error {
	id = strings.TrimSpace(id)
	if len(id) != 36 {
		return fmt.Errorf("invalid uuid length %d (expected 36): %q", len(id), id)
	}
	if id[8] != '-' || id[13] != '-' || id[18] != '-' || id[23] != '-' {
		return fmt.Errorf("invalid uuid hyphens: %q", id)
	}
	return nil
}
