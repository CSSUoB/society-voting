CREATE TABLE "users" (
  "id" VARCHAR NOT NULL, 
  "name" VARCHAR, 
  "password_hash" BLOB, 
  PRIMARY KEY ("id")
)

--bun:split

CREATE TABLE "elections" (
  "id" INTEGER NOT NULL, 
  "role_name" VARCHAR, 
  "description" VARCHAR, 
  "is_active" BOOLEAN, 
  PRIMARY KEY ("id")
)

--bun:split

CREATE TABLE "candidates" (
  "user_id" VARCHAR, 
  "election_id" INTEGER, 
  PRIMARY KEY ("user_id", "election_id")
)

--bun:split

CREATE TABLE "ballot_entry" (
  "id" INTEGER NOT NULL, 
  "election_id" INTEGER, 
  "name" VARCHAR, 
  "is_ron" BOOLEAN, 
  PRIMARY KEY ("id")
)

--bun:split

CREATE TABLE "votes" (
  "id" INTEGER NOT NULL, 
  "election_id" INTEGER, 
  "user_id" VARCHAR, 
  "choices" VARCHAR, 
  PRIMARY KEY ("id")
)
