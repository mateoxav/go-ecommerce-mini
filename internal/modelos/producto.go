package modelos

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type Producto struct {
	id        string
	nombre    string
	precio    float64
	stock     int
	categoria string
	activo    bool
}

func NuevoProducto(nombre string, precio float64, stock int, categoria string) (Producto, error) {
	return ReconstruirProducto(generarID("PROD"), nombre, precio, stock, categoria, true)
}

func ReconstruirProducto(id string, nombre string, precio float64, stock int, categoria string, activo bool) (Producto, error) {
	id = strings.TrimSpace(id)
	nombre = strings.TrimSpace(nombre)
	categoria = strings.TrimSpace(categoria)

	if id == "" {
		return Producto{}, errors.New("el id del producto es obligatorio")
	}
	if !ValidarIDProducto(id) {
		return Producto{}, ErrorIDProductoInvalido()
	}
	if nombre == "" {
		return Producto{}, errors.New("el nombre del producto es obligatorio")
	}
	if precio < 0 {
		return Producto{}, errors.New("el precio del producto no puede ser negativo")
	}
	if stock < 0 {
		return Producto{}, errors.New("el stock del producto no puede ser negativo")
	}
	if categoria == "" {
		return Producto{}, errors.New("la categoría del producto es obligatoria")
	}

	return Producto{
		id:        id,
		nombre:    nombre,
		precio:    precio,
		stock:     stock,
		categoria: categoria,
		activo:    activo,
	}, nil
}

func (p Producto) ID() string        { return p.id }
func (p Producto) Nombre() string    { return p.nombre }
func (p Producto) Precio() float64   { return p.precio }
func (p Producto) Stock() int        { return p.stock }
func (p Producto) Categoria() string { return p.categoria }
func (p Producto) Activo() bool      { return p.activo }

func (p Producto) ActivoEntero() int {
	if p.activo {
		return 1
	}
	return 0
}

func (p Producto) ConStock(nuevoStock int) (Producto, error) {
	if nuevoStock < 0 {
		return Producto{}, errors.New("el stock resultante no puede ser negativo")
	}
	p.stock = nuevoStock
	return p, nil
}

func ActivoDesdeEntero(valor int) bool {
	return valor == 1
}

func generarID(prefijo string) string {
	return fmt.Sprintf("%s-%d", prefijo, time.Now().UnixNano())
}
