package web

import "net/http"

func (s *Servidor) listarClientes(w http.ResponseWriter, r *http.Request) {
	clientes, err := s.clientes.ListarClientes(r.Context())
	if err != nil {
		escribirError(w, estadoHTTPParaError(err), err)
		return
	}
	escribirJSON(w, http.StatusOK, clientesARespuesta(clientes))
}

func (s *Servidor) registrarCliente(w http.ResponseWriter, r *http.Request) {
	var solicitud registrarClienteSolicitud
	if err := decodificarJSON(w, r, &solicitud); err != nil {
		escribirError(w, http.StatusBadRequest, err)
		return
	}

	cliente, err := s.clientes.RegistrarCliente(r.Context(), solicitud.Nombre, solicitud.Email, solicitud.Telefono)
	if err != nil {
		escribirError(w, estadoHTTPParaError(err), err)
		return
	}

	escribirJSON(w, http.StatusCreated, clienteARespuesta(cliente))
}

func (s *Servidor) buscarCliente(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	cliente, err := s.clientes.BuscarCliente(r.Context(), id)
	if err != nil {
		escribirError(w, estadoHTTPParaError(err), err)
		return
	}
	escribirJSON(w, http.StatusOK, clienteARespuesta(cliente))
}

func (s *Servidor) eliminarCliente(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.clientes.EliminarCliente(r.Context(), id); err != nil {
		escribirError(w, estadoHTTPParaError(err), err)
		return
	}
	escribirJSON(w, http.StatusOK, respuestaMensaje{Mensaje: "cliente eliminado de forma lógica"})
}
