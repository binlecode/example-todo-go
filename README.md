# simple todos REST service

## project overview

What's in this app:

- GORM as an ORM to interact with database
- GORM database automigration to create tables and seed data
- sqlite3 file db
- Request router using gorilla/mux
- Logrus for logging

## run in local

```sh
go run .
```

Health check endpoint:

```sh
curl http://localhost:9000/health
```

## project bootstrap

```sh
go get -u github.com/gorilla/mux
go get -u github.com/sirupsen/logrus
go get -u github.com/joho/godotenv 

go get -u github.com/mattn/go-sqlite3
go get -u gorm.io/gorm
go get -u gorm.io/driver/sqlite
```

