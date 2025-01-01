CREATE TABLE "measurements"
(
    "id"         SERIAL PRIMARY KEY,
    "date"       TEXT NOT NULL,
    "time"       TEXT NOT NULL,
    "status_api" TEXT NOT NULL,
    "status_db"  TEXT NOT NULL,
    "mode"       INT  NOT NULL
);
