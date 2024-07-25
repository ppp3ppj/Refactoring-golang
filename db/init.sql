CREATE TABLE "Person" (
  "key" TEXT PRIMARY KEY,
  "name" TEXT NOT NULL,
  "description" TEXT NOT NULL,
  "image" TEXT NOT NULL,
  "traits" JSONB NOT NULL DEFAULT '[]',
  "tags" TEXT[] NOT NULL DEFAULT '{}'
);
