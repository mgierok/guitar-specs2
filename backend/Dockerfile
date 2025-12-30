FROM golang:1.25-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/seed ./cmd/seed

FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=builder /out/api /app/api
COPY --from=builder /out/seed /app/seed

EXPOSE 8080

ENTRYPOINT ["/app/api"]
