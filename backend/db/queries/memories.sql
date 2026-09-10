-- name: CreateAgentMemory :one
insert into agent_memories (
    id, session_id, agent_id, memory_type, content,
    reliability, importance, source_agent_id, root_source_agent_id,
    related_claim_ids, created_at_game_second
) values (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
)
returning *;

-- name: ListAgentMemories :many
select * from agent_memories
where session_id = $1 and agent_id = $2
order by created_at_game_second desc
limit $3;

-- name: UpsertRelationship :exec
insert into relationships (
    session_id, from_agent_id, to_agent_id, trust
) values (
    $1, $2, $3, $4
)
on conflict (from_agent_id, to_agent_id) do update set
    trust = excluded.trust,
    revision = relationships.revision + 1;

-- name: GetRelationship :one
select * from relationships
where session_id = $1 and from_agent_id = $2 and to_agent_id = $3;

-- name: ListRelationships :many
select * from relationships
where session_id = $1;
