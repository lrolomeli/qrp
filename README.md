# QRP — QR Redirect Platform

Backend en Go para redirección de códigos QR con registro de clics y dashboard de estadísticas.

## Stack

- **Backend:** Go 1.23+ (stdlib, sin frameworks)
- **DB:** SQLite vía modernc.org/sqlite (pure Go, sin CGO)
- **Infra:** Docker Compose + Caddy reverse proxy

## Endpoints

| Método | Ruta | Descripción |
|--------|------|-------------|
| `GET` | `/menu/{code}` | Redirige al destino del QR (302) |
| `GET` | `/api/dashboard` | Dashboard HTML con estadísticas |
| `GET` | `/api/health` | Health check |

## Configuración de QR codes

Edita `codes.json`:

```json
{
  "mesa1": "https://tusitio.com/menu/mesa1",
  "wifi":  "https://tusitio.com/wifi"
}
```

Luego reinicia el contenedor:

```bash
docker compose restart
```

## Despliegue

```bash
# Construir y levantar
docker compose up -d --build

# El servicio se conecta a proxy-net (Caddy)
# Caddy debe tener una ruta como:
#
#   lrlpapp.store {
#       tls internal
#       reverse_proxy qrp-backend:8080
#   }
```

## Variables de entorno

| Variable | Default | Descripción |
|----------|---------|-------------|
| `PORT` | `8080` | Puerto del servidor HTTP |
| `DB_PATH` | `/data/qrp.db` | Ruta a la base SQLite |
| `CODES_PATH` | `/data/codes.json` | Ruta al archivo de QR codes |
