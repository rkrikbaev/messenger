#!/bin/bash

if [ -e /usr/local/share/ca-certificates/extra/ ]; 
then 
	echo "Folder already exists"
else
	mkdir /usr/local/share/ca-certificates/extra
fi

cp -a /app/sdk/production/production/*.pem /etc/ssl/certs/

cd /app/sdk/production/production
for f in *.pem; do 
    mv -- "$f" "${f%.pem}.crt"
done
cd ..

cp -a /app/sdk/production/production/*.crt /usr/local/share/ca-certificates/extra/

source update-ca-certificates

cd /app/sdk/production/production
for f in *.crt; do 
    mv -- "$f" "${f%.crt}.pem"
done
cd ..
