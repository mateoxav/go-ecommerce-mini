package web

import "net/http"

func (s *Servidor) listarPedidos(w http.ResponseWriter, r *http.Request) {
	filtroEstado := r.URL.Query().Get("estado")
	pedidos, err := s.pedidos.ListarPedidos(r.Context(), filtroEstado)
	if err != nil {
		escribirError(w, estadoHTTPParaError(err), err)
		return
	}
	escribirJSON(w, http.StatusOK, pedidosARespuesta(pedidos))
}

func (s *Servidor) crearPedido(w http.ResponseWriter, r *http.Request) {
	var solicitud crearPedidoSolicitud
	if err := decodificarJSON(w, r, &solicitud); err != nil {
		escribirError(w, http.StatusBadRequest, err)
		return
	}

	pedido, err := s.pedidos.CrearPedido(r.Context(), solicitud.ClienteID)
	if err != nil {
		escribirError(w, estadoHTTPParaError(err), err)
		return
	}

	escribirJSON(w, http.StatusCreated, pedidoARespuesta(pedido))
}

func (s *Servidor) buscarPedido(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	pedido, err := s.pedidos.BuscarPedido(r.Context(), id)
	if err != nil {
		escribirError(w, estadoHTTPParaError(err), err)
		return
	}
	escribirJSON(w, http.StatusOK, pedidoARespuesta(pedido))
}

func (s *Servidor) agregarItemPedido(w http.ResponseWriter, r *http.Request) {
	pedidoID := r.PathValue("id")
	var solicitud agregarItemSolicitud
	if err := decodificarJSON(w, r, &solicitud); err != nil {
		escribirError(w, http.StatusBadRequest, err)
		return
	}

	if err := s.pedidos.AgregarItem(r.Context(), pedidoID, solicitud.ProductoID, solicitud.Cantidad); err != nil {
		escribirError(w, estadoHTTPParaError(err), err)
		return
	}

	total, err := s.pedidos.CalcularTotal(r.Context(), pedidoID)
	if err != nil {
		escribirError(w, estadoHTTPParaError(err), err)
		return
	}

	escribirJSON(w, http.StatusOK, map[string]any{
		"mensaje":   "item agregado y stock actualizado",
		"pedido_id": pedidoID,
		"total":     total,
	})
}

func (s *Servidor) calcularTotalPedido(w http.ResponseWriter, r *http.Request) {
	pedidoID := r.PathValue("id")
	total, err := s.pedidos.CalcularTotal(r.Context(), pedidoID)
	if err != nil {
		escribirError(w, estadoHTTPParaError(err), err)
		return
	}
	escribirJSON(w, http.StatusOK, map[string]any{
		"pedido_id": pedidoID,
		"total":     total,
	})
}

func (s *Servidor) cambiarEstadoPedido(w http.ResponseWriter, r *http.Request) {
	pedidoID := r.PathValue("id")
	var solicitud cambiarEstadoSolicitud
	if err := decodificarJSON(w, r, &solicitud); err != nil {
		escribirError(w, http.StatusBadRequest, err)
		return
	}

	if err := s.pedidos.CambiarEstado(r.Context(), pedidoID, solicitud.Estado); err != nil {
		escribirError(w, estadoHTTPParaError(err), err)
		return
	}

	pedido, err := s.pedidos.BuscarPedido(r.Context(), pedidoID)
	if err != nil {
		escribirError(w, estadoHTTPParaError(err), err)
		return
	}

	escribirJSON(w, http.StatusOK, pedidoARespuesta(pedido))
}
