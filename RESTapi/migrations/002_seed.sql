-- Seed: 5 users + ~12 todos (mixed completed flags).
-- Idempotent: safe to re-apply, and runs once via docker-entrypoint-initdb.d.
INSERT INTO users (name, email) VALUES
  ('Abhi Patel', 'abhi@example.com'),
  ('Priya Shah', 'priya@example.com'),
  ('Rohan Mehta', 'rohan@example.com'),
  ('Sneha Iyer', 'sneha@example.com'),
  ('Arjun Nair', 'arjun@example.com')
ON CONFLICT (email) DO NOTHING;

-- Resolve ids by email so seed is stable regardless of SERIAL values.
WITH u AS (
  SELECT id, email FROM users WHERE email IN (
    'abhi@example.com','priya@example.com','rohan@example.com',
    'sneha@example.com','arjun@example.com'
  )
)
INSERT INTO todos (user_id, title, completed)
SELECT u.id, t.title, t.completed
FROM u JOIN (VALUES
  ('abhi@example.com', 'Learn Go net package', false),
  ('abhi@example.com', 'Build custom HTTP parser', true),
  ('abhi@example.com', 'Wire Postgres store', false),
  ('priya@example.com', 'Design REST routes', false),
  ('priya@example.com', 'Write API docs', true),
  ('rohan@example.com', 'Set up Docker Compose', true),
  ('rohan@example.com', 'Add healthchecks', false),
  ('sneha@example.com', 'Test query filters', false),
  ('sneha@example.com', 'Review PATCH semantics', false),
  ('arjun@example.com', 'Seed dummy data', true),
  ('arjun@example.com', 'Verify cascade delete', false),
  ('arjun@example.com', 'Cleanup containers', false)
) AS t(email, title, completed) ON t.email = u.email
WHERE NOT EXISTS (
  SELECT 1 FROM todos d WHERE d.user_id = u.id AND d.title = t.title
);
