FROM golang:1.22
SHELL ["/bin/bash", "-c"]
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
    ca-certificates \
    bash \
    dos2unix \
    "
RUN apt-get update && apt-get install -y ${DEPS}

# download go modules
WORKDIR /app

RUN mkdir bin

COPY /src/ src
COPY /sdk/ sdk

ADD ./build.sh /app/build.sh
ADD ./config.yaml /app/config.yaml
ADD ./.env /app/.env

# install kalkancryptwr
RUN dos2unix /app/sdk/production/install_production.sh
#RUN chmod +x /app/sdk/production/install_production.sh
RUN source /app/sdk/production/install_production.sh
RUN cp /app/sdk/libkalkancryptwr-64.so.2.0.3 /usr/lib
RUN mv /usr/lib/libkalkancryptwr-64.so.2.0.3 /usr/lib/libkalkancryptwr-64.so

# execute build.sh file
RUN dos2unix /app/build.sh
RUN chmod +x /app/build.sh
RUN /app/build.sh

ENTRYPOINT [ "/bin/bash" ]