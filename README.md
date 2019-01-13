# books backend

For migration we use https://github.com/golang-migrate/migrate/tree/master/cli

For up migration against database:
```
migrate -path migrations/ -database "postgres://DB_USER@localhost:5432/DB_NAME?sslmode=disable&password=DB_PASSWORD" up
```

##### Necessary ENV variables

* `DB_HOST` default `127.0.0.1`
* `DB_PORT` default  `5432`
* `DB_NAME` required
* `DB_USER` required
* `DB_PASSWOR` required
* `DB_SSL_MODE` default `disabled`
* `PORT` default `8080`


### Docker

**Build:**
``` shell
docker build -t <image name> .
```
**Run:**
``` shell
docker run --env-file <path to env file> <image name>
```

### Docker-compose
Additional [ENV variables](#necessary-env-variables) for docker-compose

* `DB_DATA` volume for _/var/lib/postgresql/data_ (required, example: `/home/user/data/`)
* `HOST_DB_PORT` postgres port for host machine (default `5433`)



**Run ready environment:**
```shell
docker-compose up --build
```
or postgres only
```shell
docker-compose up postgres
```
