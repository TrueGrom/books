# books-backend

For migration we use https://github.com/golang-migrate/migrate/tree/master/cli

For up migration against database:
migrate -path migrations/ -database "postgres://DB_USER@localhost:5432/DB_NAME?sslmode=disable&password=DB_PASSWORD" up

