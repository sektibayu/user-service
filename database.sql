CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
	id TEXT PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
	full_name VARCHAR (50) NOT NULL,
  phone_number VARCHAR (50) UNIQUE NOT NULL,
  hash_password TEXT NOT NULL,
  salt VARCHAR (12) NOT NULL,
  login_success_count INT NOT NULL DEFAULT 0,
  created_at TIMESTAMP DEFAULT current_timestamp NOT NULL,
  updated_at TIMESTAMP DEFAULT current_timestamp NOT NULL
);

CREATE UNIQUE INDEX ON users (phone_number);
