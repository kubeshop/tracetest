CREATE TABLE IF NOT EXISTS "server"  (
	id varchar(100) NOT NULL PRIMARY KEY,
  created_at timestamp NOT NULL DEFAULT NOW()
);
