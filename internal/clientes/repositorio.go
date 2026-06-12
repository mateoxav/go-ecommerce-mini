package clientes

import (
	"context"

	"github.com/mateoxav/go-ecommerce-mini/internal/modelos"
)

type Repositorio interface {
	Crear(ctx context.Context, cliente modelos.Cliente) error
	BuscarPorID(ctx context.Context, id string) (modelos.Cliente, error)
	Listar(ctx context.Context) ([]modelos.Cliente, error)
}
