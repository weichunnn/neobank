# multistage to avoid copying all the packages used which is heavy, take only binary file
# build stage
FROM golang:1.18 AS builder

WORKDIR /app

COPY . .

RUN go build -o main main.go


# prod stage
FROM golang:1.18-alpine

WORKDIR /app

COPY --from=builder /app/main .
EXPOSE 8080

CMD ["/app/main"]