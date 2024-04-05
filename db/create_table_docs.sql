CREATE TABLE "DOCUMENTS" (
	"ID"	INTEGER,
	"document"	TEXT NOT NULL,
	"response"	TEXT,
	"datetime"	TEXT NOT NULL UNIQUE,
	"sentMessageDate"	TEXT,
	"createdAtDate"	TEXT,
	"state"	TEXT,
	PRIMARY KEY("ID" AUTOINCREMENT)
);