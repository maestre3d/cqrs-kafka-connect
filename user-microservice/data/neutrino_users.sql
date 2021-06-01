CREATE TABLE IF NOT EXISTS users (
	user_id varchar(128) NOT NULL,
	username varchar(64) NOT NULL,
	display_name varchar(128) NOT NULL,
	PRIMARY KEY(user_id)
);