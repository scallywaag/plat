FROM golang:1.22-alpine AS builder
WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o plat .


FROM alpine:3.19

COPY --from=builder /app/plat /plat
EXPOSE 8080

RUN adduser -D app
USER app

CMD ["/plat"]

