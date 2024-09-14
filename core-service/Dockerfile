# Этап сборки
FROM golang:1.23-alpine AS builder
RUN apk --no-cache add git
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Компилируем приложение
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/core/main.go


FROM alpine:latest
WORKDIR /app

# Копируем бинарный файл из этапа сборки
COPY --from=builder /app/main /app/config/*.yml ./
# RUN apk --no-cache add ca-certificates

EXPOSE 50151

ENTRYPOINT [ "./main"]
CMD ["-config", "dev.yml"]