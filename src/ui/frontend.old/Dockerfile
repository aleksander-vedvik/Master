# syntax=docker/dockerfile:1

FROM golang:1.20 as builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY ./frontend/src/public ./public
COPY ./frontend/src/views ./views
COPY ./frontend/src ./frontend/src
#COPY ./frontend/src/main.go ./src/
#COPY ./frontend/src/node.go ./src/
#COPY ./frontend/src/const.go ./src/

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/docker-simplewebsite ./frontend/src/

#COPY ./frontend/src ./frontend/src
#"ui/frontend/src/storage"

#ENTRYPOINT ["/app/docker-simplewebsite"]

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/public /public
COPY --from=builder /app/views /views
COPY --from=builder /app/docker-simplewebsite /docker-simplewebsite
CMD ["/docker-simplewebsite"]