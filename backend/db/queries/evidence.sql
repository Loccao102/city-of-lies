-- name: CreateEvidenceItem :one
insert into evidence_items (
    id, session_id, scenario_evidence_id, name, description,
    location_id, reliability, source_category, discovered, discovered_at_game_second
) values (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
returning *;

-- name: ListDiscoveredEvidence :many
select * from evidence_items
where session_id = $1 and discovered = true
order by discovered_at_game_second asc;

-- name: ListEvidenceByLocation :many
select * from evidence_items
where session_id = $1 and location_id = $2;

-- name: MarkEvidenceDiscovered :exec
update evidence_items
set discovered = true, discovered_at_game_second = $3
where session_id = $1 and id = $2;

-- name: CreateEvidenceClaimLink :exec
insert into evidence_claim_links (
    evidence_id, claim_id, relation, weight
) values (
    $1, $2, $3, $4
)
on conflict do nothing;

-- name: ListEvidenceClaimLinks :many
select ecl.*, ei.scenario_evidence_id, sc.scenario_claim_id
from evidence_claim_links ecl
join evidence_items ei on ecl.evidence_id = ei.id
join session_claims sc on ecl.claim_id = sc.id
where ei.session_id = $1;
