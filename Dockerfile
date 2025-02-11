# Etapa de construcción
FROM golang:alpine3.21 AS builder

# Establecer el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copiar los archivos go.mod y go.sum y descargar las dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar el código fuente al contenedor
COPY . .

# Construir la aplicación
RUN go build -o crud ./cmd/app

# Etapa final
FROM alpine:latest

EXPOSE 8080

ENV USER=docker
ENV UID=1100
ENV GROUP=docker
ENV GID=1200

# Establecer el directorio de trabajo dentro del contenedor
WORKDIR /root/

RUN addgroup -S --gid "$GID" "$GROUP" && \
    adduser \
    --disabled-password \
    --ingroup "$GROUP" \
    --no-create-home \
    --uid "$UID" \
    "$USER"

# Copiar el binario construido desde la etapa de construcción
COPY --from=builder /app/crud .

# Cambiar el propietario del binario al usuario creado
RUN chown "$USER":"$GROUP" ./crud


USER $USER
# Exponer el puerto en el que la aplicación escucha
EXPOSE 8080

# Comando para ejecutar la aplicación
CMD ["./crud"]