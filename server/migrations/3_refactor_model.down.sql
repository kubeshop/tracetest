DROP TABLE definitions;
DROP TABLE runs;

CREATE TABLE IF NOT EXISTS assertions  (
	id UUID NOT NULL PRIMARY KEY,
	test_id UUID NOT NULL,
	assertion json NOT NULL
);

CREATE TABLE IF NOT EXISTS results  (
	id UUID NOT NULL PRIMARY KEY,
	test_id UUID NOT NULL,
	result json NOT NULL
);
