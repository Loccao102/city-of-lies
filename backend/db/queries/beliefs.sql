-- name: CreateSessionClaim :one
insert into session_claims (
    id, session_id, scenario_claim_id, display_text, classification, narrative_role, narrative_weight
) values (
    $1, $2, $3, $4, $5, $6, $7
)
returning *;

-- name: ListSessionClaims :many
select * from session_claims
where session_id = $1;

-- name: GetSessionClaimByScenarioId :one
select * from session_claims
where session_id = $1 and scenario_claim_id = $2;

-- name: UpsertAgentBelief :exec
insert into agent_beliefs (
    session_id, agent_id, claim_id, confidence,
    first_source_agent_id, last_source_agent_id, root_source_agent_id,
    evidence_resistance, updated_at_game_second
) values (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
on conflict (agent_id, claim_id) do update set
    confidence = excluded.confidence,
    last_source_agent_id = excluded.last_source_agent_id,
    root_source_agent_id = coalesce(agent_beliefs.root_source_agent_id, excluded.root_source_agent_id),
    evidence_resistance = excluded.evidence_resistance,
    updated_at_game_second = excluded.updated_at_game_second,
    revision = agent_beliefs.revision + 1;

-- name: ListAgentBeliefs :many
select * from agent_beliefs
where session_id = $1 and agent_id = $2;

-- name: ListAllSessionBeliefs :many
select ab.*, sc.scenario_claim_id, sc.narrative_role, sc.narrative_weight
from agent_beliefs ab
join session_claims sc on ab.claim_id = sc.id
where ab.session_id = $1;
