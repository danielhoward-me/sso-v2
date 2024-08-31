CREATE OR REPLACE FUNCTION update_last_updated_column() 
RETURNS TRIGGER AS $$
BEGIN
    NEW.last_updated = now();
    RETURN NEW; 
END;
$$ language 'plpgsql';

-- Should match server/types.d.ts ProfilePictureType
CREATE TYPE profile_picture AS ENUM('custom', 'identicon', 'monsterid', 'wavatar', 'retro', 'robohash');
CREATE TABLE users (
	id SERIAL PRIMARY KEY,
	uuid UUID NOT NULL UNIQUE,
	username TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL,
	email TEXT NOT NULL UNIQUE,
	profile_picture profile_picture NOT NULL DEFAULT 'wavatar',
	created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	last_updated TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TRIGGER update_user_last_updated BEFORE UPDATE ON users FOR EACH ROW EXECUTE PROCEDURE update_last_updated_column();

CREATE TABLE clients (
	id SERIAL PRIMARY KEY,
	uuid UUID NOT NULL UNIQUE,
	name TEXT NOT NULL UNIQUE,
	secret TEXT NOT NULL,
	show_confirmation_prompt BOOLEAN NOT NULL
);
CREATE TABLE client_redirects (
	client_id INTEGER NOT NULL REFERENCES clients(id),
	redirect TEXT NOT NULL,
	PRIMARY KEY (client_id, redirect)
);

CREATE TABLE auth_codes (
	id SERIAL PRIMARY KEY,
	code TEXT NOT NULL UNIQUE,
	client_id INTEGER NOT NULL REFERENCES clients(id),
	user_id INTEGER NOT NULL REFERENCES users(id),
	redirect_uri TEXT NOT NULL,
	created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	expires TIMESTAMP NOT NULL
);
CREATE TABLE refresh_tokens (
	id SERIAL PRIMARY KEY,
	token TEXT NOT NULL UNIQUE,
	client_id INTEGER NOT NULL REFERENCES clients(id),
	user_id INTEGER NOT NULL REFERENCES users(id),
	created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	expires TIMESTAMP NOT NULL
);
CREATE TABLE access_tokens (
	id SERIAL PRIMARY KEY,
	token TEXT NOT NULL UNIQUE,
	refresh_token_id INTEGER NOT NULL REFERENCES refresh_tokens(id),
	created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	expires TIMESTAMP NOT NULL,
	last_used TIMESTAMP
);
