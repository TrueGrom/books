CREATE OR REPLACE FUNCTION pseudo_encrypt(VALUE INT) RETURNS INT AS
$$
DECLARE
  l1 INT;
  l2 INT;
  r1 INT;
  r2 INT;
  i  INT := 0;
BEGIN
  l1 := (VALUE >> 12) & (4096 - 1);
  r1 := VALUE & (4096 - 1);
  WHILE i < 3
    LOOP
      l2 := r1;
      r2 := l1 # ((((1366 * r1 + 150889) % 714025) / 714025.0) * (4096 - 1))::INT;
      l1 := l2;
      r1 := r2;
      i := i + 1;
    END LOOP;
  RETURN ((l1 << 12) + r1);
END;
$$ LANGUAGE plpgsql
   STRICT
   IMMUTABLE;

CREATE SEQUENCE for_book_id START WITH 777777 INCREMENT BY 7;

CREATE TABLE IF NOT EXISTS books (
  id          INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  url_id      BIGINT DEFAULT pseudo_encrypt(nextval('for_book_id')::INTEGER),
  title       TEXT NOT NULL,
  publisher   TEXT,
  annotation  TEXT,
  isbn        TEXT,
  year        INTEGER,
  pages       INTEGER,
  circulation INTEGER,
  weight      INTEGER,
  cover       TEXT,
  size        TEXT,
  image       TEXT,
  rating      TEXT,
  UNIQUE (url_id),
  UNIQUE (image)
);

CREATE TABLE IF NOT EXISTS authors (
  id   INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  name TEXT UNIQUE
);

CREATE TABLE IF NOT EXISTS illustrators (
  id   INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  name TEXT UNIQUE
);

CREATE TABLE IF NOT EXISTS translators (
  id   INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  name TEXT UNIQUE
);

CREATE TABLE IF NOT EXISTS editors (
  id   INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  name TEXT UNIQUE
);

CREATE TABLE IF NOT EXISTS books_authors (
  author_id INTEGER REFERENCES authors (id),
  book_id   INTEGER REFERENCES books (id),
  UNIQUE (author_id, book_id)
);

CREATE TABLE IF NOT EXISTS books_illustrators (
  illustrator_id INTEGER REFERENCES illustrators (id),
  book_id        INTEGER REFERENCES books (id),
  UNIQUE (illustrator_id, book_id)
);

CREATE TABLE IF NOT EXISTS books_translators (
  translator_id INTEGER REFERENCES translators (id),
  book_id       INTEGER REFERENCES books (id),
  UNIQUE (translator_id, book_id)
);

CREATE TABLE IF NOT EXISTS books_editors (
  editor_id INTEGER REFERENCES editors (id),
  book_id   INTEGER REFERENCES books (id),
  UNIQUE (editor_id, book_id)
);

CREATE OR REPLACE FUNCTION full_text_search(VALUE TEXT, LIM INT)
  RETURNS TABLE (
    url_id       BIGINT,
    title        TEXT,
    publisher    TEXT,
    annotation   TEXT,
    isbn         TEXT,
    year         INTEGER,
    pages        INTEGER,
    circulation  INTEGER,
    weight       INTEGER,
    cover        TEXT,
    size         TEXT,
    rating       TEXT,
    image        TEXT
  ) AS
$$
BEGIN
  RETURN QUERY
    SELECT
      books.url_id,
      books.title,
      books.publisher,
      books.annotation,
      books.isbn,
      books.year,
      books.pages,
      books.circulation,
      books.weight,
      books.cover,
      books.size,
      books.rating,
      books.image
    FROM books
    WHERE to_tsvector('russian', books.title) @@ plainto_tsquery('russian', VALUE)
    LIMIT LIM;
END;
$$
  LANGUAGE plpgsql;

CREATE INDEX ON books (title);
CREATE INDEX ON books (isbn);
CREATE INDEX ON books USING gin (to_tsvector('russian', title));
CREATE INDEX ON books_authors (author_id);
CREATE INDEX ON books_authors (book_id);
CREATE INDEX ON books_editors (editor_id);
CREATE INDEX ON books_editors (book_id);
CREATE INDEX ON books_translators (translator_id);
CREATE INDEX ON books_translators (book_id);
CREATE INDEX ON books_illustrators (illustrator_id);
CREATE INDEX ON books_illustrators (book_id);
