FROM golang:alpine3.18 as builder
WORKDIR /golang-postgres-api
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-w -s" .
FROM busybox
WORKDIR /golang-postgres-api
COPY --from=builder /golang-postgres-api /usr/bin/
ENTRYPOINT ["golang-postgres-api"]