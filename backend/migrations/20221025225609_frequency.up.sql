CREATE TABLE IF NOT EXISTS "frequency" (
  "book" integer, 
  "word" varchar(255), 
  "count" integer, 
  primary key ("book", "word")
)
