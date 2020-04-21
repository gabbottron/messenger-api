# This is based on my generic golang two-stage
# Docker build. Set variables in docker-compose.

### ------ BUILD STAGE -----------------------------
FROM golang:1.13 as build-stage
ARG PROJECT_NAME
MAINTAINER Geoffrey Abbott - gabbottron@gmail.com
WORKDIR /wrk/${PROJECT_NAME}-api
COPY . .
RUN CGO_ENABLED=0 go build -o ./bin/${PROJECT_NAME}-api -a -installsuffix cgo ./src/${PROJECT_NAME}-api/main.go

### ------- RUN STAGE ------------------------------
FROM alpine:latest as run-stage
ARG PROJECT_NAME
RUN apk --no-cache add ca-certificates
RUN apk --no-cache add postgresql-client
WORKDIR /app/
COPY --from=build-stage /wrk/${PROJECT_NAME}-api/bin/${PROJECT_NAME}-api .
