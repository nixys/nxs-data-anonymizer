FROM golang:1.19-bullseye as build-env

WORKDIR /usr/src/app
ADD ./ /usr/src/app

RUN go build -o /nxs-data-anonymizer

FROM alpine

RUN apk update --no-cache && apk add --no-cache tar postgresql-client mysql-client libc6-compat s3cmd

COPY --from=build-env /nxs-data-anonymizer /

