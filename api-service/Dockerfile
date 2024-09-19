# Stage 1: Build stage
FROM golang:1.23.1 AS builder

WORKDIR /app

# Go proxy ni o'rnatish
ENV GOPROXY=direct

# Copy Go modules manifests first to leverage Docker caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -C ./cmd -a -installsuffix cgo -o /app/myapp .

# Stage 2: Final stage
FROM alpine:latest

WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/myapp .

# Copy the configuration files
COPY --from=builder /app/internal/pkg/casbin/model.conf ./internal/pkg/casbin/
COPY --from=builder /app/internal/pkg/casbin/policy.csv ./internal/pkg/casbin/

# Optionally copy the .env file
COPY .env .

# Expose port 8080
EXPOSE 8080

# Command to run the executable
CMD ["./myapp"]
