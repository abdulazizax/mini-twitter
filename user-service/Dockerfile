# Stage 1: Build stage
FROM golang:1.23.1 AS builder

WORKDIR /app

# Go proxy ni o'rnatish
ENV GOPROXY=direct

# Loyihaning qolgan qismini nusxalash
COPY . .

# Go modulyatsiya fayllarini yuklash
RUN go mod download

# Agar kerak bo'lsa, .env faylini nusxalash
COPY .env .

# Dastur build qilish
RUN CGO_ENABLED=0 GOOS=linux go build -o myapp ./cmd

# Stage 2: Final stage
FROM alpine:latest

WORKDIR /app

# Birinchi bosqichdan (builder) compiled binary faylini nusxalash
COPY --from=builder /app/myapp .

# Agar kerak bo'lsa, .env faylini nusxalash
COPY --from=builder /app/.env .

# Port 8081 ni expose qilish
EXPOSE 8081

# Dastur executable faylini ishga tushirish
CMD ["./myapp"]
