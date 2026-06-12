package pedidos

import (
	"context"

	"github.com/mateoxav/go-ecommerce-mini/internal/modelos"
)

type Repositorio interface {
	Crear(ctx context.Context, pedido modelos.Pedido) error
	BuscarPorID(ctx context.Context, id string) (modelos.Pedido, error)
	AgregarItem(ctx context.Context, pedidoID string, productoID string, cantidad int) error
	CalcularTotal(ctx context.Context, pedidoID string) (float64, error)
	CambiarEstado(ctx context.Context, pedidoID string, estado string) error
	Listar(ctx context.Context, filtroEstado string) ([]modelos.Pedido, error)
}

type RepositorioClientesLectura interface {
	BuscarPorID(ctx context.Context, id string) (modelos.Cliente, error)
}
