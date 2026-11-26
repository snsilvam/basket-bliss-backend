# Basket Bliss — Backend (Go + Firestore)

API REST para los flujos v1 de la arquitectura:

| # | Flujo | Endpoint |
|---|---|---|
| 1 | Registrar el merchant (onboarding) | `POST /api/merchants` · `GET /api/merchants` |
| 2 | Registrar ancheta (y sus productos) | `POST /api/merchants/{merchantId}/products` · `POST /api/merchants/{merchantId}/anchetas` |
| 3 | Registrar customer | `POST /api/merchants/{merchantId}/customers` |
| 4 | Ver anchetas en venta | `GET /api/merchants/{merchantId}/anchetas` (incluye `available`) |
| 5 | Comprar ancheta | `POST /api/merchants/{merchantId}/customers/{customerId}/orders` |

También: `GET .../products`, `GET .../customers`, `GET .../customers/{customerId}/orders` y `GET /api/health`.

Sin pasarela de pagos: la orden queda en `pending_payment` y el stock ya descontado.

## Estructura

Sigue [Organizing a Go module](https://go.dev/doc/modules/layout) para un servicio:

```
cmd/api/main.go          punto de entrada: config, Firestore, dependencias y HTTP
internal/
  merchants/             handler.go · service.go · repository.go · model.go
  products/
  anchetas/
  customers/
  orders/                la compra es una transacción de Firestore
  platform/
    config/              variables de entorno
    database/            cliente Firestore + rutas de colecciones
    apperr/              errores de negocio (400/404/409)
    web/                 JSON, errores HTTP, log de peticiones
```

Modelo en Firestore:

```
merchants/{merchantId}
  products/{productId}          { name, stock }
  anchetas/{anchetaId}          { name, description, price, items[{productId, productName, quantity}] }
  customers/{customerId}        { name, email }
    orders/{orderId}            { anchetaId, anchetaName, quantity, unitPrice, total, status }
```

`available` de una ancheta = mínimo de `stock / cantidad` entre sus productos. Comprar descuenta
`cantidad × unidades` del stock de cada producto dentro de una transacción, así que dos compras
simultáneas no pueden vender el mismo stock.

## Ejecutar

Requiere Go 1.26 (si tienes uno anterior, `go` descarga la toolchain solo).

Contra Firestore en GCP (usa Application Default Credentials):

```powershell
gcloud auth application-default login
$env:GCP_PROJECT_ID = "basket-bliss-515621"
$env:FIRESTORE_DATABASE_ID = "basket-bliss-database"
go run ./cmd/api
```

Contra el emulador local:

```powershell
gcloud emulators firestore start --host-port=localhost:8081   # en otra terminal
$env:FIRESTORE_EMULATOR_HOST = "localhost:8081"
$env:GCP_PROJECT_ID = "demo-basket-bliss"
go run ./cmd/api
```

| Variable | Default | |
|---|---|---|
| `PORT` | `8080` | |
| `HOST` | (todas las interfaces) | `127.0.0.1` en local evita el aviso del firewall de Windows |
| `GCP_PROJECT_ID` | se detecta de las credenciales | |
| `FIRESTORE_DATABASE_ID` | `(default)` | |
| `FIRESTORE_EMULATOR_HOST` | — | si está definida, se usa el emulador |

```powershell
go test ./...
```
