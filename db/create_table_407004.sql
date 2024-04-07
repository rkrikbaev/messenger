CREATE TABLE "407004" (
	"ID"	INTEGER,
	"temperature"	REAL,
	"density"	REAL,
	"volume"	REAL,
	"tankLevel"	REAL,
	"mass"	REAL,
	"datetime"	TEXT UNIQUE,
	"createdAtDate"	TEXT,
	PRIMARY KEY("ID" AUTOINCREMENT)
);