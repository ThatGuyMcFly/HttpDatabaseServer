#syntax=docker/dockerfile:1

FROM golang:1.26
WORKDIR /HttpDatabaseServer

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o httpserver github.com/ThatGuyMcFly/HttpDatabaseServer/cmd/server

CMD ["./httpserver"]

EXPOSE 8080