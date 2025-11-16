FROM golang:1.23-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o plat .


FROM gcr.io/distroless/static:nonroot

COPY --from=builder /app/plat /plat

USER nonroot

EXPOSE 8080

ENTRYPOINT ["/plat"]

