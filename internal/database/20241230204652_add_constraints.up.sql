-- ##### Add foreign key constraint to ballot_entry ######

PRAGMA foreign_keys=OFF;

--bun:split

CREATE TABLE "ballot_entry_prime" (
  "id" INTEGER NOT NULL, 
  "election_id" INTEGER NOT NULL, 
  "name" VARCHAR, 
  "is_ron" BOOLEAN, 
  PRIMARY KEY ("id"),
  FOREIGN KEY ("election_id") REFERENCES "elections" ("id") ON DELETE CASCADE
);

--bun:split

INSERT INTO "ballot_entry_prime" ("id", "election_id", "name", "is_ron") 
SELECT 
  "id", 
  "election_id", 
  "name", 
  "is_ron"
FROM 
  "ballot_entry";

--bun:split

DROP TABLE "ballot_entry";

--bun:split

ALTER TABLE "ballot_entry_prime" RENAME TO "ballot_entry";

--bun:split

PRAGMA foreign_keys_check;

--bun:split



-- ##### Add foreign key constraint to candidates ######

CREATE TABLE "candidates_prime" (
  "user_id" VARCHAR, 
  "election_id" INTEGER, 
  PRIMARY KEY ("user_id", "election_id"),
  FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE,
  FOREIGN KEY ("election_id") REFERENCES "elections" ("id") ON DELETE CASCADE
);

--bun:split

INSERT INTO "candidates_prime" ("user_id", "election_id") 
SELECT 
  "user_id", 
  "election_id"
FROM 
  "candidates";

--bun:split

DROP TABLE "candidates";

--bun:split

ALTER TABLE "candidates_prime" RENAME TO "candidates";

--bun:split

PRAGMA foreign_keys_check;

--bun:split



-- ##### Add foreign key constraint to votes ######

CREATE TABLE "votes_prime" (
  "id" INTEGER NOT NULL, 
  "poll_id" INTEGER NOT NULL, 
  "user_id" VARCHAR NOT NULL, 
  "choices" VARCHAR, 
  PRIMARY KEY ("id"),
  FOREIGN KEY ("poll_id") REFERENCES "polls" ("id") ON DELETE CASCADE,
  FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE RESTRICT
);

--bun:split

INSERT INTO "votes_prime" ("id", "poll_id", "user_id", "choices") 
SELECT 
  "id", 
  "poll_id", 
  "user_id", 
  "choices"
FROM 
  "votes";

--bun:split

DROP TABLE "votes";

--bun:split

ALTER TABLE "votes_prime" RENAME TO "votes";

--bun:split

PRAGMA foreign_keys_check;

--bun:split

PRAGMA foreign_keys=ON;
