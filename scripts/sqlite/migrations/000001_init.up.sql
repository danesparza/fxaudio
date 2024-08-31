/* Turn on foreign key support (I can't believe we have to do this) */
/* More information: https://www.sqlite.org/foreignkeys.html */
PRAGMA foreign_keys = ON;

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

