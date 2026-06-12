package pedidos

import (
	"context"
	"fmt"
	"strings"

	"github.com/mateoxav/go-ecommerce-mini/internal/modelos"
)

type Servicio struct {
	repo         Repositorio
	repoClientes RepositorioClientesLectura
}

func NuevoServicio(repo Repositorio, repoClientes RepositorioClientesLectura) *Servicio {
	return &Servicio{repo: repo, repoClientes: repoClientes}
}

func (s *Servicio) CrearPedido(ctx context.Context, clienteID string) (modelos.Pedido, error) {
	if _, err := s.repoClientes.BuscarPorID(ctx, clienteID); err != nil {
		return modelos.Pedido{}, fmt.Errorf("validar cliente del pedido: %w", err)
	}

	pedido, err := modelos.NuevoPedido(clienteID)
	if err != nil {
		return modelos.Pedido{}, fmt.Errorf("validar pedido: %w", err)
	}

	if err := s.repo.Crear(ctx, pedido); err != nil {
		return modelos.Pedido{}, err
	}

	return pedido, nil
}

func (s *Servicio) AgregarItem(ctx context.Context, pedidoID string, productoID string, cantidad int) error {
	if strings.TrimSpace(pedidoID) == "" || strings.TrimSpace(productoID) == "" {
		return fmt.Errorf("pedidoID y productoID son obligatorios")
	}
	if cantidad <= 0 {
		return fmt.Errorf("la cantidad debe ser mayor que cero")
	}
	return s.repo.AgregarItem(ctx, pedidoID, productoID, cantidad)
}

func (s *Servicio) CalcularTotal(ctx context.Context, pedidoID string) (float64, error) {
	return s.repo.CalcularTotal(ctx, pedidoID)
}

func (s *Servicio) CambiarEstado(ctx context.Context, pedidoID string, estado string) error {
	estado = strings.TrimSpace(strings.ToLower(estado))
	if !modelos.EstadoPedidoValido(estado) {
		return fmt.Errorf("estado inválido. Estados permitidos: pendiente, enviado, entregado, cancelado")
	}
	return s.repo.CambiarEstado(ctx, pedidoID, estado)
}

func (s *Servicio) ListarPedidos(ctx context.Context, filtroEstado string) ([]modelos.Pedido, error) {
	filtroEstado = strings.TrimSpace(strings.ToLower(filtroEstado))
	if filtroEstado != "" && !modelos.EstadoPedidoValido(filtroEstado) {
		return nil, fmt.Errorf("filtro de estado inválido")
	}
	return s.repo.Listar(ctx, filtroEstado)
}

func (s *Servicio) BuscarPedido(ctx context.Context, id string) (modelos.Pedido, error) {
	return s.repo.BuscarPorID(ctx, id)
}
