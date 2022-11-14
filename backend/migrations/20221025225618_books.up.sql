CREATE TABLE IF NOT EXISTS "books" (
  "bookId" integer not null primary key autoincrement, 
  "author" varchar(255), 
  "title" varchar(255), 
  "cover" varchar(255), 
  "filepath" integer, 
  "favorite" boolean default '0', 
  "segmented_file" varchar(255), 
  "has_read" boolean default '0', 
  unique ("bookId"),
  unique ("author", "title")
)
