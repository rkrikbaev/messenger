#!/bin/bash

rm -rf ~/.cache/go-build
rm -rf /app/bin

mkdir /app/bin

go mod init test.go
go mod tidy

go get github.com/joho/godotenv v1.5.1
go get github.com/google/uuid
go get github.com/lib/pq v1.10.0

cd /app/src

go build -o /app/bin/app app.go

ls -l /app/bin/app

echo "Build complete"

sleep 5

rm -rf /app/src
rm -rf /app/build.sh
rm -rf /app/go.mod
rm -rf /app/go.sum

echo "Clean up complete"

sleep 5
