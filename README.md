# simple todos REST service

## project overview

What's in this app:

- GORM as an ORM to interact with database
- GORM database auto-migration to create tables and seed data
- sqlite3 file db
- Request router using gorilla/mux
- Logrus for logging

## run in local

```sh
go run .
```

```sh
# run a local docker postgres instance
docker run --name example-todo-go-postgres \
  -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres
  
# test connection with container psql cli
docker exec -it example-todo-go-postgres psql -h localhost -p 5432 \
  -U postgres -d postgres
```


Health check endpoint:

```sh
curl http://localhost:9000/health
```

CORS check endpoint:

```sh
curl -D - -H 'Origin: http://foo.com' http://localhost:9000/health
```


## JWT auth

env var `SECRET_KEY` is used to sign JWT token.

`<host>/auth/authorize` endpoint is used to authenticate by user credentials
and get JWT token.

```sh
# use curl to get JWT token
curl -X POST -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"password1"}' \
  http://localhost:9000/auth/authorize
```

`<host>/auth/refresh` endpoint is used to refresh JWT token.

```sh

```




## project bootstrap

```sh
go get -u github.com/gorilla/mux
go get -u github.com/gorilla/handlers
go get -u github.com/sirupsen/logrus
go get -u github.com/joho/godotenv
# jwt
go get -u github.com/golang-jwt/jwt/v5

go get -u github.com/mattn/go-sqlite3
# lock gorm version to 1.22.2 to avoid sqlite3 driver issue 
go get -u gorm.io/gorm@v1.22.2
go get -u gorm.io/driver/sqlite
go get -u gorm.io/driver/postgres
go get -u github.com/lib/pq
```

