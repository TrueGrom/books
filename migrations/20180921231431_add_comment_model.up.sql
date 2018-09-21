CREATE TABLE comment_models
(
  "id" SERIAL,
  "created_at" TIMESTAMP WITH TIME ZONE,
  "updated_at" TIMESTAMP WITH TIME ZONE,
  "text" TEXT,
  "user_id" INTEGER REFERENCES user_models (id),
  "book_id" INTEGER REFERENCES books (id),
  PRIMARY KEY ("id")
)
