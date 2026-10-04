# Stage 1: Сборка бинарного файла
FROM golang:1.27.1-alpine AS builder

WORKDIR /app

# Устанавливаем системные зависимости, если понадобятся (например, git или gcc)
RUN apk add --no-cache git

# Копируем go.mod и go.sum для кэширования слоев зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем бинарник из cmd/main.go
# CGO_ENABLED=0 создает статически скомпонованный бинарный файл без внешних C-зависимостей
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main ./cmd/main.go

# Stage 2: Минимальный финальный образ
FROM alpine:latest

WORKDIR /app

# Добавляем сертификаты CA для работы HTTPS-запросов (если приложение будет делать внешние запросы)
RUN apk --no-cache add ca-certificates tzdata

# Копируем скомпилированный бинарный файл из этапа сборки
COPY --from=builder /app/main .

# Запускаем приложение
CMD ["./main"]