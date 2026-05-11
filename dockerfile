# imagen base ligera con Go
FROM golang:1.21-bullseye AS builder

# dependencias externas necesarias
RUN apt-get update && apt-get install -y \
    ffmpeg \
    imagemagick \
    exiv2 \
    gifsicle \
    && rm -rf /var/lib/apt/lists/*

# Crea directorio de trabajo
WORKDIR /app

# Copia los archivos del bot
COPY . .

# librerías GO + binario
RUN go mod tidy && go build -o bot ./cmd/moe-sticker-bot

# Imagen final 
FROM debian:bullseye-slim

# dependencias necesarias en runtime
RUN apt-get update && apt-get install -y \
    ffmpeg \
    imagemagick \
    exiv2 \
    gifsicle \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copia el binario compilado
COPY --from=builder /app/bot .

# Variables de entorno
ENV TELEGRAM_TOKEN=${TELEGRAM_TOKEN}
ENV DATABASE_URL=${DATABASE_URL}

# Comando final para ejecutar el bot
CMD ["./bot"]
