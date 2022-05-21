# syntax=docker/dockerfile:1

FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./

RUN cd ./cmd/main_service && go build -o /employee_manage

RUN apk update && apk add bash

RUN chmod 775 wait-for-it.sh && chmod 775 copy_migrate_and_validator.sh

RUN ./copy_migrate_and_validator.sh

EXPOSE 8088