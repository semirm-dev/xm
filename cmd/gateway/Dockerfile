# build app
FROM golang:1.23.2-alpine3.20 AS base_build

WORKDIR /app

COPY go.* ./
RUN go mod download

# instead of copying everything, we should/could copy only releveant code (cmd/gateway, internal)
# for the simplicity keep it as it is
COPY . .
RUN go build -v -o gateway cmd/gateway/main.go

# create runtime
FROM alpine:3.20

WORKDIR /app

RUN apk add ca-certificates

COPY --from=base_build /app/gateway .

EXPOSE 8080

ENTRYPOINT ["/app/gateway"]
