FROM golang:1.23-alpine AS builder
WORKDIR /app

# Cache and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build the binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o used2book ./cmd/used2book

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/used2book .

EXPOSE 80
CMD ["./used2book"]
