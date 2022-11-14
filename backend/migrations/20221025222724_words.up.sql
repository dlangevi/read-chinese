CREATE TABLE IF NOT EXISTS "words" (
  "word" varchar(255), 
  "has_flash_card" boolean, 
  "has_sentence" boolean, 
  "interval" integer, 
  "created_at" datetime not null default CURRENT_TIMESTAMP, 
  "updated_at" datetime not null default CURRENT_TIMESTAMP, 
  primary key ("word")
)

