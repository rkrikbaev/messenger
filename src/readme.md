# Parse csv to DB
docker run -it --restart=always -v d:\isun_log\app:/app --name=work_container go-env:1.20

## Build binary file
docker run --rm -it -v d:\isun_log\app:/app --name=work_container go-env:1.20 go build main.go