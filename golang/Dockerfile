FROM golang:1.27.1-alpine

WORKDIR /app

RUN apk add --no-cache git ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# При старте контейнер берет свежий код из монтируемой папки и запускает его
CMD ["go", "run", "./cmd/main.go"]