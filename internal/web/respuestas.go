package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func escribirJSON(w http.ResponseWriter, estado int, datos any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(estado)
	if datos == nil {
		return
	}
	_ = json.NewEncoder(w).Encode(datos)
}

func escribirError(w http.ResponseWriter, estado int, err error) {
	mensaje := "error interno del servidor"
	if err != nil {
		mensaje = err.Error()
	}
	escribirJSON(w, estado, respuestaError{Error: mensaje})
}

func decodificarJSON(w http.ResponseWriter, r *http.Request, destino any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	decodificador := json.NewDecoder(r.Body)
	decodificador.DisallowUnknownFields()

	if err := decodificador.Decode(destino); err != nil {
		return fmt.Errorf("cuerpo JSON inválido: %w", err)
	}
	if err := decodificador.Decode(&struct{}{}); err != io.EOF {
		return errors.New("el cuerpo JSON contiene más de un objeto")
	}
	return nil
}

func estadoHTTPParaError(err error) int {
	if err == nil {
		return http.StatusOK
	}
	mensaje := strings.ToLower(err.Error())

	if strings.Contains(mensaje, "no encontrado") || strings.Contains(mensaje, "inexistente") {
		return http.StatusNotFound
	}

	erroresCliente := []string{
		"inválido",
		"invalido",
		"obligatorio",
		"negativo",
		"mayor que cero",
		"stock insuficiente",
		"no puede ser cero",
		"no tiene un formato válido",
		"debe iniciar",
		"debe contener",
		"criterio de ordenamiento",
		"filtro de estado",
	}
	for _, fragmento := range erroresCliente {
		if strings.Contains(mensaje, fragmento) {
			return http.StatusBadRequest
		}
	}

	return http.StatusInternalServerError
}
