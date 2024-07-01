CREATE TABLE IF NOT EXISTS guild_members
(
    id bigserial PRIMARY KEY NOT NULL,  
    guild_id integer REFERENCES guilds (id) NOT NULL,
    profile_id integer REFERENCES users_profile (id) NOT NULL,
    user_id integer REFERENCES users (id) NOT NULL,
    is_active bool NOT NULL,
    joined_at timestamp DEFAULT now() NOT NULL
);