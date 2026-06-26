package web

import (
	"context"

	"github.com/mateoxav/go-ecommerce-mini/internal/modelos"
)

type ServicioProductos interface {
	CrearProducto(ctx context.Context, nombre string, precio float64, stock int, categoria string) (modelos.Producto, error)
	BuscarProducto(ctx context.Context, id string) (modelos.Producto, error)
	ListarProductos(ctx context.Context) ([]modelos.Producto, error)
	ActualizarStock(ctx context.Context, id string, cambio int) error
	EliminarProducto(ctx context.Context, id string) error
}

type ServicioClientes interface {
	RegistrarCliente(ctx context.Context, nombre string, email string, telefono string) (modelos.Cliente, error)
	BuscarCliente(ctx context.Context, id string) (modelos.Cliente, error)
	ListarClientes(ctx context.Context) ([]modelos.Cliente, error)
	EliminarCliente(ctx context.Context, id string) error
}

type ServicioPedidos interface {
	CrearPedido(ctx context.Context, clienteID string) (modelos.Pedido, error)
	AgregarItem(ctx context.Context, pedidoID string, productoID string, cantidad int) error
	CalcularTotal(ctx context.Context, pedidoID string) (float64, error)
	CambiarEstado(ctx context.Context, pedidoID string, estado string) error
	ListarPedidos(ctx context.Context, filtroEstado string) ([]modelos.Pedido, error)
	BuscarPedido(ctx context.Context, id string) (modelos.Pedido, error)
}

type ServicioInventario interface {
	VerificarStock(ctx context.Context, productoID string, cantidad int) (bool, error)
	AlertasStockBajo(ctx context.Context, umbral int) ([]modelos.Producto, error)
	ReponerStock(ctx context.Context, productoID string, cantidad int) error
	GenerarReporte(ctx context.Context, ordenarPor string) (string, error)
}
