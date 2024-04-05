#!/bin/bash

if [ -e /usr/local/share/ca-certificates/extra/ ]; 
then 
	echo "Folder already exists"
else
	mkdir /usr/local/share/ca-certificates/extra
fi

cp -a /app/sdk/test/test/*.pem /etc/ssl/certs/

cd /app/sdk/test/test
for f in *.pem; do 
    mv -- "$f" "${f%.pem}.crt"
done
cd ..

cp -a /app/sdk/test/test/*.crt /usr/local/share/ca-certificates/extra/

source update-ca-certificates

cd /app/sdk/test/test
for f in *.crt; do 
    mv -- "$f" "${f%.crt}.pem"
done
cd ..
