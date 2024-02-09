FROM golang:1.21-bullseye as build-env

ARG APP_VERSION

WORKDIR /usr/src/app
ADD ./ /usr/src/app

RUN go build -ldflags "-X github.com/nixys/nxs-data-anonymizer/ctx.version=$APP_VERSION" -o /nxs-data-anonymizer

FROM alpine

RUN apk update --no-cache && apk add --no-cache tar postgresql-client mysql-client libc6-compat s3cmd

COPY --from=build-env /nxs-data-anonymizer /

