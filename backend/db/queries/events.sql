-- name: AppendWorldEvent :one
insert into world_events (
    id, session_id, sequence, event_type, game_second,
    actor_agent_id, target_agent_id, payload
) values (
    $1, $2, $3, $4, $5, $6, $7, $8
)
returning *;

-- name: ListWorldEventsSince :many
select * from world_events
where session_id = $1 and sequence > $2
order by sequence asc
limit $3;

-- name: GetLatestEventSequence :one
select coalesce(max(sequence), 0)::bigint as max_sequence
from world_events
where session_id = $1;
