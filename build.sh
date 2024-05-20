#!/bin/bash

rm -rf ~/.cache/go-build
rm -rf /app/bin

mkdir /app/bin

go mod init app.go
go mod tidy

go get github.com/joho/godotenv v1.5.1
go get github.com/google/uuid
go get github.com/lib/pq v1.10.0 

#  Install KalkanCrypt
go get github.com/gokalkan/gokalkan
source /app/sdk/production/install_production.sh
cp /app/sdk/libkalkancryptwr-64.so.2.0.3 /usr/lib
mv /usr/lib/libkalkancryptwr-64.so.2.0.3 /usr/lib/libkalkancryptwr-64.so

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
