# File Sharing API with Gofile Integration

Esta es una API REST desarrollada en Go que permite a los usuarios cargar, descargar, listar, eliminar y compartir archivos utilizando la API de Gofile como backend de almacenamiento. La API está protegida con autenticación JWT.

## Requisitos previos

- **Go** instalado (versión 1.16 o superior).
- Tener configurada la API Key de Gofile (si utilizas cuentas de usuario).
- **Git** instalado para clonar el proyecto.

-  Iniciar proyecto GO
go mod init file-sharing-api

- Instalar dependencias
go mod tidy



## Instalación

### 1. Clonar el repositorio

Primero, clona este repositorio en tu máquina local:

```bash
git clone https://github.com/grs89/ARQ-SOFTWARE2024
cd file-sharing-api


### Subir un archivo
```bash
curl -X POST http://localhost:8080/upload \
  -H "Authorization: Bearer <your_jwt_token>" \
  -F "file=@/path/to/your/file.txt"

### Descarhar un archivo
curl -X GET http://localhost:8080/download/<fileId> \
  -H "Authorization: Bearer <your_jwt_token>

### Eliminar un archivo
curl -X DELETE http://localhost:8080/files/<fileId> \
  -H "Authorization: Bearer <your_jwt_token>"
