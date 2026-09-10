-- name: CreateSession :one
insert into game_sessions (
    id, scenario_id, status, simulation_seed, game_second,
    player_credibility, false_narrative_ratio, evidence_strength, revision
) values (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
returning *;

-- name: GetSession :one
select * from game_sessions
where id = $1;

-- name: UpdateSessionState :exec
update game_sessions
set
    status = $2,
    game_second = $3,
    player_credibility = $4,
    false_narrative_ratio = $5,
    evidence_strength = $6,
    revision = revision + 1,
    updated_at = now(),
    ended_at = case when $7::boolean then now() else ended_at end
where id = $1;

-- name: ListActiveSessions :many
select * from game_sessions
where status = 'running';
