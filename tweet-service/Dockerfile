# Stage 1: Build stage
FROM golang:1.23.1 AS builder

WORKDIR /app

# Go proxy ni o'rnatish
ENV GOPROXY=direct

# Go modulyatsiya fayllarini dastlab nusxalash (Docker caching dan foydalanish uchun)
COPY go.mod go.sum ./
RUN go mod download

# Loyihaning qolgan qismini nusxalash
COPY . .

# Dastur build qilish
RUN CGO_ENABLED=0 GOOS=linux go build -C ./cmd -a -installsuffix cgo -o /app/myapp .

# Stage 2: Final stage
FROM alpine:latest

WORKDIR /app

# Birinchi bosqichdan (builder) compiled binary faylini nusxalash
COPY --from=builder /app/myapp .

# Agar kerak bo'lsa, .env faylini nusxalash
COPY .env .

# Port 8082 ni expose qilish
EXPOSE 8082

# Dastur executable faylini ishga tushirish
CMD ["./myapp"]
