CREATE TABLE users (
   id serial PRIMARY KEY,
   name VARCHAR (50) NOT NULL,
   password VARCHAR (250) NOT NULL,
   email VARCHAR (355) UNIQUE NOT NULL,
   created_at TIMESTAMP NOT NULL,
   updated_at TIMESTAMP,
   deleted_at TIMESTAMP,
   last_login TIMESTAMP
);