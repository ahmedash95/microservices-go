CREATE TABLE comments (
   id serial PRIMARY KEY,
   post_id INT NOT NULL CHECK (post_id >= 0),
   user_id INT NOT NULL CHECK (user_id >= 0),
   comment TEXT NOT NULL,
   created_at TIMESTAMP NOT NULL,
   updated_at TIMESTAMP,
   deleted_at TIMESTAMP
);