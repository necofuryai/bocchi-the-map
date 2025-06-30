# Build stage
FROM --platform=$BUILDPLATFORM golang:1.24-alpine@sha256:c694a4d291a13a9f9d94933395673494fc2cc9d4777b85df3a7e70b3492d3574 AS builder

# Build arguments for cross-compilation
ARG TARGETOS
ARG TARGETARCH

# Install git and ca-certificates for dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with proper cross-compilation
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -a -installsuffix cgo -o main ./cmd/api

# Final stage
FROM --platform=$TARGETPLATFORM alpine:3.19@sha256:c5b1261d6d3e43071626931fc004f70149baeba2c8ec672bd4f27761f8e1ad6b

# Install ca-certificates and wget for HTTPS requests and health checks
RUN apk --no-cache add ca-certificates wget

# Create non-root user
RUN addgroup -g 1001 -S bocchi && \
    adduser -S bocchi -u 1001

# Set working directory and ensure ownership
WORKDIR /home/bocchi
RUN chown bocchi:bocchi /home/bocchi

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Change ownership to non-root user
RUN chown bocchi:bocchi main

# Switch to non-root user
USER bocchi

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./main"]