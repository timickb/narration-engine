create table blogs (
    id                  uuid not null primary key,
    author_id           uuid not null,
    author_email        text not null,
    name                text not null,
    subscribers_count   integer not null,
    publications_count  integer not null,
    donations_count     integer not null,
    created_at          timestamptz not null default now(),
    updated_at          timestamptz not null default now()
);

create table publications (
    id                  uuid not null primary key,
    blog_id             uuid not null references blogs(id),
    author_id           uuid not null,
    title               text not null,
    body                text not null default '',
    status              text not null,
    created_at          timestamptz not null default now(),
    updated_at          timestamptz not null default now()
);