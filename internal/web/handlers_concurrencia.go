package web

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/mateoxav/go-ecommerce-mini/internal/modelos"
)

type resultadoConcurrente struct {
	clave string
	valor any
	err   error
}

func (s *Servidor) resumenInventarioConcurrente(w http.ResponseWriter, r *http.Request) {
	inicio := time.Now()
	ctx := r.Context()

	umbralTexto := strings.TrimSpace(r.URL.Query().Get("umbral"))
	if umbralTexto == "" {
		umbralTexto = "5"
	}
	umbral, err := strconv.Atoi(umbralTexto)
	if err != nil {
		escribirError(w, http.StatusBadRequest, fmt.Errorf("el umbral debe ser un número entero válido"))
		return
	}

	orden := strings.TrimSpace(r.URL.Query().Get("orden"))
	if orden == "" {
		orden = "nombre"
	}

	canalResultados := make(chan resultadoConcurrente, 3)
	var grupo sync.WaitGroup
	grupo.Add(3)

	go func() {
		defer grupo.Done()
		productos, err := s.productos.ListarProductos(ctx)
		canalResultados <- resultadoConcurrente{clave: "productos", valor: productos, err: err}
	}()

	go func() {
		defer grupo.Done()
		productos, err := s.inventario.AlertasStockBajo(ctx, umbral)
		canalResultados <- resultadoConcurrente{clave: "stock_bajo", valor: productos, err: err}
	}()

	go func() {
		defer grupo.Done()
		reporte, err := s.inventario.GenerarReporte(ctx, orden)
		canalResultados <- resultadoConcurrente{clave: "reporte", valor: reporte, err: err}
	}()

	go func() {
		grupo.Wait()
		close(canalResultados)
	}()

	var (
		productos          []ProductoRespuesta
		productosStockBajo []ProductoRespuesta
		reporte            string
		errores            []string
	)

	for resultado := range canalResultados {
		if resultado.err != nil {
			errores = append(errores, fmt.Sprintf("%s: %v", resultado.clave, resultado.err))
			continue
		}

		switch resultado.clave {
		case "productos":
			lista, ok := resultado.valor.([]modelos.Producto)
			if ok {
				productos = productosARespuesta(lista)
			}
		case "stock_bajo":
			lista, ok := resultado.valor.([]modelos.Producto)
			if ok {
				productosStockBajo = productosARespuesta(lista)
			}
		case "reporte":
			texto, ok := resultado.valor.(string)
			if ok {
				reporte = texto
			}
		}
	}

	respuesta := resumenConcurrenteRespuesta{
		Mensaje:            "resumen generado con goroutines y canales",
		TotalProductos:     len(productos),
		ProductosStockBajo: productosStockBajo,
		Reporte:            reporte,
		Errores:            errores,
		TiempoMS:           time.Since(inicio).Milliseconds(),
		Concurrencia:       concurrenciaRespuesta{GoroutinesEjecutadas: 3, CanalUsado: "canalResultados"},
	}

	if len(errores) > 0 {
		respuesta.Mensaje = "resumen concurrente incompleto"
		escribirJSON(w, http.StatusInternalServerError, respuesta)
		return
	}

	escribirJSON(w, http.StatusOK, respuesta)
}
