-- Seed: 5 users.
-- Idempotent: safe to re-apply, and runs once via docker-entrypoint-initdb.d.
INSERT INTO users (name, email) VALUES
  ('Abhi Patel', 'abhi@example.com'),
  ('Priya Shah', 'priya@example.com'),
  ('Rohan Mehta', 'rohan@example.com'),
  ('Sneha Iyer', 'sneha@example.com'),
  ('Arjun Nair', 'arjun@example.com')
ON CONFLICT (email) DO NOTHING;
