package domain

import "errors"

var (
	ErrSessionNotFound      = errors.New("game session not found")
	ErrSessionNotRunning    = errors.New("session is not in running state")
	ErrAgentNotFound        = errors.New("agent not found")
	ErrAgentBusy            = errors.New("agent is busy")
	ErrLocationNotFound     = errors.New("location not found")
	ErrEvidenceNotFound     = errors.New("evidence not found")
	ErrEvidenceNotDiscovered = errors.New("evidence has not been discovered by player")
	ErrClaimNotFound        = errors.New("claim not found")
	ErrInvalidGroundTruth   = errors.New("invalid ground truth submission")
	ErrInvalidCorrection    = errors.New("invalid correction parameters")
	ErrScenarioNotFound     = errors.New("scenario not found")
	ErrInvalidConfig        = errors.New("invalid configuration")
	ErrUnauthorized         = errors.New("unauthorized action")
)
