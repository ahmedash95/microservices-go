CREATE TABLE posts (
   id serial PRIMARY KEY,
   title VARCHAR (50) NOT NULL,
   content TEXT NOT NULL,
   is_draft BOOLEAN NOT NULL DEFAULT 't',
   image_url TEXT NOT NULL,
   publish_date TIMESTAMP NOT NULL,
   user_id INT NOT NULL CHECK (user_id >= 0),
   created_at TIMESTAMP NOT NULL,
   updated_at TIMESTAMP,
   deleted_at TIMESTAMP
);