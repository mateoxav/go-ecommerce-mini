# go-ecommerce-mini

Sistema de Gestión de E-Commerce desarrollado en Go para el Proyecto Final de Programación Orientada a Objetos. Esta versión integra la evolución del proyecto: primero fue una aplicación CLI con SQLite y ahora se convierte principalmente en un backend web REST JSON, conservando la CLI como evidencia del progreso del desarrollo.

## Datos generales

- **Repositorio:** https://github.com/mateoxav/go-ecommerce-mini
- **Materia:** Programación Orientada a Objetos
- **Proyecto:** Sistema de Gestión de E-Commerce
- **Fecha:** Junio de 2026
- **Lenguaje:** Go 1.26.1
- **Persistencia:** SQLite local mediante `database/sql` y `github.com/mattn/go-sqlite3`
- **Servidor web:** Biblioteca estándar `net/http`
- **Serialización:** JSON mediante `encoding/json`

## Objetivo del programa

Implementar un sistema de gestión de e-commerce que permita administrar productos, clientes, pedidos e inventario, incorporando servicios web REST, persistencia local, manejo de errores, estructuras de datos, principios de orientación a objetos en Go, pruebas de software y concurrencia mediante goroutines y canales.

## Funcionalidades principales

- Backend REST JSON con más de 8 servicios web.
- CRUD funcional de productos con eliminación lógica.
- Registro, búsqueda, listado y eliminación lógica de clientes.
- Creación de pedidos, agregado de ítems, cálculo de total y cambio de estado.
- Control de inventario con verificación de stock, alertas de stock bajo, reposición y reportes.
- Endpoint concurrente que ejecuta tareas paralelas usando goroutines y canales.
- CLI conservada como modo alternativo para evidenciar la evolución del proyecto.
- Base de datos SQLite inicializada automáticamente.
- Tests unitarios, pruebas HTTP con `httptest` y flujo de aceptación sobre la API.

## Requisitos

- Go 1.26.1 o una versión compatible con el proyecto.
- GCC/CGO habilitado para compilar `github.com/mattn/go-sqlite3`.
- SQLite no requiere servidor externo porque se usa como archivo local.

## Instalación

```bash
git clone https://github.com/mateoxav/go-ecommerce-mini.git
cd go-ecommerce-mini
go mod tidy
```

## Ejecución del backend web

El modo principal del proyecto final es el servidor web:

```bash
go run .
```

También se puede ejecutar explícitamente así:

```bash
go run . web
```

Por defecto el servidor queda disponible en:

```text
http://localhost:8080
```

Para cambiar el puerto:

```bash
PUERTO=3000 go run . web
```

Para cambiar la ruta de la base de datos:

```bash
DB_RUTA=datos/ecommerce.db go run . web
```

## Ejecución de la CLI

La CLI se conserva como vestigio funcional del avance anterior:

```bash
go run . cli
```

## Servicios web REST JSON

### Salud del sistema

| Método | Ruta | Descripción |
|---|---|---|
| GET | `/api/salud` | Verifica que el backend esté activo. |

### Productos

| Método | Ruta | Descripción |
|---|---|---|
| GET | `/api/productos` | Lista productos activos. |
| POST | `/api/productos` | Crea un producto. |
| GET | `/api/productos/{id}` | Busca un producto por ID. |
| PUT | `/api/productos/{id}/stock` | Actualiza el stock por cambio positivo o negativo. |
| DELETE | `/api/productos/{id}` | Elimina un producto de forma lógica. |

Ejemplo:

```bash
curl -X POST http://localhost:8080/api/productos \
  -H "Content-Type: application/json" \
  -d '{"nombre":"Laptop","precio":850,"stock":10,"categoria":"Tecnología"}'
```

### Clientes

| Método | Ruta | Descripción |
|---|---|---|
| GET | `/api/clientes` | Lista clientes activos. |
| POST | `/api/clientes` | Registra un cliente. |
| GET | `/api/clientes/{id}` | Busca un cliente por ID. |
| DELETE | `/api/clientes/{id}` | Elimina un cliente de forma lógica. |

Ejemplo:

```bash
curl -X POST http://localhost:8080/api/clientes \
  -H "Content-Type: application/json" \
  -d '{"nombre":"Ana Pérez","email":"ana@example.com","telefono":"0991234567"}'
```

### Pedidos

| Método | Ruta | Descripción |
|---|---|---|
| GET | `/api/pedidos` | Lista pedidos. Acepta filtro `?estado=pendiente`. |
| POST | `/api/pedidos` | Crea un pedido asociado a un cliente. |
| GET | `/api/pedidos/{id}` | Busca un pedido por ID. |
| POST | `/api/pedidos/{id}/items` | Agrega un producto al pedido y descuenta stock. |
| GET | `/api/pedidos/{id}/total` | Calcula el total del pedido. |
| PUT | `/api/pedidos/{id}/estado` | Cambia el estado del pedido. |

Ejemplo:

```bash
curl -X POST http://localhost:8080/api/pedidos \
  -H "Content-Type: application/json" \
  -d '{"cliente_id":"CLI-123"}'
```

### Inventario

| Método | Ruta | Descripción |
|---|---|---|
| GET | `/api/inventario/stock?id=PROD-123&cantidad=2` | Verifica si hay stock y muestra cuánto stock existe. |
| GET | `/api/inventario/stock-bajo?umbral=5` | Lista productos con stock menor al umbral. |
| POST | `/api/inventario/reponer` | Repone stock de un producto. |
| GET | `/api/inventario/reporte?orden=stock` | Genera reporte ordenado por nombre, precio o stock. |

Ejemplo:

```bash
curl "http://localhost:8080/api/inventario/stock?id=PROD-123&cantidad=2"
```

### Concurrencia

| Método | Ruta | Descripción |
|---|---|---|
| GET | `/api/concurrencia/resumen-inventario?umbral=5&orden=stock` | Ejecuta en paralelo tres tareas: listar productos, calcular alertas de stock bajo y generar reporte. Usa goroutines y un canal de resultados. |

Ejemplo:

```bash
curl "http://localhost:8080/api/concurrencia/resumen-inventario?umbral=5&orden=stock"
```

## Relación con las unidades de la materia

| Unidad | Aplicación en el proyecto |
|---|---|
| Unidad 1 | Sintaxis de Go, condicionales, ciclos, manejo de errores y menú CLI. |
| Unidad 2 | Arrays para estados válidos, slices para listados, maps para selección de opciones y organización de datos. |
| Unidad 3 | Structs, métodos, constructores, encapsulación, interfaces y separación por paquetes. |
| Unidad 4 | Servicios web REST, JSON, concurrencia con goroutines/canales y testing. |

## Pruebas

Ejecutar todas las pruebas:

```bash
go test ./...
```

Pruebas incluidas:

- Validaciones del dominio en `internal/modelos`.
- Pruebas HTTP de endpoints usando `net/http/httptest`.
- Prueba de aceptación del flujo producto → cliente → pedido → item → stock.
- Prueba del endpoint concurrente de resumen de inventario.
- Prueba de integración de inicialización de tablas SQLite en memoria.

## Estructura del proyecto

```text
main.go
internal/
  cli/            # CLI conservada de la entrega anterior
  clientes/       # Servicio y repositorio de clientes
  inventario/     # Servicio de inventario
  modelos/        # Structs, constructores y validaciones del dominio
  pedidos/        # Servicio y repositorio de pedidos
  persistencia/   # Conexión e inicialización SQLite
  productos/      # Servicio y repositorio de productos
  web/            # Servidor REST JSON, handlers, DTOs, middleware y tests
docs/
  entrega_2.md
```

## Nota sobre clases en Go

Go no utiliza clases tradicionales. En este proyecto, los objetos del dominio se representan mediante `structs`, constructores, métodos e interfaces. Esta estructura permite aplicar conceptos de POO como encapsulación, responsabilidad única y desacoplamiento sin romper el estilo idiomático de Go.
