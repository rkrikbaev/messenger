module app

go 1.20

require (
	dbsql v0.0.0
	github.com/fatih/structs v1.1.0
	github.com/gokalkan/gokalkan v1.3.1
	github.com/google/uuid v1.6.0 // indirect
	github.com/joho/godotenv v1.5.1
	github.com/mattn/go-sqlite3 v1.14.22
	gopkg.in/yaml.v2 v2.4.0
	postgresdb v0.0.0
	utils v0.0.0
)

replace utils => ./routes/utils

replace dbsql => ./routes/dbsql

replace postgresdb => ./routes/postgresdb
