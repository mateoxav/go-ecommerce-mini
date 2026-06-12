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
	clienteID = strings.TrimSpace(clienteID)
	if !modelos.ValidarIDCliente(clienteID) {
		return modelos.Pedido{}, modelos.ErrorIDClienteInvalido()
	}
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
	pedidoID = strings.TrimSpace(pedidoID)
	productoID = strings.TrimSpace(productoID)
	if !modelos.ValidarIDPedido(pedidoID) {
		return modelos.ErrorIDPedidoInvalido()
	}
	if !modelos.ValidarIDProducto(productoID) {
		return modelos.ErrorIDProductoInvalido()
	}
	if cantidad <= 0 {
		return fmt.Errorf("la cantidad debe ser mayor que cero")
	}
	return s.repo.AgregarItem(ctx, pedidoID, productoID, cantidad)
}

func (s *Servicio) CalcularTotal(ctx context.Context, pedidoID string) (float64, error) {
	pedidoID = strings.TrimSpace(pedidoID)
	if !modelos.ValidarIDPedido(pedidoID) {
		return 0, modelos.ErrorIDPedidoInvalido()
	}
	if _, err := s.repo.BuscarPorID(ctx, pedidoID); err != nil {
		return 0, err
	}
	return s.repo.CalcularTotal(ctx, pedidoID)
}

func (s *Servicio) CambiarEstado(ctx context.Context, pedidoID string, estado string) error {
	pedidoID = strings.TrimSpace(pedidoID)
	if !modelos.ValidarIDPedido(pedidoID) {
		return modelos.ErrorIDPedidoInvalido()
	}
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
	id = strings.TrimSpace(id)
	if !modelos.ValidarIDPedido(id) {
		return modelos.Pedido{}, modelos.ErrorIDPedidoInvalido()
	}
	return s.repo.BuscarPorID(ctx, id)
}
