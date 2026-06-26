package web

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func (s *Servidor) verificarStock(w http.ResponseWriter, r *http.Request) {
	productoID := strings.TrimSpace(r.URL.Query().Get("id"))
	cantidadTexto := strings.TrimSpace(r.URL.Query().Get("cantidad"))
	cantidad, err := strconv.Atoi(cantidadTexto)
	if err != nil {
		escribirError(w, http.StatusBadRequest, fmt.Errorf("la cantidad debe ser un número entero válido"))
		return
	}

	disponible, err := s.inventario.VerificarStock(r.Context(), productoID, cantidad)
	if err != nil {
		escribirError(w, estadoHTTPParaError(err), err)
		return
	}

	producto, err := s.productos.BuscarProducto(r.Context(), productoID)
	if err != nil {
		escribirError(w, estadoHTTPParaError(err), err)
		return
	}

	escribirJSON(w, http.StatusOK, stockRespuesta{
		ProductoID:            productoID,
		CantidadSolicitada:    cantidad,
		StockDisponibleActual: producto.Stock(),
		Disponible:            disponible,
	})
}

func (s *Servidor) alertasStockBajo(w http.ResponseWriter, r *http.Request) {
	umbralTexto := strings.TrimSpace(r.URL.Query().Get("umbral"))
	if umbralTexto == "" {
		umbralTexto = "5"
	}
	umbral, err := strconv.Atoi(umbralTexto)
	if err != nil {
		escribirError(w, http.StatusBadRequest, fmt.Errorf("el umbral debe ser un número entero válido"))
		return
	}

	productos, err := s.inventario.AlertasStockBajo(r.Context(), umbral)
	if err != nil {
		escribirError(w, estadoHTTPParaError(err), err)
		return
	}

	escribirJSON(w, http.StatusOK, map[string]any{
		"umbral":    umbral,
		"productos": productosARespuesta(productos),
	})
}

func (s *Servidor) reponerStock(w http.ResponseWriter, r *http.Request) {
	var solicitud reponerStockSolicitud
	if err := decodificarJSON(w, r, &solicitud); err != nil {
		escribirError(w, http.StatusBadRequest, err)
		return
	}

	if err := s.inventario.ReponerStock(r.Context(), solicitud.ProductoID, solicitud.Cantidad); err != nil {
		escribirError(w, estadoHTTPParaError(err), err)
		return
	}

	producto, err := s.productos.BuscarProducto(r.Context(), solicitud.ProductoID)
	if err != nil {
		escribirError(w, estadoHTTPParaError(err), err)
		return
	}

	escribirJSON(w, http.StatusOK, productoARespuesta(producto))
}

func (s *Servidor) generarReporteInventario(w http.ResponseWriter, r *http.Request) {
	orden := strings.TrimSpace(r.URL.Query().Get("orden"))
	if orden == "" {
		orden = "nombre"
	}
	reporte, err := s.inventario.GenerarReporte(r.Context(), orden)
	if err != nil {
		escribirError(w, estadoHTTPParaError(err), err)
		return
	}

	escribirJSON(w, http.StatusOK, reporteInventarioRespuesta{
		Orden:   orden,
		Reporte: reporte,
	})
}
