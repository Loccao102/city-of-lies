# Stage 1: Build
FROM golang:alpine AS builder

WORKDIR /app

COPY backend/go.mod backend/go.sum* ./
RUN go mod download

COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/sim ./cmd/sim

# Stage 2: Runtime
FROM alpine:3.20

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /app/bin/api /app/bin/api
COPY --from=builder /app/bin/sim /app/bin/sim
COPY data/ /app/data/

EXPOSE 8080

ENTRYPOINT ["/app/bin/api"]
