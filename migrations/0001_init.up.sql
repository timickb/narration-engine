-- Экземпляры сценариев
create table instances (
    id                  uuid not null primary key,
    parent_id           uuid,

    scenario_name       text not null,
    scenario_version    text not null,
    current_state       text not null,
    previous_state      text not null,
    context             jsonb not null,
    blocking_key        text,
    retries             integer not null,
    failed              boolean not null,

    created_at          timestamptz not null default now(),
    updated_at          timestamptz not null default now(),

    constraint instance_instances_fk
        foreign key(parent_id) references instances(id)
);

-- История переходов между состояниями экземпляров
create table transitions (
    id                   uuid not null primary key,
    instance_id          uuid not null,

    state_from           text not null,
    state_to             text not null,
    event_name           text not null,
    params               jsonb not null,
    failed               boolean not null,
    error                text,

    created_at           timestamptz not null default now(),

    constraint transitions_instances_fk
        foreign key(instance_id) references instances(id)
);

-- Очередь событий на выполнение
create table pending_events (
    id                   uuid not null primary key,
    instance_id          uuid not null,
    event_name           text not null,
    params               jsonb not null,
    external             bool not null,
    created_at           timestamptz not null default now(),
    executed_at          timestamptz
);