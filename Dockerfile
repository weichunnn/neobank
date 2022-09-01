# multistage to avoid copying all the packages used which is heavy, take only binary file
# build stage
FROM golang:1.19.0-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz

# prod stage
FROM alpine:3.13
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY db/migration ./migration
COPY wait-for.sh .
COPY start.sh .

COPY app.env .

EXPOSE 8000
CMD ["/app/main"]
ENTRYPOINT [ "/app/start.sh" ]