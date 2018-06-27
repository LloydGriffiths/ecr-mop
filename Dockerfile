FROM frolvlad/alpine-glibc

RUN apk update \
    && apk add ca-certificates \
    && rm -rf /var/cache/apk/*

COPY build/ecr-mop /usr/local/bin/ecr-mop
