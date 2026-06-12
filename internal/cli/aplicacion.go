package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mateoxav/go-ecommerce-mini/internal/clientes"
	"github.com/mateoxav/go-ecommerce-mini/internal/inventario"
	"github.com/mateoxav/go-ecommerce-mini/internal/pedidos"
	"github.com/mateoxav/go-ecommerce-mini/internal/productos"
)

type Aplicacion struct {
	entrada            *bufio.Reader
	servicioProductos  *productos.Servicio
	servicioClientes   *clientes.Servicio
	servicioPedidos    *pedidos.Servicio
	servicioInventario *inventario.Servicio
}

func NuevaAplicacion(
	servicioProductos *productos.Servicio,
	servicioClientes *clientes.Servicio,
	servicioPedidos *pedidos.Servicio,
	servicioInventario *inventario.Servicio,
) *Aplicacion {
	return &Aplicacion{
		entrada:            bufio.NewReader(os.Stdin),
		servicioProductos:  servicioProductos,
		servicioClientes:   servicioClientes,
		servicioPedidos:    servicioPedidos,
		servicioInventario: servicioInventario,
	}
}

func (a *Aplicacion) Ejecutar(ctx context.Context) error {
	fmt.Println("Sistema de Gestión de E-Commerce")

	acciones := map[string]func(context.Context) error{
		"1": a.menuProductos,
		"2": a.menuClientes,
		"3": a.menuPedidos,
		"4": a.menuInventario,
	}

	for {
		fmt.Println("\nMENÚ PRINCIPAL")
		fmt.Println("1. Gestión de productos")
		fmt.Println("2. Gestión de clientes")
		fmt.Println("3. Gestión de pedidos")
		fmt.Println("4. Gestión de inventario")
		fmt.Println("5. Salir")

		opcion := a.leerTexto("Seleccione una opción: ")
		if opcion == "5" {
			fmt.Println("Cerrando sistema. Base de datos guardada localmente.")
			return nil
		}

		accion, existe := acciones[opcion]
		if !existe {
			fmt.Println("Opción no válida")
			continue
		}

		if err := accion(ctx); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}
}

func (a *Aplicacion) menuProductos(ctx context.Context) error {
	for {
		fmt.Println("\nPRODUCTOS")
		fmt.Println("1. Crear producto")
		fmt.Println("2. Listar productos")
		fmt.Println("3. Buscar producto")
		fmt.Println("4. Actualizar stock")
		fmt.Println("5. Eliminar producto")
		fmt.Println("6. Volver")

		switch a.leerTexto("Seleccione una opción: ") {
		case "1":
			return a.crearProducto(ctx)
		case "2":
			return a.listarProductos(ctx)
		case "3":
			return a.buscarProducto(ctx)
		case "4":
			return a.actualizarStock(ctx)
		case "5":
			return a.eliminarProducto(ctx)
		case "6":
			return nil
		default:
			fmt.Println("Opción no válida")
		}
	}
}

func (a *Aplicacion) crearProducto(ctx context.Context) error {
	nombre := a.leerTexto("Nombre: ")
	precio := a.leerFloat("Precio: ")
	stock := a.leerEntero("Stock inicial: ")
	categoria := a.leerTexto("Categoría: ")

	producto, err := a.servicioProductos.CrearProducto(ctx, nombre, precio, stock, categoria)
	if err != nil {
		return err
	}

	fmt.Printf("Producto creado con ID: %s\n", producto.ID())
	return nil
}

func (a *Aplicacion) listarProductos(ctx context.Context) error {
	productos, err := a.servicioProductos.ListarProductos(ctx)
	if err != nil {
		return err
	}

	if len(productos) == 0 {
		fmt.Println("No hay productos registrados")
		return nil
	}

	fmt.Printf("%-24s %-20s %-12s %-10s %-12s\n", "ID", "Nombre", "Categoría", "Precio", "Stock")
	for _, p := range productos {
		fmt.Printf("%-24s %-20s %-12s $%-9.2f %-12d\n", p.ID(), p.Nombre(), p.Categoria(), p.Precio(), p.Stock())
	}
	return nil
}

func (a *Aplicacion) buscarProducto(ctx context.Context) error {
	id := a.leerTexto("ID del producto: ")
	producto, err := a.servicioProductos.BuscarProducto(ctx, id)
	if err != nil {
		return err
	}

	fmt.Printf("Producto: %s | Categoría: %s | Precio: $%.2f | Stock: %d\n",
		producto.Nombre(), producto.Categoria(), producto.Precio(), producto.Stock())
	return nil
}

func (a *Aplicacion) actualizarStock(ctx context.Context) error {
	id := a.leerTexto("ID del producto: ")
	cambio := a.leerEntero("Cambio de stock (+ para reponer, - para descontar): ")
	if err := a.servicioProductos.ActualizarStock(ctx, id, cambio); err != nil {
		return err
	}
	fmt.Println("Stock actualizado")
	return nil
}

func (a *Aplicacion) eliminarProducto(ctx context.Context) error {
	id := a.leerTexto("ID del producto: ")
	if err := a.servicioProductos.EliminarProducto(ctx, id); err != nil {
		return err
	}
	fmt.Println("Producto eliminado de forma lógica")
	return nil
}

func (a *Aplicacion) menuClientes(ctx context.Context) error {
	for {
		fmt.Println("\nCLIENTES")
		fmt.Println("1. Registrar cliente")
		fmt.Println("2. Listar clientes")
		fmt.Println("3. Buscar cliente")
		fmt.Println("4. Volver")

		switch a.leerTexto("Seleccione una opción: ") {
		case "1":
			return a.registrarCliente(ctx)
		case "2":
			return a.listarClientes(ctx)
		case "3":
			return a.buscarCliente(ctx)
		case "4":
			return nil
		default:
			fmt.Println("Opción no válida")
		}
	}
}

func (a *Aplicacion) registrarCliente(ctx context.Context) error {
	nombre := a.leerTexto("Nombre: ")
	email := a.leerTexto("Email: ")
	telefono := a.leerTexto("Teléfono: ")

	cliente, err := a.servicioClientes.RegistrarCliente(ctx, nombre, email, telefono)
	if err != nil {
		return err
	}

	fmt.Printf("Cliente registrado con ID: %s\n", cliente.ID())
	return nil
}

func (a *Aplicacion) listarClientes(ctx context.Context) error {
	clientes, err := a.servicioClientes.ListarClientes(ctx)
	if err != nil {
		return err
	}

	if len(clientes) == 0 {
		fmt.Println("No hay clientes registrados")
		return nil
	}

	fmt.Printf("%-24s %-22s %-28s %-15s\n", "ID", "Nombre", "Email", "Teléfono")
	for _, c := range clientes {
		fmt.Printf("%-24s %-22s %-28s %-15s\n", c.ID(), c.Nombre(), c.Email(), c.Telefono())
	}
	return nil
}

func (a *Aplicacion) buscarCliente(ctx context.Context) error {
	id := a.leerTexto("ID del cliente: ")
	cliente, err := a.servicioClientes.BuscarCliente(ctx, id)
	if err != nil {
		return err
	}

	fmt.Printf("Cliente: %s | Email: %s | Teléfono: %s\n", cliente.Nombre(), cliente.Email(), cliente.Telefono())
	return nil
}

func (a *Aplicacion) menuPedidos(ctx context.Context) error {
	for {
		fmt.Println("\nPEDIDOS")
		fmt.Println("1. Crear pedido")
		fmt.Println("2. Agregar item a pedido")
		fmt.Println("3. Calcular total")
		fmt.Println("4. Cambiar estado")
		fmt.Println("5. Listar pedidos")
		fmt.Println("6. Volver")

		switch a.leerTexto("Seleccione una opción: ") {
		case "1":
			return a.crearPedido(ctx)
		case "2":
			return a.agregarItemPedido(ctx)
		case "3":
			return a.calcularTotalPedido(ctx)
		case "4":
			return a.cambiarEstadoPedido(ctx)
		case "5":
			return a.listarPedidos(ctx)
		case "6":
			return nil
		default:
			fmt.Println("Opción no válida")
		}
	}
}

func (a *Aplicacion) crearPedido(ctx context.Context) error {
	clienteID := a.leerTexto("ID del cliente: ")
	pedido, err := a.servicioPedidos.CrearPedido(ctx, clienteID)
	if err != nil {
		return err
	}
	fmt.Printf("Pedido creado con ID: %s\n", pedido.ID())
	return nil
}

func (a *Aplicacion) agregarItemPedido(ctx context.Context) error {
	pedidoID := a.leerTexto("ID del pedido: ")
	productoID := a.leerTexto("ID del producto: ")
	cantidad := a.leerEntero("Cantidad: ")

	if err := a.servicioPedidos.AgregarItem(ctx, pedidoID, productoID, cantidad); err != nil {
		return err
	}
	fmt.Println("Item agregado y stock actualizado")
	return nil
}

func (a *Aplicacion) calcularTotalPedido(ctx context.Context) error {
	pedidoID := a.leerTexto("ID del pedido: ")
	total, err := a.servicioPedidos.CalcularTotal(ctx, pedidoID)
	if err != nil {
		return err
	}
	fmt.Printf("Total del pedido: $%.2f\n", total)
	return nil
}

func (a *Aplicacion) cambiarEstadoPedido(ctx context.Context) error {
	pedidoID := a.leerTexto("ID del pedido: ")
	estado := a.leerTexto("Nuevo estado (pendiente/enviado/entregado/cancelado): ")

	if err := a.servicioPedidos.CambiarEstado(ctx, pedidoID, estado); err != nil {
		return err
	}
	fmt.Println("Estado actualizado")
	return nil
}

func (a *Aplicacion) listarPedidos(ctx context.Context) error {
	filtro := a.leerTexto("Filtrar por estado, o vacío para todos: ")
	pedidos, err := a.servicioPedidos.ListarPedidos(ctx, filtro)
	if err != nil {
		return err
	}

	if len(pedidos) == 0 {
		fmt.Println("No hay pedidos registrados")
		return nil
	}

	fmt.Printf("%-24s %-24s %-12s %-12s %-25s\n", "ID", "Cliente", "Total", "Estado", "Fecha")
	for _, p := range pedidos {
		fmt.Printf("%-24s %-24s $%-11.2f %-12s %-25s\n", p.ID(), p.ClienteID(), p.Total(), p.Estado(), p.Fecha())
	}
	return nil
}

func (a *Aplicacion) menuInventario(ctx context.Context) error {
	for {
		fmt.Println("\nINVENTARIO")
		fmt.Println("1. Verificar stock")
		fmt.Println("2. Alertas de stock bajo")
		fmt.Println("3. Reponer stock")
		fmt.Println("4. Generar reporte")
		fmt.Println("5. Volver")

		switch a.leerTexto("Seleccione una opción: ") {
		case "1":
			return a.verificarStock(ctx)
		case "2":
			return a.alertasStockBajo(ctx)
		case "3":
			return a.reponerStock(ctx)
		case "4":
			return a.generarReporteInventario(ctx)
		case "5":
			return nil
		default:
			fmt.Println("Opción no válida")
		}
	}
}

func (a *Aplicacion) verificarStock(ctx context.Context) error {
	productoID := a.leerTexto("ID del producto: ")
	cantidad := a.leerEntero("Cantidad a verificar: ")

	disponible, err := a.servicioInventario.VerificarStock(ctx, productoID, cantidad)
	if err != nil {
		return err
	}
	if disponible {
		fmt.Println("Hay stock suficiente")
	} else {
		fmt.Println("No hay stock suficiente")
	}
	return nil
}

func (a *Aplicacion) alertasStockBajo(ctx context.Context) error {
	umbral := a.leerEntero("Umbral de stock bajo: ")
	productos, err := a.servicioInventario.AlertasStockBajo(ctx, umbral)
	if err != nil {
		return err
	}

	if len(productos) == 0 {
		fmt.Println("No hay alertas de stock bajo")
		return nil
	}

	fmt.Println("Productos con stock bajo:")
	for _, p := range productos {
		fmt.Printf("- %s | Stock: %d\n", p.Nombre(), p.Stock())
	}
	return nil
}

func (a *Aplicacion) reponerStock(ctx context.Context) error {
	productoID := a.leerTexto("ID del producto: ")
	cantidad := a.leerEntero("Cantidad a reponer: ")

	if err := a.servicioInventario.ReponerStock(ctx, productoID, cantidad); err != nil {
		return err
	}
	fmt.Println("Stock repuesto")
	return nil
}

func (a *Aplicacion) generarReporteInventario(ctx context.Context) error {
	orden := a.leerTexto("Ordenar por nombre, precio o stock: ")
	reporte, err := a.servicioInventario.GenerarReporte(ctx, orden)
	if err != nil {
		return err
	}
	fmt.Println(reporte)
	return nil
}

func (a *Aplicacion) leerTexto(etiqueta string) string {
	fmt.Print(etiqueta)
	texto, _ := a.entrada.ReadString('\n')
	return strings.TrimSpace(texto)
}

func (a *Aplicacion) leerEntero(etiqueta string) int {
	for {
		valor := a.leerTexto(etiqueta)
		numero, err := strconv.Atoi(valor)
		if err == nil {
			return numero
		}
		fmt.Println("Ingrese un número entero válido")
	}
}

func (a *Aplicacion) leerFloat(etiqueta string) float64 {
	for {
		valor := a.leerTexto(etiqueta)
		numero, err := strconv.ParseFloat(valor, 64)
		if err == nil {
			return numero
		}
		fmt.Println("Ingrese un número decimal válido")
	}
}
