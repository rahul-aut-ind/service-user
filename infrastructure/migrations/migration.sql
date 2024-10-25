# Migration for creating the users table and populating with 10 users
CREATE TABLE IF NOT EXISTS users (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255),
    email VARCHAR(255) UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
# Insert 10 users
INSERT INTO users (name, email) VALUES
  ('User1', 'user1@example.com'),
  ('User2', 'user2@example.com'),
  ('User3', 'user3@example.com'),
  ('User4', 'user4@example.com'),
  ('User5', 'user5@example.com'),
  ('User6', 'user6@example.com'),
  ('User7', 'user7@example.com'),
  ('User8', 'user8@example.com'),
  ('User9', 'user9@example.com'),
  ('User10', 'user10@example.com')
;