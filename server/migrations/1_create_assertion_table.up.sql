CREATE TABLE IF NOT EXISTS assertions  (
	id UUID NOT NULL PRIMARY KEY,
	test_id UUID NOT NULL,
	assertion json NOT NULL
);