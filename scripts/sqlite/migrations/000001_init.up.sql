
/* Database creation and initialization */
create table media
(
    id          TEXT not null
        constraint media_pk
            primary key,
    filepath    TEXT not null,
    description TEXT,
    created     INTEGER default current_timestamp,
    tags        TEXT /* JSON string array */
);

