-- name: CreateConversation :one
insert into conversations (
    id, session_id, agent_id, started_at_game_second
) values (
    $1, $2, $3, $4
)
returning *;

-- name: GetConversationByAgent :one
select * from conversations
where session_id = $1 and agent_id = $2
order by started_at_game_second desc
limit 1;

-- name: CreateConversationMessage :one
insert into conversation_messages (
    id, conversation_id, speaker_type, content, intent,
    created_at_game_second, model_provider, model_name, latency_ms
) values (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
returning *;

-- name: ListConversationMessages :many
select * from conversation_messages
where conversation_id = $1
order by created_at_game_second asc;

-- name: CreatePublication :one
insert into publications (
    id, session_id, challenged_claim_id, message, strength,
    player_credibility_before, player_credibility_after, published_at_game_second
) values (
    $1, $2, $3, $4, $5, $6, $7, $8
)
returning *;

-- name: CreatePublicationEvidence :exec
insert into publication_evidence (
    publication_id, evidence_id
) values (
    $1, $2
);

-- name: ListPublications :many
select p.*, sc.scenario_claim_id
from publications p
join session_claims sc on p.challenged_claim_id = sc.id
where p.session_id = $1
order by p.published_at_game_second desc;
