CREATE OR REPLACE FUNCTION pseudo_encrypt(VALUE int)
  RETURNS INT AS $$
DECLARE
  l1 INT;
  l2 INT;
  r1 INT;
  r2 INT;
  i  INT :=0;
BEGIN
  l1 := (VALUE >> 12) & (4096 - 1);
  r1 := VALUE & (4096 - 1);
  WHILE i < 3 LOOP
    l2 := r1;
    r2 := l1 # ((((1366 * r1 + 150889) % 714025) / 714025.0) * (4096 - 1)) :: INT;
    l1 := l2;
    r1 := r2;
    i := i + 1;
  END LOOP;
  RETURN ((l1 << 12) + r1);
END;
$$
LANGUAGE plpgsql
STRICT
IMMUTABLE;

DROP SEQUENCE IF EXISTS for_book_id;
CREATE SEQUENCE for_book_id
  START WITH 777777
  INCREMENT BY 7;

CREATE TABLE books (
  id           INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  url_id       BIGINT  DEFAULT pseudo_encrypt(nextval('for_book_id') :: int),
  title        TEXT NOT NULL,
  publisher    TEXT,
  annotation   TEXT,
  isbn         TEXT,
  authors      TEXT ARRAY,
  illustrators TEXT ARRAY,
  translators  TEXT ARRAY,
  editors      TEXT ARRAY,
  year         INTEGER,
  pages        INTEGER,
  image_file   TEXT,
  UNIQUE (url_id),
  UNIQUE (image_file)
);

CREATE INDEX ON books (isbn);
CREATE INDEX ON books
USING GIN (to_tsvector('russian', title));

CREATE OR REPLACE FUNCTION full_text_search(VALUE TEXT, LIM INT)
  RETURNS TABLE(
    url_id       BIGINT,
    title        TEXT,
    publisher    TEXT,
    annotation   TEXT,
    isbn         TEXT,
    authors      TEXT ARRAY,
    illustrators TEXT ARRAY,
    translators  TEXT ARRAY,
    editors      TEXT ARRAY,
    year         INTEGER,
    pages        INTEGER,
    image_file   TEXT
  ) AS $$
BEGIN
  RETURN QUERY
  SELECT books.url_id,
         books.title,
         books.publisher,
         books.annotation,
         books.isbn,
         books.authors,
         books.illustrators,
         books.translators,
         books.editors,
         books.year,
         books.pages,
         books.image_file
  FROM books
  WHERE to_tsvector('russian', books.title) @@ plainto_tsquery('russian', VALUE)
  LIMIT LIM;
END;
$$
LANGUAGE plpgsql;

CREATE TABLE IF NOT EXISTS user_models
(
  id       INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  username TEXT,
  email    TEXT,
  bio      VARCHAR(1024),
  image    TEXT,
  password TEXT NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS uix_user_models_username
  ON user_models (username);

CREATE UNIQUE INDEX IF NOT EXISTS uix_user_models_email
  ON user_models (email);


CREATE TABLE books_user_models (
  user_model_id    int REFERENCES user_models (id) ON UPDATE CASCADE ON DELETE CASCADE,
  book_model_id int REFERENCES books (id) ON UPDATE CASCADE,
  CONSTRAINT books_user_models_pkey PRIMARY KEY (user_model_id, book_model_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS uix_books_user_models_user_id_book_id
  ON books_user_models (user_model_id, book_model_id);


CREATE INDEX books_user_models_index_user_id ON public.books_user_models (user_model_id);
CREATE INDEX books_user_models_index_book_id ON public.books_user_models (book_model_id);