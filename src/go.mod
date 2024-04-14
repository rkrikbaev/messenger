module main

go 1.20

require (
	github.com/gokalkan/gokalkan v1.3.1
	github.com/google/uuid v1.6.0 // indirect
	github.com/joho/godotenv v1.5.1
	github.com/mattn/go-sqlite3 v1.14.22
	utils v0.0.0
	dbsql v0.0.0
)

replace utils => ./routes/utils
replace dbsql => ./routes/dbsql
