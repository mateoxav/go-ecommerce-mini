package clientes

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

func (s *Servicio) RegistrarCliente(ctx context.Context, nombre string, email string, telefono string) (modelos.Cliente, error) {
	cliente, err := modelos.NuevoCliente(nombre, email, telefono)
	if err != nil {
		return modelos.Cliente{}, fmt.Errorf("validar cliente: %w", err)
	}

	if err := s.repo.Crear(ctx, cliente); err != nil {
		return modelos.Cliente{}, err
	}

	return cliente, nil
}

func (s *Servicio) BuscarCliente(ctx context.Context, id string) (modelos.Cliente, error) {
	id = strings.TrimSpace(id)
	if !modelos.ValidarIDCliente(id) {
		return modelos.Cliente{}, modelos.ErrorIDClienteInvalido()
	}
	return s.repo.BuscarPorID(ctx, id)
}

func (s *Servicio) ListarClientes(ctx context.Context) ([]modelos.Cliente, error) {
	return s.repo.Listar(ctx)
}

func (s *Servicio) EliminarCliente(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if !modelos.ValidarIDCliente(id) {
		return modelos.ErrorIDClienteInvalido()
	}
	return s.repo.EliminarLogico(ctx, id)
}
