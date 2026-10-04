FROM python:3.11-alpine

# Установка системных зависимостей для сборки C-расширений (например, psycopg2)
RUN apk add --no-cache gcc musl-dev libpq-dev

WORKDIR /app

# Копируем и устанавливаем зависимости
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Копируем исходный код
COPY . .

CMD ["python", "main.py"]