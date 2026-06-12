package productos

import (
	"context"
	"fmt"
	"strings"

	"github.com/mateoxav/go-ecommerce-mini/internal/modelos"
)

type Servicio struct {
	repo Repositorio
}

func NuevoServicio(repo Repositorio) *Servicio {
	return &Servicio{repo: repo}
}

func (s *Servicio) CrearProducto(ctx context.Context, nombre string, precio float64, stock int, categoria string) (modelos.Producto, error) {
	producto, err := modelos.NuevoProducto(nombre, precio, stock, categoria)
	if err != nil {
		return modelos.Producto{}, fmt.Errorf("validar producto: %w", err)
	}

	if err := s.repo.Crear(ctx, producto); err != nil {
		return modelos.Producto{}, err
	}

	return producto, nil
}

func (s *Servicio) BuscarProducto(ctx context.Context, id string) (modelos.Producto, error) {
	id = strings.TrimSpace(id)
	if !modelos.ValidarIDProducto(id) {
		return modelos.Producto{}, modelos.ErrorIDProductoInvalido()
	}
	return s.repo.BuscarPorID(ctx, id)
}

func (s *Servicio) ListarProductos(ctx context.Context) ([]modelos.Producto, error) {
	return s.repo.ListarActivos(ctx)
}

func (s *Servicio) ActualizarStock(ctx context.Context, id string, cambio int) error {
	id = strings.TrimSpace(id)
	if !modelos.ValidarIDProducto(id) {
		return modelos.ErrorIDProductoInvalido()
	}
	if cambio == 0 {
		return fmt.Errorf("el cambio de stock no puede ser cero")
	}
	return s.repo.ActualizarStock(ctx, id, cambio)
}

func (s *Servicio) EliminarProducto(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if !modelos.ValidarIDProducto(id) {
		return modelos.ErrorIDProductoInvalido()
	}
	return s.repo.EliminarLogico(ctx, id)
}
