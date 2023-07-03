FROM golang:1.20.3 as builder

RUN apt-get update && apt-get install -y upx

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-w -s' -o=/tmp/bin/main .

RUN upx -5 /tmp/bin/main

# Use alpine image for smaller images
FROM golang:1.20.3-alpine3.17

COPY --from=builder /tmp/bin/main /usr/local/bin/main

EXPOSE 8080

CMD ["/usr/local/bin/main"]