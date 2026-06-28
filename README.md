# go-ecommerce-mini

Sistema de Gestión de E-Commerce desarrollado en Go para el Proyecto Final de Programación Orientada a Objetos. El proyecto representa la integración del trabajo realizado durante las 8 semanas de la materia. Inicialmente fue una aplicación de consola con SQLite y, en esta última etapa, se amplió principalmente como un backend web REST JSON, conservando la CLI como evidencia de la evolución del sistema.

## Datos generales

* **Repositorio:** https://github.com/mateoxav/go-ecommerce-mini
* **Materia:** Programación Orientada a Objetos
* **Proyecto:** Sistema de Gestión de E-Commerce
* **Fecha:** Junio de 2026
* **Lenguaje principal:** Go 1.26
* **Frontend complementario:** React + Vite + TypeScript + Chakra UI
* **Persistencia:** SQLite local mediante `database/sql` y `github.com/mattn/go-sqlite3`
* **Servidor web:** Biblioteca estándar `net/http`
* **Serialización de datos:** JSON mediante `encoding/json`
* **Tipo de aplicación:** Sistema web de gestión para e-commerce

## Objetivo del programa

El objetivo del proyecto es implementar un sistema de gestión de e-commerce que permita administrar productos, clientes, pedidos e inventario, integrando los conocimientos revisados durante las 4 unidades de la materia.

El sistema permite realizar operaciones básicas de administración, exponer servicios web REST, serializar datos mediante JSON, persistir información en SQLite, aplicar estructuras de datos de Go, usar conceptos de orientación a objetos mediante structs, métodos, constructores e interfaces, y demostrar concurrencia mediante goroutines y canales.

## Justificación de la aplicación seleccionada

Se eligió un sistema de gestión de e-commerce porque es un caso práctico y realista que permite aplicar varios conceptos importantes de programación orientada a objetos. Este tipo de sistema necesita manejar productos, clientes, pedidos, stock, reportes y operaciones de consulta o actualización, por lo que resulta adecuado para demostrar estructuras de datos, separación por módulos, persistencia, servicios web y pruebas.

Además, el dominio de e-commerce facilita la creación de más de 8 servicios web diferentes, lo cual se ajusta a los requisitos de la actividad final.

## Funcionalidades principales

* Backend REST JSON con más de 8 servicios web.
* CRUD funcional de productos con eliminación lógica.
* Registro, búsqueda, listado y eliminación lógica de clientes.
* Creación de pedidos asociados a clientes.
* Agregado de ítems a pedidos con descuento automático de stock.
* Cálculo de total de pedidos.
* Cambio de estado de pedidos.
* Control de inventario con verificación de stock, alertas de stock bajo, reposición y reportes.
* Endpoint concurrente que ejecuta tareas paralelas usando goroutines y canales.
* CLI conservada como modo alternativo para evidenciar el avance progresivo del proyecto.
* Base de datos SQLite inicializada automáticamente.
* Frontend web simple para probar visualmente las funcionalidades del backend.
* Tests unitarios, pruebas HTTP con `httptest`, prueba de aceptación y prueba del endpoint concurrente.

## Requisitos

* Go 1.26.x o una versión compatible con el proyecto.
* GCC/CGO habilitado para compilar `github.com/mattn/go-sqlite3`.
* Node.js compatible con Vite.
* pnpm para instalar y ejecutar el frontend.
* SQLite no requiere servidor externo porque se usa como archivo local.

## Instalación del backend

```bash
git clone https://github.com/mateoxav/go-ecommerce-mini.git
cd go-ecommerce-mini
go mod tidy
```

## Ejecución del backend web

El modo principal del proyecto final es el backend web:

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

La CLI se conserva como vestigio funcional de las entregas anteriores. No es el modo principal de la entrega final, pero se mantiene porque muestra cómo fue progresando el sistema.

```bash
go run . cli
```

## Instalación y ejecución del frontend

El frontend se encuentra en la carpeta `frontend/` y fue creado para probar el backend de forma visual.

```bash
cd frontend
pnpm install
```

Crear el archivo `.env` a partir del ejemplo:

```bash
cp .env.example .env
```

En Windows PowerShell:

```powershell
Copy-Item .env.example .env
```

El archivo `.env` debe contener:

```env
VITE_API_URL=http://localhost:8080/api
```

Ejecutar el frontend:

```bash
pnpm dev
```

Por defecto estará disponible en:

```text
http://127.0.0.1:5173
```

## Servicios web REST JSON

El proyecto cumple el requisito de crear al menos 8 servicios web. En esta versión se implementan más de 8 endpoints, todos orientados a trabajar con datos en formato JSON.

### Salud del sistema

| Método | Ruta         | Descripción                          |
| ------ | ------------ | ------------------------------------ |
| GET    | `/api/salud` | Verifica que el backend esté activo. |

### Productos

| Método | Ruta                        | Descripción                                        |
| ------ | --------------------------- | -------------------------------------------------- |
| GET    | `/api/productos`            | Lista productos activos.                           |
| POST   | `/api/productos`            | Crea un producto.                                  |
| GET    | `/api/productos/{id}`       | Busca un producto por ID.                          |
| PUT    | `/api/productos/{id}/stock` | Actualiza el stock por cambio positivo o negativo. |
| DELETE | `/api/productos/{id}`       | Elimina un producto de forma lógica.               |

Ejemplo para crear producto:

```bash
curl -X POST http://localhost:8080/api/productos \
  -H "Content-Type: application/json" \
  -d '{"nombre":"Laptop","precio":850,"stock":10,"categoria":"Tecnología"}'
```

Ejemplo para listar productos:

```bash
curl http://localhost:8080/api/productos
```

### Clientes

| Método | Ruta                 | Descripción                         |
| ------ | -------------------- | ----------------------------------- |
| GET    | `/api/clientes`      | Lista clientes activos.             |
| POST   | `/api/clientes`      | Registra un cliente.                |
| GET    | `/api/clientes/{id}` | Busca un cliente por ID.            |
| DELETE | `/api/clientes/{id}` | Elimina un cliente de forma lógica. |

Ejemplo para registrar cliente:

```bash
curl -X POST http://localhost:8080/api/clientes \
  -H "Content-Type: application/json" \
  -d '{"nombre":"Ana Pérez","email":"ana@example.com","telefono":"0991234567"}'
```

### Pedidos

| Método | Ruta                       | Descripción                                       |
| ------ | -------------------------- | ------------------------------------------------- |
| GET    | `/api/pedidos`             | Lista pedidos. Acepta filtro `?estado=pendiente`. |
| POST   | `/api/pedidos`             | Crea un pedido asociado a un cliente.             |
| GET    | `/api/pedidos/{id}`        | Busca un pedido por ID.                           |
| POST   | `/api/pedidos/{id}/items`  | Agrega un producto al pedido y descuenta stock.   |
| GET    | `/api/pedidos/{id}/total`  | Calcula el total del pedido.                      |
| PUT    | `/api/pedidos/{id}/estado` | Cambia el estado del pedido.                      |

Ejemplo para crear pedido:

```bash
curl -X POST http://localhost:8080/api/pedidos \
  -H "Content-Type: application/json" \
  -d '{"cliente_id":"CLI-123"}'
```

Ejemplo para agregar un producto al pedido:

```bash
curl -X POST http://localhost:8080/api/pedidos/PED-123/items \
  -H "Content-Type: application/json" \
  -d '{"producto_id":"PROD-123","cantidad":2}'
```

### Inventario

| Método | Ruta                                           | Descripción                                          |
| ------ | ---------------------------------------------- | ---------------------------------------------------- |
| GET    | `/api/inventario/stock?id=PROD-123&cantidad=2` | Verifica si hay stock y muestra cuánto stock existe. |
| GET    | `/api/inventario/stock-bajo?umbral=5`          | Lista productos con stock menor al umbral.           |
| POST   | `/api/inventario/reponer`                      | Repone stock de un producto.                         |
| GET    | `/api/inventario/reporte?orden=stock`          | Genera reporte ordenado por nombre, precio o stock.  |

Ejemplo para verificar stock:

```bash
curl "http://localhost:8080/api/inventario/stock?id=PROD-123&cantidad=2"
```

Ejemplo para reponer stock:

```bash
curl -X POST http://localhost:8080/api/inventario/reponer \
  -H "Content-Type: application/json" \
  -d '{"producto_id":"PROD-123","cantidad":5}'
```

### Concurrencia

| Método | Ruta                                                        | Descripción                                                                          |
| ------ | ----------------------------------------------------------- | ------------------------------------------------------------------------------------ |
| GET    | `/api/concurrencia/resumen-inventario?umbral=5&orden=stock` | Ejecuta tareas de inventario en paralelo usando goroutines y un canal de resultados. |

Ejemplo:

```bash
curl "http://localhost:8080/api/concurrencia/resumen-inventario?umbral=5&orden=stock"
```

Este endpoint ejecuta de forma concurrente:

1. Listado de productos.
2. Cálculo de productos con stock bajo.
3. Generación del reporte de inventario.

La respuesta incluye información sobre las goroutines ejecutadas y el canal utilizado, para que la concurrencia pueda demostrarse de forma visible en la entrega.

## Aplicación de JSON

La serialización y deserialización de datos se realiza con JSON. Las solicitudes `POST` y `PUT` reciben cuerpos JSON, mientras que las respuestas del backend también se devuelven en JSON.

Ejemplo de JSON para crear un producto:

```json
{
  "nombre": "Laptop",
  "precio": 850,
  "stock": 10,
  "categoria": "Tecnología"
}
```

Ejemplo de respuesta JSON:

```json
{
  "id": "PROD-123",
  "nombre": "Laptop",
  "precio": 850,
  "stock": 10,
  "categoria": "Tecnología",
  "activo": true
}
```

## Relación con las unidades de la materia

| Unidad   | Aplicación en el proyecto                                                                                             |
| -------- | --------------------------------------------------------------------------------------------------------------------- |
| Unidad 1 | Sintaxis de Go, variables, condicionales, ciclos, funciones, manejo de errores y menú CLI.                            |
| Unidad 2 | Arrays para estados válidos, slices para listados dinámicos, maps para selección de opciones y organización de datos. |
| Unidad 3 | Structs, métodos, constructores, encapsulación, interfaces, separación por paquetes y principios SOLID básicos.       |
| Unidad 4 | Servicios web REST, serialización JSON, concurrencia con goroutines/canales y testing.                                |

## Programación orientada a objetos en Go

Go no utiliza clases tradicionales como otros lenguajes orientados a objetos. Por eso, en este proyecto los objetos del dominio se representan mediante:

* `structs`, como `Producto`, `Cliente`, `Pedido` e `ItemPedido`.
* Constructores, como `NuevoProducto`, `NuevoCliente` y `NuevoPedido`.
* Métodos para acceder a datos de forma encapsulada.
* Interfaces para desacoplar servicios y repositorios.
* Paquetes separados por responsabilidad.

Esta estructura permite aplicar conceptos de POO respetando el estilo propio de Go.

## Concurrencia

La concurrencia se implementa principalmente en el endpoint:

```text
GET /api/concurrencia/resumen-inventario
```

Dentro de este servicio se usan goroutines para ejecutar tareas en paralelo y un canal para recoger los resultados. También se usa `sync.WaitGroup` para esperar a que las tareas concurrentes terminen antes de construir la respuesta final.

Esto permite demostrar de forma práctica los temas de la Unidad 4:

* Introducción a la concurrencia.
* Goroutines.
* Canales.

## Pruebas

Ejecutar todas las pruebas:

```bash
go test ./...
```

Pruebas incluidas:

* Validaciones del dominio en `internal/modelos`.
* Pruebas HTTP de endpoints usando `net/http/httptest`.
* Prueba de aceptación del flujo producto → cliente → pedido → item → stock.
* Prueba del endpoint concurrente de resumen de inventario.
* Prueba de integración de inicialización de tablas SQLite en memoria.

## Tipos de pruebas realizadas

| Tipo de prueba         | Aplicación en el proyecto                                                                                        |
| ---------------------- | ---------------------------------------------------------------------------------------------------------------- |
| Pruebas unitarias      | Validan funciones concretas del dominio, como nombres, teléfonos y prefijos de ID.                               |
| Pruebas de integración | Verifican la inicialización de SQLite y el funcionamiento conjunto de algunos componentes.                       |
| Pruebas HTTP           | Validan handlers del backend REST usando `httptest`.                                                             |
| Prueba de aceptación   | Comprueba un flujo completo de uso: crear producto, crear cliente, crear pedido, agregar ítem y verificar stock. |
| Prueba de concurrencia | Valida el endpoint que ejecuta tareas paralelas con goroutines y canales.                                        |

## Estructura del proyecto

```text
main.go
go.mod
README.md
internal/
  cli/            # CLI conservada de la entrega anterior
  clientes/       # Servicio y repositorio de clientes
  inventario/     # Servicio de inventario
  modelos/        # Structs, constructores y validaciones del dominio
  pedidos/        # Servicio y repositorio de pedidos
  persistencia/   # Conexión e inicialización SQLite
  productos/      # Servicio y repositorio de productos
  web/            # Servidor REST JSON, handlers, DTOs, middleware y tests
frontend/
  src/
    api/          # Cliente HTTP para consumir el backend
    components/   # Componentes reutilizables de interfaz
    pages/        # Pantallas principales
    theme/        # Configuración visual
    types/        # Tipos TypeScript del dominio
docs/
  entrega_2.md
```

## Estado del proyecto

El proyecto cumple con los puntos principales solicitados en la entrega final:

* Código fuente en repositorio GitHub.
* README con explicación del sistema, objetivo, funcionalidades y fecha.
* Aplicación web con servicios REST.
* Más de 8 servicios web.
* Serialización mediante JSON.
* Integración de las 4 unidades de la materia.
* Concurrencia mediante goroutines y canales.
* Pruebas de software.
* Frontend complementario para demostrar el funcionamiento de la API.
* CLI conservada como evidencia del avance progresivo.

## Nota final

Este proyecto fue desarrollado con un enfoque académico, pero buscando mantener una estructura clara y cercana a buenas prácticas reales de desarrollo. La idea principal fue no solo cumplir con los requisitos, sino también construir un sistema que pueda explicarse y demostrarse de forma ordenada en el video final.
