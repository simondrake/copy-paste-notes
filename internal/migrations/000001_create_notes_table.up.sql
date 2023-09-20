CREATE TABLE IF NOT EXISTS "notes" (
  "id" INTEGER NOT NULL PRIMARY KEY,
  "create_timestamp" TEXT,
  "title" TEXT NOT NULL UNIQUE,
  "description" TEXT NOT NULL
  );
