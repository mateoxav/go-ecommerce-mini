# Frontend · Go E-Commerce Mini

Frontend desarrollado con React, TypeScript, Vite, pnpm y Chakra UI 3.36.0 para consumir el backend REST JSON del proyecto `go-ecommerce-mini`.

## Tecnologías

- React + TypeScript
- Vite
- pnpm
- Chakra UI 3.36.0
- next-themes para alternar tema claro/oscuro
- Fetch API para consumir el backend

## Estructura

```text
frontend/
├── index.html
├── package.json
├── vite.config.ts
├── tsconfig.json
├── .env.example
└── src/
    ├── api/
    ├── components/
    ├── pages/
    ├── theme/
    ├── types/
    ├── App.tsx
    ├── main.tsx
    └── styles.css
```

## Instalación

Desde la carpeta `frontend`:

```bash
pnpm install
```

Crear el archivo `.env` a partir del ejemplo:

```bash
cp .env.example .env
```

En Windows PowerShell puedes usar:

```powershell
Copy-Item .env.example .env
```

Contenido esperado:

```env
VITE_API_URL=http://localhost:8080/api
```

## Ejecución

Primero inicia el backend Go desde la raíz del repositorio:

```bash
go run . web
```

Luego inicia el frontend:

```bash
cd frontend
pnpm dev
```

La interfaz queda disponible en:

```text
http://127.0.0.1:5173
```

## Build de producción

```bash
pnpm build
```

Vista previa del build:

```bash
pnpm preview
```

## Pantallas incluidas

- Dashboard
- Productos
- Clientes
- Pedidos
- Inventario
- Concurrencia

## Funcionalidades cubiertas

### Dashboard

- Verifica el estado del backend con `GET /api/salud`.
- Muestra resumen de productos, clientes, pedidos y alertas.

### Productos

- Listar productos.
- Crear productos.
- Buscar producto por ID.
- Actualizar stock con cambio positivo o negativo.
- Eliminar producto de forma lógica.

### Clientes

- Listar clientes.
- Registrar clientes.
- Buscar cliente por ID.
- Eliminar cliente de forma lógica.

### Pedidos

- Listar pedidos.
- Filtrar pedidos por estado.
- Crear pedido.
- Buscar pedido.
- Agregar ítem al pedido.
- Calcular total.
- Cambiar estado.

### Inventario

- Verificar stock y mostrar stock disponible actual.
- Consultar alertas de stock bajo.
- Reponer stock.
- Generar reporte ordenado por nombre, precio o stock.

### Concurrencia

- Ejecuta `GET /api/concurrencia/resumen-inventario`.
- Muestra total de productos, productos con stock bajo, reporte, tiempo de ejecución, goroutines y canal usado.

## Nota

Este frontend asume que el backend está corriendo en `http://localhost:8080` y que mantiene las rutas REST JSON generadas en la versión final del backend.
