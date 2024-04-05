CREATE TABLE "407009" (
	"ID"	INTEGER,
	"temperature"	REAL,
	"density"	REAL,
	"volume"	REAL,
	"massflowbegin"	INTEGER,
	"massflowend"	INTEGER,
	"mass"	INTEGER,
	"datetime"	TEXT UNIQUE,
	"createdAtDate"	TEXT,
	PRIMARY KEY("ID" AUTOINCREMENT)
);