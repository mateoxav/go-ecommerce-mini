package pedidos

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

func (r *RepositorioSQLite) Crear(ctx context.Context, pedido modelos.Pedido) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO pedidos (id, cliente_id, total, estado, fecha)
		VALUES (?, ?, ?, ?, ?)
	`, pedido.ID(), pedido.ClienteID(), pedido.Total(), pedido.Estado(), pedido.Fecha())
	if err != nil {
		return fmt.Errorf("crear pedido: %w", err)
	}
	return nil
}

func (r *RepositorioSQLite) BuscarPorID(ctx context.Context, id string) (modelos.Pedido, error) {
	var (
		pedidoID  string
		clienteID string
		total     float64
		estado    string
		fecha     string
	)

	err := r.db.QueryRowContext(ctx, `
		SELECT id, cliente_id, total, estado, fecha
		FROM pedidos
		WHERE id = ?
	`, id).Scan(&pedidoID, &clienteID, &total, &estado, &fecha)
	if err != nil {
		if err == sql.ErrNoRows {
			return modelos.Pedido{}, fmt.Errorf("pedido no encontrado: %w", err)
		}
		return modelos.Pedido{}, fmt.Errorf("buscar pedido: %w", err)
	}

	pedido, err := modelos.ReconstruirPedido(pedidoID, clienteID, total, estado, fecha)
	if err != nil {
		return modelos.Pedido{}, fmt.Errorf("reconstruir pedido: %w", err)
	}

	return pedido, nil
}

func (r *RepositorioSQLite) AgregarItem(ctx context.Context, pedidoID string, productoID string, cantidad int) (err error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("iniciar transacción: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	var precio float64
	var stock int
	err = tx.QueryRowContext(ctx, `
		SELECT precio, stock
		FROM productos
		WHERE id = ? AND activo = 1
	`, productoID).Scan(&precio, &stock)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("producto no encontrado: %w", err)
		}
		return fmt.Errorf("consultar producto para pedido: %w", err)
	}

	if cantidad <= 0 {
		return fmt.Errorf("la cantidad debe ser mayor que cero")
	}
	if stock < cantidad {
		return fmt.Errorf("stock insuficiente: disponible %d, solicitado %d", stock, cantidad)
	}

	var existePedido int
	if err = tx.QueryRowContext(ctx, `SELECT COUNT(1) FROM pedidos WHERE id = ?`, pedidoID).Scan(&existePedido); err != nil {
		return fmt.Errorf("verificar pedido: %w", err)
	}
	if existePedido == 0 {
		return fmt.Errorf("pedido no encontrado")
	}

	if _, err = tx.ExecContext(ctx, `
		INSERT INTO items_pedido (pedido_id, producto_id, cantidad, precio_unitario)
		VALUES (?, ?, ?, ?)
	`, pedidoID, productoID, cantidad, precio); err != nil {
		return fmt.Errorf("agregar item al pedido: %w", err)
	}

	if _, err = tx.ExecContext(ctx, `
		UPDATE productos
		SET stock = stock - ?
		WHERE id = ? AND stock - ? >= 0
	`, cantidad, productoID, cantidad); err != nil {
		return fmt.Errorf("descontar stock: %w", err)
	}

	if _, err = tx.ExecContext(ctx, `
		UPDATE pedidos
		SET total = (
			SELECT COALESCE(SUM(cantidad * precio_unitario), 0)
			FROM items_pedido
			WHERE pedido_id = ?
		)
		WHERE id = ?
	`, pedidoID, pedidoID); err != nil {
		return fmt.Errorf("recalcular total del pedido: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("confirmar transacción: %w", err)
	}

	return nil
}

func (r *RepositorioSQLite) CalcularTotal(ctx context.Context, pedidoID string) (float64, error) {
	var total float64
	err := r.db.QueryRowContext(ctx, `
		SELECT COALESCE(SUM(cantidad * precio_unitario), 0)
		FROM items_pedido
		WHERE pedido_id = ?
	`, pedidoID).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("calcular total: %w", err)
	}
	return total, nil
}

func (r *RepositorioSQLite) CambiarEstado(ctx context.Context, pedidoID string, estado string) error {
	resultado, err := r.db.ExecContext(ctx, `
		UPDATE pedidos
		SET estado = ?
		WHERE id = ?
	`, estado, pedidoID)
	if err != nil {
		return fmt.Errorf("cambiar estado: %w", err)
	}

	filas, err := resultado.RowsAffected()
	if err != nil {
		return fmt.Errorf("verificar cambio de estado: %w", err)
	}
	if filas == 0 {
		return fmt.Errorf("pedido no encontrado")
	}

	return nil
}

func (r *RepositorioSQLite) Listar(ctx context.Context, filtroEstado string) ([]modelos.Pedido, error) {
	consulta := `
		SELECT id, cliente_id, total, estado, fecha
		FROM pedidos
	`
	args := []any{}

	if filtroEstado != "" {
		consulta += ` WHERE estado = ? `
		args = append(args, filtroEstado)
	}
	consulta += ` ORDER BY fecha DESC `

	rows, err := r.db.QueryContext(ctx, consulta, args...)
	if err != nil {
		return nil, fmt.Errorf("listar pedidos: %w", err)
	}
	defer rows.Close()

	pedidos := make([]modelos.Pedido, 0)
	for rows.Next() {
		var (
			id        string
			clienteID string
			total     float64
			estado    string
			fecha     string
		)
		if err := rows.Scan(&id, &clienteID, &total, &estado, &fecha); err != nil {
			return nil, fmt.Errorf("leer pedido: %w", err)
		}
		pedido, err := modelos.ReconstruirPedido(id, clienteID, total, estado, fecha)
		if err != nil {
			return nil, fmt.Errorf("reconstruir pedido: %w", err)
		}
		pedidos = append(pedidos, pedido)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("recorrer pedidos: %w", err)
	}

	return pedidos, nil
}
