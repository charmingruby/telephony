CREATE TABLE IF NOT EXISTS users
(
    id bigserial PRIMARY KEY NOT NULL,
    uuid varchar NOT NULL,
    first_name varchar NOT NULL,
    last_name varchar NOT NULL,
    email varchar NOT NULL,
    password_hash varchar NOT NULL,
    created_at timestamp DEFAULT now() NOT NULL,
    updated_at timestamp DEFAULT now() NOT NULL,
    deleted_at timestamp
);

CREATE UNIQUE INDEX users_uuid_uindex
    ON users (uuid);

CREATE UNIQUE INDEX users_email_uindex
    ON users (email);

CREATE TABLE IF NOT EXISTS users_profile 
(
    id bigserial PRIMARY KEY NOT NULL,
    uuid varchar NOT NULL,
    display_name varchar NOT NULL,
    bio varchar NOT NULL,
    guilds_quantity integer NOT NULL,
    messages_quantity integer NOT NULL,
    user_id integer REFERENCES users (id),
    created_at timestamp DEFAULT now() NOT NULL,
    updated_at timestamp DEFAULT now() NOT NULL,
    deleted_at timestamp
);

CREATE UNIQUE INDEX users_profile_display_name_uindex
    ON users_profile (display_name);


CREATE TABLE IF NOT EXISTS examples
(
    id varchar PRIMARY KEY NOT NULL,
    name varchar NOT NULL,
    created_at timestamp DEFAULT now() NOT NULL
);