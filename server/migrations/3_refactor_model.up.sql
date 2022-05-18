DROP TABLE assertions;
DROP TABLE results;

CREATE TABLE IF NOT EXISTS definitions (
	test_id UUID NOT NULL PRIMARY KEY,
	"definition" json NOT NULL
);

CREATE TABLE IF NOT EXISTS runs (
	id UUID NOT NULL PRIMARY KEY,
	test_id UUID NOT NULL,
	run json NOT NULL
);
