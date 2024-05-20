FROM golang:1.20

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
    iputils-ping \
    postgresql-client \
    curl \
    "
RUN apt-get update && apt-get install -y ${DEPS}

# download go modules
WORKDIR /app

RUN mkdir bin

COPY /src/ src

COPY ./build.sh /app/build.sh

# execute build.sh file
RUN chmod +x /app/build.sh
RUN /app/build.sh

ENTRYPOINT [ "/app/bin/app" ]