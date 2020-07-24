FROM golang:1.14.6-alpine3.12

# git
RUN apk add git

# chrome
RUN apk add chromium chromium-chromedriver

# go install
WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["jclockedio"]
