CREATE TABLE IF NOT EXISTS channels
(
    id bigserial PRIMARY KEY NOT NULL,
    name varchar NOT NULL,
    messages_quantity varchar NOT NULL,
    guild_id integer REFERENCES guilds (id),
    owner_id integer REFERENCES users_profile (id),
    created_at timestamp DEFAULT now() NOT NULL,
    updated_at timestamp DEFAULT now() NOT NULL,
    deleted_at timestamp
);

CREATE UNIQUE INDEX channels_name_uindex
    ON guilds (name);