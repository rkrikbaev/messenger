CREATE TABLE messenger."DOCUMENTS" (
    "ID" SERIAL PRIMARY KEY,
    "document" TEXT,
    "created" TEXT NOT NULL UNIQUE,
    "sent" TEXT,
    "state" TEXT
);