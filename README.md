# QRP — QR Redirect Platform

Backend en Go para redirección de códigos QR con registro de clics y dashboard de estadísticas.

## Cómo funciona

1. Generas un QR code apuntando a `https://lrlpapp.store`
2. El backend registra cada clic (IP, user-agent, timestamp)
3. Redirige automáticamente a la URL configurada en `REDIRECT_URL`

Para cambiar el destino, solo editas la variable de entorno y reinicias.

## Endpoints

| Método | Ruta | Descripción |
|--------|------|-------------|
| `GET` | `/` | Registra clic + redirige a `REDIRECT_URL` (302) |
| `GET` | `/api/dashboard` | Dashboard HTML con estadísticas |
| `GET` | `/api/health` | Health check |

## Despliegue

```bash
# Crear .env con REDIRECT_URL
cp .env.example .env
# Editar .env con tu dominio destino

# Construir y levantar
docker compose up -d --build

# Recargar Caddy (infra-proxy)
docker compose -f ../infra-proxy/docker-compose.yml up -d --build caddy
```

## Variables de entorno

| Variable | Default | Descripción |
|----------|---------|-------------|
| `REDIRECT_URL` | `https://ejemplo.com` | Destino de la redirección |
| `PORT` | `8080` | Puerto del servidor HTTP |
| `DB_PATH` | `/data/qrp.db` | Ruta a la base SQLite |
