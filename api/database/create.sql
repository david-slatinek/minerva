CREATE TABLE "songs" (
     "id" TEXT NOT NULL UNIQUE,
     "title" TEXT NOT NULL,
     "duration" TEXT NOT NULL,
     "release" TEXT NOT NULL,
     "author" TEXT NOT NULL,
     PRIMARY KEY("id")
);
