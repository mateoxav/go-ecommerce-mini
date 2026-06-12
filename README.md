# go-ecommerce-mini

Sistema de Gestión de E-Commerce desarrollado en Go como avance del Aprendizaje Autónomo 2.

## Alcance de esta entrega

Esta versión implementa una CLI funcional con persistencia SQLite local. Incluye:

- Menú principal y submenús por módulo.
- CRUD funcional de productos.
- Registro, búsqueda y listado de clientes.
- Creación de pedidos, agregado de ítems, cálculo de total y cambio de estado.
- Control de inventario con verificación de stock, alertas, reposición y reporte ordenado.
- Persistencia en `ecommerce.db` usando SQLite.
- Uso de structs, constructores, métodos, interfaces, arrays, slices y maps.
- Manejo de errores con retornos idiomáticos `error`.

## Requisitos

- Go 1.26.1 o compatible con el proyecto.
- GCC/CGO habilitado para compilar `github.com/mattn/go-sqlite3`.

## Instalación

```bash
git clone https://github.com/mateoxav/go-ecommerce-mini.git
cd go-ecommerce-mini
go mod tidy
go run .
```

## Estructura

```text
main.go
internal/
  cli/
  clientes/
  inventario/
  modelos/
  pedidos/
  persistencia/
  productos/
docs/
```

## Relación con los temas de la unidad

| Tema | Aplicación |
|---|---|
| Arrays | Estados fijos de pedido: pendiente, enviado, entregado, cancelado. |
| Slices | Listados dinámicos de productos, clientes, pedidos y alertas. |
| Maps | Enrutamiento del menú CLI y criterios de ordenamiento del inventario. |
| Structs | Producto, Cliente, Pedido e ItemPedido. |
| Métodos | Getters, servicios y repositorios con receivers. |
| Constructores | NuevoProducto, NuevoCliente, NuevoPedido y constructores de servicios/repositorios. |
| Interfaces | Repositorios desacoplados para aplicar SOLID. |
| Manejo de errores | Retorno y propagación de `error` con contexto. |

## Nota sobre clases en Go

Go no utiliza clases tradicionales. En este proyecto, los objetos del dominio se representan mediante `structs`, constructores, métodos e interfaces, que permiten aplicar principios de programación orientada a objetos respetando el estilo propio del lenguaje.
