CREATE TABLE "poll_types" (
  "id" INTEGER NOT NULL, 
  "name" VARCHAR NOT NULL, 
  PRIMARY KEY ("id")
);

--bun:split

CREATE TABLE "polls" (
  "id" INTEGER NOT NULL, 
  "poll_type_id" INTEGER NOT NULL, 
  "is_active" BOOLEAN NOT NULL, 
  "is_concluded" BOOLEAN NOT NULL, 
  PRIMARY KEY ("id"), 
  FOREIGN KEY ("poll_type_id") REFERENCES "poll_types" ("id")
);

--bun:split

CREATE TABLE "poll_outcomes" (
  "id" INTEGER NOT NULL, 
  "poll_id" INTEGER NOT NULL, 
  "date" TIMESTAMP NOT NULL DEFAULT current_timestamp, 
  "ballots" INTEGER NOT NULL, 
  "is_published" BOOLEAN NOT NULL, 
  PRIMARY KEY ("id"), 
  UNIQUE ("poll_id"), 
  FOREIGN KEY ("poll_id") REFERENCES "polls" ("id") ON DELETE CASCADE
);

--bun:split

CREATE TABLE "referendums" (
  "id" INTEGER NOT NULL, 
  "title" VARCHAR NOT NULL, 
  "question" VARCHAR NOT NULL, 
  "description" VARCHAR NOT NULL, 
  PRIMARY KEY ("id"), 
  FOREIGN KEY ("id") REFERENCES "polls" ("id") ON DELETE CASCADE
);

--bun:split

CREATE TABLE "referendum_outcomes" (
  "id" INTEGER NOT NULL, 
  "votes_for" INTEGER NOT NULL, 
  "votes_against" INTEGER NOT NULL, 
  "votes_abstain" INTEGER NOT NULL, 
  PRIMARY KEY ("id"), 
  FOREIGN KEY ("id") REFERENCES "poll_outcomes" ("id") ON DELETE CASCADE
);

--bun:split

ALTER TABLE "votes" RENAME COLUMN "election_id" TO "poll_id";


--bun:split

INSERT INTO "poll_types" ("id", "name") 
VALUES 
  ("0", "Orphaned"), 
  ("1", "Election"), 
  ("2", "Referendum");

--bun:split



-- ###### Migrate election data into polls ######

INSERT INTO "polls" (
  "id", "poll_type_id", "is_active", "is_concluded"
) 
SELECT 
  "id", 
  1, 
  "is_active", 
  "is_concluded"
FROM 
  "elections";

--bun:split



-- ##### Migrate election outcomes data into poll outcomes and update table ######

PRAGMA foreign_keys=OFF;

--bun:split

INSERT INTO "poll_outcomes" ("id", "poll_id", "date", "ballots", "is_published") 
SELECT 
  "id", 
  "election_id", 
  "date", 
  "ballots",
  "is_published" 
FROM 
  "election_outcomes";

--bun:split

CREATE TABLE "election_outcomes_prime" (
  "id" INTEGER NOT NULL, 
  "rounds" INTEGER NOT NULL, 
  PRIMARY KEY ("id"), 
  FOREIGN KEY ("id") REFERENCES "poll_outcomes" ("id") ON DELETE CASCADE
);

--bun:split

INSERT INTO "election_outcomes_prime" ("id", "rounds") 
SELECT 
  "id", 
  "rounds"
FROM 
  "election_outcomes";

--bun:split

DROP TABLE "election_outcomes";

--bun:split

ALTER TABLE "election_outcomes_prime" RENAME TO "election_outcomes";

--bun:split

PRAGMA foreign_keys_check;



-- ##### Update election foreign key reference and remove columns ######

CREATE TABLE "elections_prime" (
  "id" INTEGER NOT NULL, 
  "role_name" VARCHAR NOT NULl, 
  "description" VARCHAR NOT NULL, 
  PRIMARY KEY ("id"), 
  FOREIGN KEY ("id") REFERENCES polls ("id") ON DELETE CASCADE
);

--bun:split

INSERT INTO "elections_prime" ("id", "role_name", "description") 
SELECT 
  "id", 
  "role_name", 
  "description"
FROM 
  "elections";

--bun:split

DROP TABLE "elections";

--bun:split

ALTER TABLE "elections_prime" RENAME TO "elections";

--bun:split

PRAGMA foreign_keys_check;

--bun:split

PRAGMA foreign_keys=ON;
