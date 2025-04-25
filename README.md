# Condominio App

**Condominio App** es una aplicación para la gestión de condominios, que incluye funcionalidades como cuotas de mantenimiento, servicios comunes, pagos, y más.

## Tabla de Contenidos
- [Tecnologías Usadas](#tecnologías-usadas)
- [Requisitos](#requisitos)
- [Instalación](#instalación)
- [Variables de Entorno](#variables-de-entorno)
- [Comandos Disponibles](#comandos-disponibles)
- [Estructura del Proyecto](#estructura-del-proyecto)
  - [Cuotas de Mantenimiento](#cuotas-de-mantenimiento)
  - [Servicios Comunes](#servicios-comunes)
  - [Pagos](#pagos)
- [Contribuciones](#contribuciones)
- [Licencia](#licencia)

---

## Tecnologías Usadas

Este proyecto utiliza las siguientes tecnologías:

- Golang
- PostgreSQL
- Docker

---

## Requisitos

- Tener instalado Docker y Docker Compose.
- Tener instalado Go (versión 1.18 o superior).
- Acceso a una terminal o línea de comandos.

---

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

---

## Uso

### Paso 1: Ejecutar el servicio localmente
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
```

---
## Variables de Entorno
Asegúrate de configurar las siguientes variables de entorno antes de ejecutar la aplicación:

```sh
DB_HOST: Host de la base de datos.
DB_USER: Usuario de la base de datos.
DB_PASSWORD: Contraseña de la base de datos.
DB_NAME: Nombre de la base de datos.
```


## Comandos Disponibles
```sh
go run main.go: #Ejecuta la aplicación localmente.

go mod tidy: #Limpia y organiza las dependencias.
docker build -t condominio-app .: #Construye la imagen Docker.
docker run -p 8080:8080 condominio-app: #Ejecuta la aplicación en un contenedor Docker.

```

## Contribuciones

¡Las contribuciones son bienvenidas! Si deseas contribuir, por favor sigue estos pasos:

1. Haz un fork del repositorio.
2. Crea una rama para tu funcionalidad (`git checkout -b feature/nueva-funcionalidad`).
3. Realiza tus cambios y haz un commit (`git commit -m 'Agrega nueva funcionalidad'`).
4. Haz un push a tu rama (`git push origin feature/nueva-funcionalidad`).
5. Abre un Pull Request.

---