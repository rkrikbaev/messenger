CREATE TABLE messenger."DOCUMENTS"(
    "ID" SERIAL PRIMARY KEY,
    "document" TEXT,
    "created_dt" TEXT NOT NULL UNIQUE,
    "sent_dt" TEXT,
    "event_dt"  TEXT,
    "state" TEXT
);