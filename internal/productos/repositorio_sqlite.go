package productos

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mateoxav/go-ecommerce-mini/internal/modelos"
)

type RepositorioSQLite struct {
	db *sql.DB
}

func NuevoRepositorioSQLite(db *sql.DB) *RepositorioSQLite {
	return &RepositorioSQLite{db: db}
}

func (r *RepositorioSQLite) Crear(ctx context.Context, producto modelos.Producto) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO productos (id, nombre, precio, stock, categoria, activo)
		VALUES (?, ?, ?, ?, ?, ?)
	`, producto.ID(), producto.Nombre(), producto.Precio(), producto.Stock(), producto.Categoria(), producto.ActivoEntero())
	if err != nil {
		return fmt.Errorf("crear producto: %w", err)
	}
	return nil
}

func (r *RepositorioSQLite) BuscarPorID(ctx context.Context, id string) (modelos.Producto, error) {
	var (
		productoID string
		nombre     string
		precio     float64
		stock      int
		categoria  string
		activo     int
	)

	err := r.db.QueryRowContext(ctx, `
		SELECT id, nombre, precio, stock, categoria, activo
		FROM productos
		WHERE id = ? AND activo = 1
	`, id).Scan(&productoID, &nombre, &precio, &stock, &categoria, &activo)
	if err != nil {
		if err == sql.ErrNoRows {
			return modelos.Producto{}, fmt.Errorf("producto no encontrado: %w", err)
		}
		return modelos.Producto{}, fmt.Errorf("buscar producto: %w", err)
	}

	producto, err := modelos.ReconstruirProducto(productoID, nombre, precio, stock, categoria, modelos.ActivoDesdeEntero(activo))
	if err != nil {
		return modelos.Producto{}, fmt.Errorf("reconstruir producto: %w", err)
	}

	return producto, nil
}

func (r *RepositorioSQLite) ListarActivos(ctx context.Context) ([]modelos.Producto, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, nombre, precio, stock, categoria, activo
		FROM productos
		WHERE activo = 1
		ORDER BY nombre ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("listar productos: %w", err)
	}
	defer rows.Close()

	productos := make([]modelos.Producto, 0)
	for rows.Next() {
		var (
			id        string
			nombre    string
			precio    float64
			stock     int
			categoria string
			activo    int
		)
		if err := rows.Scan(&id, &nombre, &precio, &stock, &categoria, &activo); err != nil {
			return nil, fmt.Errorf("leer producto: %w", err)
		}
		producto, err := modelos.ReconstruirProducto(id, nombre, precio, stock, categoria, modelos.ActivoDesdeEntero(activo))
		if err != nil {
			return nil, fmt.Errorf("reconstruir producto: %w", err)
		}
		productos = append(productos, producto)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("recorrer productos: %w", err)
	}

	return productos, nil
}

func (r *RepositorioSQLite) ActualizarStock(ctx context.Context, id string, cambio int) error {
	resultado, err := r.db.ExecContext(ctx, `
		UPDATE productos
		SET stock = stock + ?
		WHERE id = ? AND activo = 1 AND stock + ? >= 0
	`, cambio, id, cambio)
	if err != nil {
		return fmt.Errorf("actualizar stock: %w", err)
	}

	filas, err := resultado.RowsAffected()
	if err != nil {
		return fmt.Errorf("verificar actualización de stock: %w", err)
	}
	if filas == 0 {
		return fmt.Errorf("no se pudo actualizar stock: producto inexistente o stock insuficiente")
	}

	return nil
}

func (r *RepositorioSQLite) EliminarLogico(ctx context.Context, id string) error {
	resultado, err := r.db.ExecContext(ctx, `
		UPDATE productos
		SET activo = 0
		WHERE id = ? AND activo = 1
	`, id)
	if err != nil {
		return fmt.Errorf("eliminar producto: %w", err)
	}

	filas, err := resultado.RowsAffected()
	if err != nil {
		return fmt.Errorf("verificar eliminación de producto: %w", err)
	}
	if filas == 0 {
		return fmt.Errorf("producto no encontrado o ya eliminado")
	}

	return nil
}
