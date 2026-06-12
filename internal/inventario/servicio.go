package inventario

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/mateoxav/go-ecommerce-mini/internal/modelos"
)

type RepositorioProductosLecturaEscritura interface {
	BuscarPorID(ctx context.Context, id string) (modelos.Producto, error)
	ListarActivos(ctx context.Context) ([]modelos.Producto, error)
	ActualizarStock(ctx context.Context, id string, cambio int) error
}

type Servicio struct {
	repoProductos RepositorioProductosLecturaEscritura
}

func NuevoServicio(repoProductos RepositorioProductosLecturaEscritura) *Servicio {
	return &Servicio{repoProductos: repoProductos}
}

func (s *Servicio) VerificarStock(ctx context.Context, productoID string, cantidad int) (bool, error) {
	productoID = strings.TrimSpace(productoID)
	if !modelos.ValidarIDProducto(productoID) {
		return false, modelos.ErrorIDProductoInvalido()
	}
	if cantidad <= 0 {
		return false, fmt.Errorf("la cantidad debe ser mayor que cero")
	}

	producto, err := s.repoProductos.BuscarPorID(ctx, productoID)
	if err != nil {
		return false, err
	}

	return producto.Stock() >= cantidad, nil
}

func (s *Servicio) AlertasStockBajo(ctx context.Context, umbral int) ([]modelos.Producto, error) {
	if umbral < 0 {
		return nil, fmt.Errorf("el umbral no puede ser negativo")
	}

	productos, err := s.repoProductos.ListarActivos(ctx)
	if err != nil {
		return nil, err
	}

	alertas := make([]modelos.Producto, 0)
	for _, producto := range productos {
		if producto.Stock() < umbral {
			alertas = append(alertas, producto)
		}
	}

	return alertas, nil
}

func (s *Servicio) ReponerStock(ctx context.Context, productoID string, cantidad int) error {
	productoID = strings.TrimSpace(productoID)
	if !modelos.ValidarIDProducto(productoID) {
		return modelos.ErrorIDProductoInvalido()
	}
	if cantidad <= 0 {
		return fmt.Errorf("la cantidad a reponer debe ser mayor que cero")
	}
	return s.repoProductos.ActualizarStock(ctx, productoID, cantidad)
}

func (s *Servicio) GenerarReporte(ctx context.Context, ordenarPor string) (string, error) {
	productos, err := s.repoProductos.ListarActivos(ctx)
	if err != nil {
		return "", err
	}

	ordenarPor = strings.TrimSpace(strings.ToLower(ordenarPor))
	if ordenarPor == "" {
		ordenarPor = "nombre"
	}

	ordenamientos := map[string]func(i, j int) bool{
		"nombre": func(i, j int) bool { return productos[i].Nombre() < productos[j].Nombre() },
		"precio": func(i, j int) bool { return productos[i].Precio() < productos[j].Precio() },
		"stock":  func(i, j int) bool { return productos[i].Stock() < productos[j].Stock() },
	}

	criterio, ok := ordenamientos[ordenarPor]
	if !ok {
		return "", fmt.Errorf("criterio de ordenamiento inválido: use nombre, precio o stock")
	}

	sort.Slice(productos, criterio)

	var reporte strings.Builder
	reporte.WriteString("\nREPORTE DE INVENTARIO\n")
	reporte.WriteString("------------------------------------------------------------\n")
	reporte.WriteString(fmt.Sprintf("%-22s %-12s %-10s %-8s\n", "Producto", "Categoría", "Precio", "Stock"))
	reporte.WriteString("------------------------------------------------------------\n")

	for _, producto := range productos {
		reporte.WriteString(fmt.Sprintf("%-22s %-12s $%-9.2f %-8d\n",
			producto.Nombre(), producto.Categoria(), producto.Precio(), producto.Stock()))
	}

	return reporte.String(), nil
}
