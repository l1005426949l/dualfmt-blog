# Dockerfile
FROM --platform=${BUILDPLATFORM} golang:1.21-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download && \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -a -o /app/file-management-service ./cmd/server

FROM --platform=${BUILDPLATFORM} alpine:latest

WORKDIR /app

COPY --from=builder /app/file-management-service /app/file-management-service
COPY --from=builder /app/configs /app/configs

EXPOSE 8000 9000

CMD ["/app/file-management-service", "-conf", "/app/configs"]
