---
title: CRUD + JWT
keywords: [auth, jwt, gorm, fiber]
description: Simple JWT authentication.
---

## Archivos Principales

- `main.go`: Punto de entrada de la aplicación. Configura las rutas y arranca el servidor.
- `cmd/api/v0/info/handlers.go`: Contiene el handler para la ruta de información del servicio.
- `cmd/api/v0/info/info.go`: Define las estructuras de datos y funciones relacionadas con la información del servicio.

## Instalación

1. Clona el repositorio:
    ```sh
    git clone <URL_DEL_REPOSITORIO>
    cd <NOMBRE_DEL_REPOSITORIO>
    ```

2. Instala las dependencias:
    ```sh
    go mod tidy
    ```

## Uso

Para ejecutar el servicio, utiliza el siguiente comando:

```sh
go run main.go


docker run -d --name myposgres -e POSTGRES_USER=user1 -e POSTGRES_PASSWORD=password postgres 

psql -U user1 --password

CREATE DATABASE crud;


Paso 2: Verificar archivos go.mod y go.sum
Asegúrate de que los archivos go.mod y go.sum estén actualizados. Puedes hacerlo ejecutando los siguientes comandos en tu entorno de desarrollo local:

go mod tidy
go mod verify

verify
Paso 3: Construir y ejecutar la aplicación
Para construir y ejecutar la aplicación utilizando Docker, sigue estos pasos:

Construir la imagen Docker:
docker build -t crud-app .


Ejecutar el contenedor Docker:

docker run -p 8080:8080 crud-app