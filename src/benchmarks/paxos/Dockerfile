# syntax=docker/dockerfile:1

FROM golang:alpine

WORKDIR /app

COPY ../go.mod .
COPY ../go.sum .
RUN go mod download

COPY . .

RUN go build -o /server ./cmd/server/main.go

# Internal port
EXPOSE 8080

CMD [ "/server" ]