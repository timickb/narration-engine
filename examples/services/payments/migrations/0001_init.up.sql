create table account (
    id                  uuid not null primary key,
    user_id             uuid not null,
    amount              decimal,
    created_at          timestamptz not null default now(),
    last_payout_at      timestamptz
);

create table invoice (
    id                  uuid not null primary key,
    description         text not null default '',
    amount              decimal,
    completed           bool,
    failed              bool,
    created_at          timestamptz not null default now()
);