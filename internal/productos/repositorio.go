package productos

import (
	"context"

	"github.com/mateoxav/go-ecommerce-mini/internal/modelos"
)

type Repositorio interface {
	Crear(ctx context.Context, producto modelos.Producto) error
	BuscarPorID(ctx context.Context, id string) (modelos.Producto, error)
	ListarActivos(ctx context.Context) ([]modelos.Producto, error)
	ActualizarStock(ctx context.Context, id string, cambio int) error
	EliminarLogico(ctx context.Context, id string) error
}
