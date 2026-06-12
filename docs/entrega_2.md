# Aprendizaje Autónomo 2 — Desarrollo del Sistema de Gestión

## Objetivo del avance

Continuar el desarrollo del Sistema de Gestión de E-Commerce iniciado en la etapa de planeación, implementando una aplicación CLI funcional en Go conectada a una base de datos SQLite local.

## Avance implementado

El sistema ya permite operar desde consola los módulos principales definidos en la arquitectura:

- Productos
- Clientes
- Pedidos
- Inventario
- Persistencia

## Funcionalidades desarrolladas

### Productos

- Crear producto.
- Listar productos activos.
- Buscar producto por ID.
- Actualizar stock.
- Eliminar producto mediante borrado lógico.

### Clientes

- Registrar cliente.
- Validar email.
- Buscar cliente por ID.
- Listar clientes.

### Pedidos

- Crear pedido asociado a un cliente.
- Agregar productos al pedido.
- Descontar stock automáticamente.
- Calcular total del pedido.
- Cambiar estado del pedido.
- Listar pedidos con filtro opcional por estado.

### Inventario

- Verificar stock disponible.
- Generar alertas de stock bajo.
- Reponer stock.
- Generar reporte ordenado por nombre, precio o stock.

## Aplicación de estructuras de datos

- Arrays: estados válidos del pedido.
- Slices: listados dinámicos consultados desde SQLite.
- Maps: menú CLI, validación de estados y criterios de ordenamiento.

## Aplicación de objetos en Go

Go no maneja clases tradicionales, por lo que el sistema usa:

- Structs para representar entidades del dominio.
- Constructores para crear objetos válidos.
- Métodos para encapsular lectura de datos.
- Interfaces para desacoplar servicios y repositorios.

## Principios SOLID aplicados

- Responsabilidad única: cada paquete se enfoca en un módulo del sistema.
- Abierto/cerrado: los servicios pueden ampliarse sin cambiar la CLI completa.
- Sustitución de Liskov: las interfaces permiten cambiar implementaciones de repositorio.
- Segregación de interfaces: cada módulo declara interfaces pequeñas según su necesidad.
- Inversión de dependencias: los servicios dependen de interfaces, no de SQLite directamente.

## Pendiente para entrega final

- Mejorar pruebas unitarias.
- Agregar capturas de ejecución.
- Refinar validaciones.
- Preparar video demostrativo.
- Completar documentación técnica del repositorio.
- Aplicación web, API y concurrencia.
