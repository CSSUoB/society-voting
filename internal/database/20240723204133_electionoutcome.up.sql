CREATE TABLE "election_outcomes" (
  "id" INTEGER NOT NULL, 
  "election_id" INTEGER NOT NULL, 
  "date" TIMESTAMP NOT NULL DEFAULT current_timestamp, 
  "ballots" INTEGER NOT NULL, 
  "rounds" INTEGER NOT NULL, 
  "is_published" BOOLEAN NOT NULL, 
  PRIMARY KEY ("id"), 
  UNIQUE ("election_id"), 
  FOREIGN KEY ("election_id") REFERENCES "elections" ("id") ON DELETE CASCADE
);

--bun:split

CREATE TABLE "election_outcome_results" (
  "id" INTEGER NOT NULL, 
  "name" VARCHAR NOT NULL, 
  "round" INTEGER NOT NULL, 
  "votes" INTEGER NOT NULL, 
  "is_rejected" BOOLEAN NOT NULL, 
  "is_elected" BOOLEAN NOT NULL, 
  "election_outcome_id" INTEGER NOT NULL, 
  PRIMARY KEY ("id"), 
  FOREIGN KEY ("election_outcome_id") REFERENCES "election_outcomes" ("id") ON DELETE CASCADE
);
