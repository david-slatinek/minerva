CREATE TABLE "measurements"
(
    "id"     SERIAL PRIMARY KEY,
    "date"   TEXT NOT NULL,
    "time"   TEXT NOT NULL,
    "status" TEXT NOT NULL,
    "mode"   INT NOT NULL
);
