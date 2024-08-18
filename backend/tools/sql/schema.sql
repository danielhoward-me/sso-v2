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
	id UUID NOT NULL PRIMARY KEY,
	username TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL,
	email VARCHAR(256) NOT NULL UNIQUE,
	profile_picture profile_picture NOT NULL DEFAULT 'wavatar',
	created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	last_updated TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TRIGGER update_user_last_updated BEFORE UPDATE ON users FOR EACH ROW EXECUTE PROCEDURE update_last_updated_column();

CREATE TABLE clients (
	id UUID NOT NULL PRIMARY KEY,
	name TEXT NOT NULL UNIQUE,
	secret TEXT NOT NULL,
	show_confirmation_prompt BOOLEAN NOT NULL
);
CREATE TABLE client_redirects (
	client_id UUID NOT NULL REFERENCES clients(id),
	redirect TEXT NOT NULL,
	PRIMARY KEY (client_id, redirect)
);
