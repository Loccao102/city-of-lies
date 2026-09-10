-- 000001_init.up.sql
create table if not exists game_sessions (
    id uuid primary key,
    scenario_id text not null,
    status text not null,
    simulation_seed bigint not null,
    game_second bigint not null default 0,
    player_credibility double precision not null check (player_credibility between 0 and 1),
    false_narrative_ratio double precision not null default 0 check (false_narrative_ratio between 0 and 1),
    evidence_strength double precision not null default 0 check (evidence_strength between 0 and 1),
    revision bigint not null default 0,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    ended_at timestamptz null
);

create table if not exists session_agents (
    id uuid primary key,
    session_id uuid not null references game_sessions(id) on delete cascade,
    scenario_agent_id text not null,
    name text not null,
    role text not null,
    current_location_id text not null,
    influence double precision not null check (influence between 0 and 1),
    credibility double precision not null check (credibility between 0 and 1),
    skepticism double precision not null check (skepticism between 0 and 1),
    sociality double precision not null check (sociality between 0 and 1),
    deception_tendency double precision not null check (deception_tendency between 0 and 1),
    receptiveness double precision not null check (receptiveness between 0 and 1),
    busy_until_game_second bigint not null default 0,
    revision bigint not null default 0,
    unique(session_id, scenario_agent_id)
);

create table if not exists session_claims (
    id uuid primary key,
    session_id uuid not null references game_sessions(id) on delete cascade,
    scenario_claim_id text not null,
    display_text text not null,
    classification text not null,
    narrative_role text not null,
    narrative_weight double precision not null default 0,
    unique(session_id, scenario_claim_id)
);

create table if not exists agent_beliefs (
    session_id uuid not null references game_sessions(id) on delete cascade,
    agent_id uuid not null references session_agents(id) on delete cascade,
    claim_id uuid not null references session_claims(id) on delete cascade,
    confidence double precision not null check (confidence between 0 and 1),
    first_source_agent_id uuid null references session_agents(id),
    last_source_agent_id uuid null references session_agents(id),
    root_source_agent_id uuid null references session_agents(id),
    evidence_resistance double precision not null default 0,
    updated_at_game_second bigint not null,
    revision bigint not null default 0,
    primary key(agent_id, claim_id)
);

create table if not exists agent_memories (
    id uuid primary key,
    session_id uuid not null references game_sessions(id) on delete cascade,
    agent_id uuid not null references session_agents(id) on delete cascade,
    memory_type text not null,
    content text not null,
    reliability double precision not null check (reliability between 0 and 1),
    importance double precision not null check (importance between 0 and 1),
    source_agent_id uuid null references session_agents(id),
    root_source_agent_id uuid null references session_agents(id),
    related_claim_ids uuid[] not null default '{}',
    created_at_game_second bigint not null
);

create table if not exists relationships (
    session_id uuid not null references game_sessions(id) on delete cascade,
    from_agent_id uuid not null references session_agents(id) on delete cascade,
    to_agent_id uuid not null references session_agents(id) on delete cascade,
    trust double precision not null check (trust between 0 and 1),
    revision bigint not null default 0,
    primary key(from_agent_id, to_agent_id)
);

create table if not exists evidence_items (
    id uuid primary key,
    session_id uuid not null references game_sessions(id) on delete cascade,
    scenario_evidence_id text not null,
    name text not null,
    description text not null,
    location_id text not null,
    reliability double precision not null check (reliability between 0 and 1),
    source_category text not null,
    discovered boolean not null default false,
    discovered_at_game_second bigint null,
    unique(session_id, scenario_evidence_id)
);

create table if not exists evidence_claim_links (
    evidence_id uuid not null references evidence_items(id) on delete cascade,
    claim_id uuid not null references session_claims(id) on delete cascade,
    relation text not null check (relation in ('supports','contradicts')),
    weight double precision not null default 1 check (weight between 0 and 1),
    primary key(evidence_id, claim_id, relation)
);

create table if not exists conversations (
    id uuid primary key,
    session_id uuid not null references game_sessions(id) on delete cascade,
    agent_id uuid not null references session_agents(id),
    started_at_game_second bigint not null,
    ended_at_game_second bigint null
);

create table if not exists conversation_messages (
    id uuid primary key,
    conversation_id uuid not null references conversations(id) on delete cascade,
    speaker_type text not null check (speaker_type in ('player','agent','system')),
    content text not null,
    intent text null,
    created_at_game_second bigint not null,
    model_provider text null,
    model_name text null,
    latency_ms integer null
);

create table if not exists publications (
    id uuid primary key,
    session_id uuid not null references game_sessions(id) on delete cascade,
    challenged_claim_id uuid not null references session_claims(id),
    message text not null,
    strength double precision not null check (strength between 0 and 1),
    player_credibility_before double precision not null,
    player_credibility_after double precision not null,
    published_at_game_second bigint not null
);

create table if not exists publication_evidence (
    publication_id uuid not null references publications(id) on delete cascade,
    evidence_id uuid not null references evidence_items(id),
    primary key(publication_id, evidence_id)
);

create table if not exists world_events (
    id uuid primary key,
    session_id uuid not null references game_sessions(id) on delete cascade,
    sequence bigint not null,
    event_type text not null,
    game_second bigint not null,
    actor_agent_id uuid null references session_agents(id),
    target_agent_id uuid null references session_agents(id),
    payload jsonb not null default '{}'::jsonb,
    occurred_at timestamptz not null default now(),
    unique(session_id, sequence)
);
