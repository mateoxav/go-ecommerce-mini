package web

import "github.com/mateoxav/go-ecommerce-mini/internal/modelos"

type respuestaError struct {
	Error string `json:"error"`
}

type respuestaMensaje struct {
	Mensaje string `json:"mensaje"`
}

type ProductoRespuesta struct {
	ID        string  `json:"id"`
	Nombre    string  `json:"nombre"`
	Precio    float64 `json:"precio"`
	Stock     int     `json:"stock"`
	Categoria string  `json:"categoria"`
	Activo    bool    `json:"activo"`
}

type ClienteRespuesta struct {
	ID            string `json:"id"`
	Nombre        string `json:"nombre"`
	Email         string `json:"email"`
	Telefono      string `json:"telefono"`
	FechaRegistro string `json:"fecha_registro"`
}

type PedidoRespuesta struct {
	ID        string  `json:"id"`
	ClienteID string  `json:"cliente_id"`
	Total     float64 `json:"total"`
	Estado    string  `json:"estado"`
	Fecha     string  `json:"fecha"`
}

type crearProductoSolicitud struct {
	Nombre    string  `json:"nombre"`
	Precio    float64 `json:"precio"`
	Stock     int     `json:"stock"`
	Categoria string  `json:"categoria"`
}

type actualizarStockSolicitud struct {
	Cambio int `json:"cambio"`
}

type registrarClienteSolicitud struct {
	Nombre   string `json:"nombre"`
	Email    string `json:"email"`
	Telefono string `json:"telefono"`
}

type crearPedidoSolicitud struct {
	ClienteID string `json:"cliente_id"`
}

type agregarItemSolicitud struct {
	ProductoID string `json:"producto_id"`
	Cantidad   int    `json:"cantidad"`
}

type cambiarEstadoSolicitud struct {
	Estado string `json:"estado"`
}

type reponerStockSolicitud struct {
	ProductoID string `json:"producto_id"`
	Cantidad   int    `json:"cantidad"`
}

type stockRespuesta struct {
	ProductoID            string `json:"producto_id"`
	CantidadSolicitada    int    `json:"cantidad_solicitada"`
	StockDisponibleActual int    `json:"stock_disponible_actual"`
	Disponible            bool   `json:"disponible"`
}

type reporteInventarioRespuesta struct {
	Orden   string `json:"orden"`
	Reporte string `json:"reporte"`
}

type resumenConcurrenteRespuesta struct {
	Mensaje            string                `json:"mensaje"`
	TotalProductos     int                   `json:"total_productos"`
	ProductosStockBajo []ProductoRespuesta   `json:"productos_stock_bajo"`
	Reporte            string                `json:"reporte"`
	Errores            []string              `json:"errores,omitempty"`
	TiempoMS           int64                 `json:"tiempo_ms"`
	Concurrencia       concurrenciaRespuesta `json:"concurrencia"`
}

type concurrenciaRespuesta struct {
	GoroutinesEjecutadas int    `json:"goroutines_ejecutadas"`
	CanalUsado           string `json:"canal_usado"`
}

func productoARespuesta(producto modelos.Producto) ProductoRespuesta {
	return ProductoRespuesta{
		ID:        producto.ID(),
		Nombre:    producto.Nombre(),
		Precio:    producto.Precio(),
		Stock:     producto.Stock(),
		Categoria: producto.Categoria(),
		Activo:    producto.Activo(),
	}
}

func productosARespuesta(productos []modelos.Producto) []ProductoRespuesta {
	respuesta := make([]ProductoRespuesta, 0, len(productos))
	for _, producto := range productos {
		respuesta = append(respuesta, productoARespuesta(producto))
	}
	return respuesta
}

func clienteARespuesta(cliente modelos.Cliente) ClienteRespuesta {
	return ClienteRespuesta{
		ID:            cliente.ID(),
		Nombre:        cliente.Nombre(),
		Email:         cliente.Email(),
		Telefono:      cliente.Telefono(),
		FechaRegistro: cliente.FechaRegistro(),
	}
}

func clientesARespuesta(clientes []modelos.Cliente) []ClienteRespuesta {
	respuesta := make([]ClienteRespuesta, 0, len(clientes))
	for _, cliente := range clientes {
		respuesta = append(respuesta, clienteARespuesta(cliente))
	}
	return respuesta
}

func pedidoARespuesta(pedido modelos.Pedido) PedidoRespuesta {
	return PedidoRespuesta{
		ID:        pedido.ID(),
		ClienteID: pedido.ClienteID(),
		Total:     pedido.Total(),
		Estado:    pedido.Estado(),
		Fecha:     pedido.Fecha(),
	}
}

func pedidosARespuesta(pedidos []modelos.Pedido) []PedidoRespuesta {
	respuesta := make([]PedidoRespuesta, 0, len(pedidos))
	for _, pedido := range pedidos {
		respuesta = append(respuesta, pedidoARespuesta(pedido))
	}
	return respuesta
}
