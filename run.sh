#!/bin/bash

#apt-get update
#apt-get install libpcsclite-dev
apt-get install gcc
apt-get install libltdl7
#apt-get install build-essential
#apt-get install glibc-devel
#apt install golang-go

rm -rf ~/.cache/go-build

#sudo source /app/sdk/SDK\ 2.0/C/Linux/ca-certs/ca-certs_new/production2015/install_production.sh
#sudo source /app/sdk/SDK\ 2.0/C/Linux/ca-certs/ca-certs_new/test2015/install_test.sh
source /app/sdk/production/install_production.sh
source /app/sdk/test/install_test.sh

cp /app/sdk/libkalkancryptwr-64.so.2.0.3 /usr/lib
mv /usr/lib/libkalkancryptwr-64.so.2.0.3 /usr/lib/libkalkancryptwr-64.so
#sudo cp -f /app/sdk/SDK\ 2.0/C/Linux/C/libs/v2.0.3/libkalkancryptwr-64.so.2.0.3 /usr/lib/libkalkancryptwr-64.so


go get github.com/gokalkan/gokalkan
go get github.com/google/uuid

## here in /app
cd /app

rm go.mod
rm go.sum

go mod init test.go
go mod tidy

go get github.com/gokalkan/gokalkan
go get github.com/google/uuid

go run app.go
#go build test.go
#./test
