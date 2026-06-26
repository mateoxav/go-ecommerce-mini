package web

import "net/http"

func (s *Servidor) listarProductos(w http.ResponseWriter, r *http.Request) {
	productos, err := s.productos.ListarProductos(r.Context())
	if err != nil {
		escribirError(w, estadoHTTPParaError(err), err)
		return
	}
	escribirJSON(w, http.StatusOK, productosARespuesta(productos))
}

func (s *Servidor) crearProducto(w http.ResponseWriter, r *http.Request) {
	var solicitud crearProductoSolicitud
	if err := decodificarJSON(w, r, &solicitud); err != nil {
		escribirError(w, http.StatusBadRequest, err)
		return
	}

	producto, err := s.productos.CrearProducto(r.Context(), solicitud.Nombre, solicitud.Precio, solicitud.Stock, solicitud.Categoria)
	if err != nil {
		escribirError(w, estadoHTTPParaError(err), err)
		return
	}

	escribirJSON(w, http.StatusCreated, productoARespuesta(producto))
}

func (s *Servidor) buscarProducto(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	producto, err := s.productos.BuscarProducto(r.Context(), id)
	if err != nil {
		escribirError(w, estadoHTTPParaError(err), err)
		return
	}
	escribirJSON(w, http.StatusOK, productoARespuesta(producto))
}

func (s *Servidor) actualizarStockProducto(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var solicitud actualizarStockSolicitud
	if err := decodificarJSON(w, r, &solicitud); err != nil {
		escribirError(w, http.StatusBadRequest, err)
		return
	}

	if err := s.productos.ActualizarStock(r.Context(), id, solicitud.Cambio); err != nil {
		escribirError(w, estadoHTTPParaError(err), err)
		return
	}

	producto, err := s.productos.BuscarProducto(r.Context(), id)
	if err != nil {
		escribirError(w, estadoHTTPParaError(err), err)
		return
	}

	escribirJSON(w, http.StatusOK, productoARespuesta(producto))
}

func (s *Servidor) eliminarProducto(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.productos.EliminarProducto(r.Context(), id); err != nil {
		escribirError(w, estadoHTTPParaError(err), err)
		return
	}
	escribirJSON(w, http.StatusOK, respuestaMensaje{Mensaje: "producto eliminado de forma lógica"})
}
