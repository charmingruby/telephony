CREATE TABLE IF NOT EXISTS guilds
(
    id bigserial PRIMARY KEY NOT NULL,
    name varchar NOT NULL,
    description varchar NOT NULL,
    tags varchar[] NOT NULL,
    channels_quantity integer NOT NULL,
    owner_id integer REFERENCES users_profile (id),
    created_at timestamp DEFAULT now() NOT NULL,
    updated_at timestamp DEFAULT now() NOT NULL,
    deleted_at timestamp
);