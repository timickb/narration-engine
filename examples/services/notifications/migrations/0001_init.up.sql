create table messages (
    id          uuid not null  primary key,
    mail_from   text not null,
    mail_to     text not null,
    subject     text not null,
    body        text not null
);