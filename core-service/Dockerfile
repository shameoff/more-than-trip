# Этап сборки
FROM golang:1.23-alpine AS builder
RUN apk --no-cache add git
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Компилируем приложение
RUN go build -o /app/main ./cmd/core/main.go

FROM busybox:latest
WORKDIR /app

# Копируем бинарный файл из этапа сборки
COPY --from=builder /app/main /app/config/ ./
# RUN apk --no-cache add ca-certificates

EXPOSE 50151

CMD ["./main", "-config", "dev.yml"]