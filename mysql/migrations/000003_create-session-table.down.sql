USE snippetbox;

-- Drop the index first
DROP INDEX sessions_expiry_idx ON sessions;

-- Drop the table
DROP TABLE sessions;