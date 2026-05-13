# Usa una imagen base de Go
FROM golang:1.22

# Instala dependencias externas necesarias para tu bot
RUN apt-get update && apt-get install -y \
    ffmpeg \
    imagemagick \
    exiv2 \
    gifsicle \
    libarchive-tools \
    && rm -rf /var/lib/apt/lists/*

# Crea directorio de trabajo
WORKDIR /app

# Copia los archivos del repositorio
COPY . .

# Descarga dependencias de Go
RUN go mod tidy

# Compila el bot
RUN go build -o bot ./cmd/moe-sticker-bot

# Variables de entorno necesarias
ENV BOT_TOKEN=${BOT_TOKEN}
ENV DB_ADDR=${DB_ADDR}
ENV DB_USER=${DB_USER}
ENV DB_PASS=${DB_PASS}
ENV DB_NAME=${DB_NAME}

# Comando de inicio
CMD ["./bot"]