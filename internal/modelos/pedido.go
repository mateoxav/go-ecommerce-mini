package modelos

import (
	"errors"
	"strings"
	"time"
)

type EstadoPedido string

const (
	EstadoPendiente EstadoPedido = "pendiente"
	EstadoEnviado   EstadoPedido = "enviado"
	EstadoEntregado EstadoPedido = "entregado"
	EstadoCancelado EstadoPedido = "cancelado"
)

// Array usado para representar un conjunto fijo de estados válidos.
var estadosPedido = [4]EstadoPedido{EstadoPendiente, EstadoEnviado, EstadoEntregado, EstadoCancelado}

func EstadosPedidoValidos() []EstadoPedido {
	return estadosPedido[:]
}

func MapaEstadosPedidoValidos() map[string]bool {
	estados := make(map[string]bool, len(estadosPedido))
	for _, estado := range estadosPedido {
		estados[string(estado)] = true
	}
	return estados
}

func EstadoPedidoValido(estado string) bool {
	estado = strings.TrimSpace(strings.ToLower(estado))
	return MapaEstadosPedidoValidos()[estado]
}

type Pedido struct {
	id        string
	clienteID string
	total     float64
	estado    EstadoPedido
	fecha     string
}

type ItemPedido struct {
	id             int64
	pedidoID       string
	productoID     string
	cantidad       int
	precioUnitario float64
}

func NuevoPedido(clienteID string) (Pedido, error) {
	return ReconstruirPedido(generarID("PED"), clienteID, 0, string(EstadoPendiente), time.Now().Format(time.RFC3339))
}

func ReconstruirPedido(id string, clienteID string, total float64, estado string, fecha string) (Pedido, error) {
	id = strings.TrimSpace(id)
	clienteID = strings.TrimSpace(clienteID)
	estado = strings.TrimSpace(strings.ToLower(estado))
	fecha = strings.TrimSpace(fecha)

	if id == "" {
		return Pedido{}, errors.New("el id del pedido es obligatorio")
	}
	if !ValidarIDPedido(id) {
		return Pedido{}, ErrorIDPedidoInvalido()
	}
	if clienteID == "" {
		return Pedido{}, errors.New("el id del cliente es obligatorio")
	}
	if !ValidarIDCliente(clienteID) {
		return Pedido{}, ErrorIDClienteInvalido()
	}
	if total < 0 {
		return Pedido{}, errors.New("el total del pedido no puede ser negativo")
	}
	if !EstadoPedidoValido(estado) {
		return Pedido{}, errors.New("el estado del pedido no es válido")
	}
	if fecha == "" {
		return Pedido{}, errors.New("la fecha del pedido es obligatoria")
	}

	return Pedido{
		id:        id,
		clienteID: clienteID,
		total:     total,
		estado:    EstadoPedido(estado),
		fecha:     fecha,
	}, nil
}

func ReconstruirItemPedido(id int64, pedidoID string, productoID string, cantidad int, precioUnitario float64) (ItemPedido, error) {
	pedidoID = strings.TrimSpace(pedidoID)
	productoID = strings.TrimSpace(productoID)

	if id < 0 {
		return ItemPedido{}, errors.New("el id del item no puede ser negativo")
	}
	if pedidoID == "" {
		return ItemPedido{}, errors.New("el id del pedido es obligatorio")
	}
	if !ValidarIDPedido(pedidoID) {
		return ItemPedido{}, ErrorIDPedidoInvalido()
	}
	if productoID == "" {
		return ItemPedido{}, errors.New("el id del producto es obligatorio")
	}
	if !ValidarIDProducto(productoID) {
		return ItemPedido{}, ErrorIDProductoInvalido()
	}
	if cantidad <= 0 {
		return ItemPedido{}, errors.New("la cantidad debe ser mayor que cero")
	}
	if precioUnitario < 0 {
		return ItemPedido{}, errors.New("el precio unitario no puede ser negativo")
	}

	return ItemPedido{
		id:             id,
		pedidoID:       pedidoID,
		productoID:     productoID,
		cantidad:       cantidad,
		precioUnitario: precioUnitario,
	}, nil
}

func (p Pedido) ID() string        { return p.id }
func (p Pedido) ClienteID() string { return p.clienteID }
func (p Pedido) Total() float64    { return p.total }
func (p Pedido) Estado() string    { return string(p.estado) }
func (p Pedido) Fecha() string     { return p.fecha }

func (i ItemPedido) ID() int64               { return i.id }
func (i ItemPedido) PedidoID() string        { return i.pedidoID }
func (i ItemPedido) ProductoID() string      { return i.productoID }
func (i ItemPedido) Cantidad() int           { return i.cantidad }
func (i ItemPedido) PrecioUnitario() float64 { return i.precioUnitario }
func (i ItemPedido) Subtotal() float64       { return float64(i.cantidad) * i.precioUnitario }
