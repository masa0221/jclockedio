# syntax=docker/dockerfile:1
FROM golang:1.17.8-alpine3.15
RUN apk add git chromium chromium-chromedriver tzdata

WORKDIR /app

COPY . ./
RUN go mod download

RUN go build -o /jclockedio

ENV CGO_ENABLED=0
ENV TZ=Asia/Tokyo

ENTRYPOINT [ "/jclockedio" ]

