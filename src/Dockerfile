FROM ubuntu:20.04

# set timezone
ENV DEBIAN_FRONTEND=noninteractive
ENV TZ="Asia/Almaty"
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# install dependencies
ARG DEPS=" \
    git \
    gcc \
    libpcsclite-dev \
    zlib1g-dev \
    libltdl7 \
    curl \
    "
RUN apt-get update && apt-get install -y ${DEPS}

# install golang
ARG GO_VERSION=1.20

RUN curl -sL https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz -o /tmp/go${GO_VERSION}.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf /tmp/go${GO_VERSION}.linux-amd64.tar.gz && rm /tmp/go${GO_VERSION}.linux-amd64.tar.gz

# golang config
ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

# download go modules
WORKDIR /app

RUN go get github.com/joho/godotenv v1.5.1
RUN go get github.com/mattn/go-sqlite3 v1.14.17

ENTRYPOINT ["./run.sh"]
