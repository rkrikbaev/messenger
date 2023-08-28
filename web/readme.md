# docker
## Build binary file
docker run --rm -it -v d:\isun_log\app:/app --name=work_container go-env:1.20 go build main.go

docker run --name web-app -d --restart=always -v d:\isun_log\web:/app -p 5000:5000 web-app:latest