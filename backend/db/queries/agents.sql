-- name: CreateSessionAgent :one
insert into session_agents (
    id, session_id, scenario_agent_id, name, role,
    current_location_id, influence, credibility, skepticism,
    sociality, deception_tendency, receptiveness, busy_until_game_second
) values (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
)
returning *;

-- name: ListSessionAgents :many
select * from session_agents
where session_id = $1
order by name;

-- name: GetSessionAgent :one
select * from session_agents
where session_id = $1 and id = $2;

-- name: GetSessionAgentByScenarioId :one
select * from session_agents
where session_id = $1 and scenario_agent_id = $2;

-- name: UpdateAgentLocation :exec
update session_agents
set current_location_id = $3, revision = revision + 1
where session_id = $1 and id = $2;

-- name: UpdateAgentBusyUntil :exec
update session_agents
set busy_until_game_second = $3, revision = revision + 1
where session_id = $1 and id = $2;
