package web

import (
	"net/http"
	"time"
)

type Servidor struct {
	productos  ServicioProductos
	clientes   ServicioClientes
	pedidos    ServicioPedidos
	inventario ServicioInventario
}

func NuevoServidor(
	productos ServicioProductos,
	clientes ServicioClientes,
	pedidos ServicioPedidos,
	inventario ServicioInventario,
) *Servidor {
	return &Servidor{
		productos:  productos,
		clientes:   clientes,
		pedidos:    pedidos,
		inventario: inventario,
	}
}

func (s *Servidor) Rutas() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/salud", s.salud)

	mux.HandleFunc("GET /api/productos", s.listarProductos)
	mux.HandleFunc("POST /api/productos", s.crearProducto)
	mux.HandleFunc("GET /api/productos/{id}", s.buscarProducto)
	mux.HandleFunc("PUT /api/productos/{id}/stock", s.actualizarStockProducto)
	mux.HandleFunc("DELETE /api/productos/{id}", s.eliminarProducto)

	mux.HandleFunc("GET /api/clientes", s.listarClientes)
	mux.HandleFunc("POST /api/clientes", s.registrarCliente)
	mux.HandleFunc("GET /api/clientes/{id}", s.buscarCliente)
	mux.HandleFunc("DELETE /api/clientes/{id}", s.eliminarCliente)

	mux.HandleFunc("GET /api/pedidos", s.listarPedidos)
	mux.HandleFunc("POST /api/pedidos", s.crearPedido)
	mux.HandleFunc("GET /api/pedidos/{id}", s.buscarPedido)
	mux.HandleFunc("POST /api/pedidos/{id}/items", s.agregarItemPedido)
	mux.HandleFunc("GET /api/pedidos/{id}/total", s.calcularTotalPedido)
	mux.HandleFunc("PUT /api/pedidos/{id}/estado", s.cambiarEstadoPedido)

	mux.HandleFunc("GET /api/inventario/stock", s.verificarStock)
	mux.HandleFunc("GET /api/inventario/stock-bajo", s.alertasStockBajo)
	mux.HandleFunc("POST /api/inventario/reponer", s.reponerStock)
	mux.HandleFunc("GET /api/inventario/reporte", s.generarReporteInventario)

	mux.HandleFunc("GET /api/concurrencia/resumen-inventario", s.resumenInventarioConcurrente)

	return conRegistro(conCORS(mux))
}

func NuevoHTTPServer(direccion string, manejador http.Handler) *http.Server {
	return &http.Server{
		Addr:              direccion,
		Handler:           manejador,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
}

func (s *Servidor) salud(w http.ResponseWriter, r *http.Request) {
	escribirJSON(w, http.StatusOK, map[string]string{
		"estado":   "ok",
		"servicio": "Sistema de Gestión de E-Commerce",
	})
}
