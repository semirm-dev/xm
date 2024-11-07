# build app
FROM golang:1.23.2-alpine3.20 AS base_build

WORKDIR /app

COPY go.* ./
RUN go mod download

# instead of copying everything, we should/could copy only releveant code (cmd/companies, internal)
# for the simplicity keep it as it is
COPY . .
RUN go build -v -o companies cmd/companies/main.go

# create runtime
FROM alpine:3.20

WORKDIR /app

RUN apk add ca-certificates

COPY --from=base_build /app/companies .

EXPOSE 8080

ENTRYPOINT ["/app/companies"]
