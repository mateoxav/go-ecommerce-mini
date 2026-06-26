package cli

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mateoxav/go-ecommerce-mini/internal/clientes"
	"github.com/mateoxav/go-ecommerce-mini/internal/inventario"
	"github.com/mateoxav/go-ecommerce-mini/internal/modelos"
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
			a.ejecutarAccion(ctx, a.crearProducto)
		case "2":
			a.ejecutarAccion(ctx, a.listarProductos)
		case "3":
			a.ejecutarAccion(ctx, a.buscarProducto)
		case "4":
			a.ejecutarAccion(ctx, a.actualizarStock)
		case "5":
			a.ejecutarAccion(ctx, a.eliminarProducto)
		case "6":
			return nil
		default:
			fmt.Println("Opción no válida")
		}
	}
}

func (a *Aplicacion) crearProducto(ctx context.Context) error {
	nombre := a.leerTextoValidado("Nombre: ", validarNoVacio("el nombre del producto es obligatorio"))
	precio := a.leerFloatValidado("Precio: ", func(valor float64) error {
		if valor < 0 {
			return fmt.Errorf("el precio no puede ser negativo")
		}
		return nil
	})
	stock := a.leerEnteroValidado("Stock inicial: ", func(valor int) error {
		if valor < 0 {
			return fmt.Errorf("el stock no puede ser negativo")
		}
		return nil
	})
	categoria := a.leerTextoValidado("Categoría: ", validarNoVacio("la categoría es obligatoria"))

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
	id := a.leerIDProductoExistente(ctx)
	producto, err := a.servicioProductos.BuscarProducto(ctx, id)
	if err != nil {
		return err
	}

	fmt.Printf("Producto: %s | Categoría: %s | Precio: $%.2f | Stock: %d\n",
		producto.Nombre(), producto.Categoria(), producto.Precio(), producto.Stock())
	return nil
}

func (a *Aplicacion) actualizarStock(ctx context.Context) error {
	id := a.leerIDProductoExistente(ctx)
	cambio := a.leerEnteroValidado("Cambio de stock (+ para reponer, - para descontar): ", func(valor int) error {
		if valor == 0 {
			return fmt.Errorf("el cambio de stock no puede ser cero")
		}
		return nil
	})
	if err := a.servicioProductos.ActualizarStock(ctx, id, cambio); err != nil {
		return err
	}
	fmt.Println("Stock actualizado")
	return nil
}

func (a *Aplicacion) eliminarProducto(ctx context.Context) error {
	id := a.leerIDProductoExistente(ctx)
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
		fmt.Println("4. Eliminar cliente")
		fmt.Println("5. Volver")

		switch a.leerTexto("Seleccione una opción: ") {
		case "1":
			a.ejecutarAccion(ctx, a.registrarCliente)
		case "2":
			a.ejecutarAccion(ctx, a.listarClientes)
		case "3":
			a.ejecutarAccion(ctx, a.buscarCliente)
		case "4":
			a.ejecutarAccion(ctx, a.eliminarCliente)
		case "5":
			return nil
		default:
			fmt.Println("Opción no válida")
		}
	}
}

func (a *Aplicacion) registrarCliente(ctx context.Context) error {
	nombre := a.leerTextoValidado("Nombre: ", func(valor string) error {
		if !modelos.ValidarNombrePersona(valor) {
			return fmt.Errorf("el nombre es obligatorio y no puede contener números")
		}
		return nil
	})
	email := a.leerTextoValidado("Email: ", func(valor string) error {
		if !modelos.ValidarEmail(valor) {
			return fmt.Errorf("el email no tiene un formato válido")
		}
		return nil
	})
	telefono := a.leerTextoValidado("Teléfono (solo números, máximo 10 dígitos; puede quedar vacío): ", func(valor string) error {
		if !modelos.ValidarTelefono(valor) {
			return fmt.Errorf("el teléfono debe contener solo números y máximo 10 dígitos")
		}
		return nil
	})

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
	id := a.leerIDClienteExistente(ctx)
	cliente, err := a.servicioClientes.BuscarCliente(ctx, id)
	if err != nil {
		return err
	}

	fmt.Printf("Cliente: %s | Email: %s | Teléfono: %s\n", cliente.Nombre(), cliente.Email(), cliente.Telefono())
	return nil
}

func (a *Aplicacion) eliminarCliente(ctx context.Context) error {
	id := a.leerIDClienteExistente(ctx)
	if err := a.servicioClientes.EliminarCliente(ctx, id); err != nil {
		return err
	}
	fmt.Println("Cliente eliminado de forma lógica")
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
			a.ejecutarAccion(ctx, a.crearPedido)
		case "2":
			a.ejecutarAccion(ctx, a.agregarItemPedido)
		case "3":
			a.ejecutarAccion(ctx, a.calcularTotalPedido)
		case "4":
			a.ejecutarAccion(ctx, a.cambiarEstadoPedido)
		case "5":
			a.ejecutarAccion(ctx, a.listarPedidos)
		case "6":
			return nil
		default:
			fmt.Println("Opción no válida")
		}
	}
}

func (a *Aplicacion) crearPedido(ctx context.Context) error {
	clienteID := a.leerIDClienteExistente(ctx)
	pedido, err := a.servicioPedidos.CrearPedido(ctx, clienteID)
	if err != nil {
		return err
	}
	fmt.Printf("Pedido creado con ID: %s\n", pedido.ID())
	return nil
}

func (a *Aplicacion) agregarItemPedido(ctx context.Context) error {
	pedidoID := a.leerIDPedidoExistente(ctx)
	productoID := a.leerIDProductoExistente(ctx)
	cantidad := a.leerEnteroValidado("Cantidad: ", func(valor int) error {
		if valor <= 0 {
			return fmt.Errorf("la cantidad debe ser mayor que cero")
		}
		disponible, err := a.servicioInventario.VerificarStock(ctx, productoID, valor)
		if err != nil {
			return err
		}
		if !disponible {
			return fmt.Errorf("no hay stock suficiente para esa cantidad")
		}
		return nil
	})

	if err := a.servicioPedidos.AgregarItem(ctx, pedidoID, productoID, cantidad); err != nil {
		return err
	}
	fmt.Println("Item agregado y stock actualizado")
	return nil
}

func (a *Aplicacion) calcularTotalPedido(ctx context.Context) error {
	pedidoID := a.leerIDPedidoExistente(ctx)
	total, err := a.servicioPedidos.CalcularTotal(ctx, pedidoID)
	if err != nil {
		return err
	}
	fmt.Printf("Total del pedido: $%.2f\n", total)
	return nil
}

func (a *Aplicacion) cambiarEstadoPedido(ctx context.Context) error {
	pedidoID := a.leerIDPedidoExistente(ctx)
	estado := a.seleccionarEstadoPedido("Nuevo estado")

	if err := a.servicioPedidos.CambiarEstado(ctx, pedidoID, estado); err != nil {
		return err
	}
	fmt.Println("Estado actualizado")
	return nil
}

func (a *Aplicacion) listarPedidos(ctx context.Context) error {
	filtro := a.seleccionarFiltroEstadoPedido()
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
			a.ejecutarAccion(ctx, a.verificarStock)
		case "2":
			a.ejecutarAccion(ctx, a.alertasStockBajo)
		case "3":
			a.ejecutarAccion(ctx, a.reponerStock)
		case "4":
			a.ejecutarAccion(ctx, a.generarReporteInventario)
		case "5":
			return nil
		default:
			fmt.Println("Opción no válida")
		}
	}
}

func (a *Aplicacion) verificarStock(ctx context.Context) error {
	productoID := a.leerIDProductoExistente(ctx)
	producto, err := a.servicioProductos.BuscarProducto(ctx, productoID)
	if err != nil {
		return err
	}

	cantidad := a.leerEnteroValidado("Cantidad a verificar: ", func(valor int) error {
		if valor <= 0 {
			return fmt.Errorf("la cantidad debe ser mayor que cero")
		}
		return nil
	})

	disponible, err := a.servicioInventario.VerificarStock(ctx, productoID, cantidad)
	if err != nil {
		return err
	}
	if disponible {
		fmt.Printf("Hay stock suficiente. Stock disponible actual: %d unidades\n", producto.Stock())
	} else {
		fmt.Printf("No hay stock suficiente. Stock disponible actual: %d unidades\n", producto.Stock())
	}
	return nil
}

func (a *Aplicacion) alertasStockBajo(ctx context.Context) error {
	umbral := a.leerEnteroValidado("Umbral de stock bajo: ", func(valor int) error {
		if valor < 0 {
			return fmt.Errorf("el umbral no puede ser negativo")
		}
		return nil
	})
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
	productoID := a.leerIDProductoExistente(ctx)
	cantidad := a.leerEnteroValidado("Cantidad a reponer: ", func(valor int) error {
		if valor <= 0 {
			return fmt.Errorf("la cantidad a reponer debe ser mayor que cero")
		}
		return nil
	})

	if err := a.servicioInventario.ReponerStock(ctx, productoID, cantidad); err != nil {
		return err
	}
	fmt.Println("Stock repuesto")
	return nil
}

func (a *Aplicacion) generarReporteInventario(ctx context.Context) error {
	orden := a.seleccionarOrdenInventario()
	reporte, err := a.servicioInventario.GenerarReporte(ctx, orden)
	if err != nil {
		return err
	}
	fmt.Println(reporte)
	return nil
}

func (a *Aplicacion) ejecutarAccion(ctx context.Context, accion func(context.Context) error) {
	if err := accion(ctx); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

func (a *Aplicacion) leerIDClienteExistente(ctx context.Context) string {
	return a.leerTextoValidado("ID del cliente: ", func(valor string) error {
		if !modelos.ValidarIDCliente(valor) {
			return modelos.ErrorIDClienteInvalido()
		}
		if _, err := a.servicioClientes.BuscarCliente(ctx, valor); err != nil {
			return err
		}
		return nil
	})
}

func (a *Aplicacion) leerIDProductoExistente(ctx context.Context) string {
	return a.leerTextoValidado("ID del producto: ", func(valor string) error {
		if !modelos.ValidarIDProducto(valor) {
			return modelos.ErrorIDProductoInvalido()
		}
		if _, err := a.servicioProductos.BuscarProducto(ctx, valor); err != nil {
			return err
		}
		return nil
	})
}

func (a *Aplicacion) leerIDPedidoExistente(ctx context.Context) string {
	return a.leerTextoValidado("ID del pedido: ", func(valor string) error {
		if !modelos.ValidarIDPedido(valor) {
			return modelos.ErrorIDPedidoInvalido()
		}
		if _, err := a.servicioPedidos.BuscarPedido(ctx, valor); err != nil {
			return err
		}
		return nil
	})
}

func (a *Aplicacion) seleccionarFiltroEstadoPedido() string {
	for {
		opcion := strings.ToLower(a.leerTexto("Filtrar pedidos por estado(e) o vacío/todos(v): "))
		switch opcion {
		case "v", "", "todos", "todo":
			return ""
		case "e", "estado":
			return a.seleccionarEstadoPedido("Estado a filtrar")
		default:
			fmt.Println("Opción no válida. Use e para elegir estado o v para ver todos.")
		}
	}
}

func (a *Aplicacion) seleccionarEstadoPedido(etiqueta string) string {
	estados := map[string]string{
		"p":         "pendiente",
		"pendiente": "pendiente",
		"e":         "enviado",
		"enviado":   "enviado",
		"t":         "entregado",
		"entregado": "entregado",
		"c":         "cancelado",
		"cancelado": "cancelado",
	}

	for {
		valor := strings.ToLower(a.leerTexto(fmt.Sprintf("%s: pendiente(p), enviado(e), entregado(t), cancelado(c): ", etiqueta)))
		estado, existe := estados[valor]
		if existe {
			return estado
		}
		fmt.Println("Estado no válido")
	}
}

func (a *Aplicacion) seleccionarOrdenInventario() string {
	ordenes := map[string]string{
		"n":      "nombre",
		"nombre": "nombre",
		"p":      "precio",
		"precio": "precio",
		"s":      "stock",
		"stock":  "stock",
	}

	for {
		valor := strings.ToLower(a.leerTexto("Ordenar por nombre(n), precio(p) o stock(s): "))
		orden, existe := ordenes[valor]
		if existe {
			return orden
		}
		fmt.Println("Opción no válida. Use n, p o s.")
	}
}

func (a *Aplicacion) leerTexto(etiqueta string) string {
	fmt.Print(etiqueta)
	texto, _ := a.entrada.ReadString('\n')
	return strings.TrimSpace(texto)
}

func (a *Aplicacion) leerTextoValidado(etiqueta string, validar func(string) error) string {
	for {
		valor := a.leerTexto(etiqueta)
		if err := validar(valor); err != nil {
			fmt.Printf("Valor inválido: %v\n", err)
			continue
		}
		return strings.TrimSpace(valor)
	}
}

func (a *Aplicacion) leerEnteroValidado(etiqueta string, validar func(int) error) int {
	for {
		valor := a.leerTexto(etiqueta)
		numero, err := strconv.Atoi(valor)
		if err != nil {
			fmt.Println("Ingrese un número entero válido")
			continue
		}
		if err := validar(numero); err != nil {
			fmt.Printf("Valor inválido: %v\n", err)
			continue
		}
		return numero
	}
}

func (a *Aplicacion) leerFloatValidado(etiqueta string, validar func(float64) error) float64 {
	for {
		valor := a.leerTexto(etiqueta)
		numero, err := strconv.ParseFloat(valor, 64)
		if err != nil {
			fmt.Println("Ingrese un número decimal válido")
			continue
		}
		if err := validar(numero); err != nil {
			fmt.Printf("Valor inválido: %v\n", err)
			continue
		}
		return numero
	}
}

func validarNoVacio(mensaje string) func(string) error {
	return func(valor string) error {
		if strings.TrimSpace(valor) == "" {
			return errors.New(mensaje)
		}
		return nil
	}
}
