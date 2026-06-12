package clientes

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

func (r *RepositorioSQLite) Crear(ctx context.Context, cliente modelos.Cliente) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO clientes (id, nombre, email, telefono, fecha_registro)
		VALUES (?, ?, ?, ?, ?)
	`, cliente.ID(), cliente.Nombre(), cliente.Email(), cliente.Telefono(), cliente.FechaRegistro())
	if err != nil {
		return fmt.Errorf("crear cliente: %w", err)
	}
	return nil
}

func (r *RepositorioSQLite) BuscarPorID(ctx context.Context, id string) (modelos.Cliente, error) {
	var (
		clienteID     string
		nombre        string
		email         string
		telefono      string
		fechaRegistro string
	)

	err := r.db.QueryRowContext(ctx, `
		SELECT id, nombre, email, telefono, fecha_registro
		FROM clientes
		WHERE id = ?
	`, id).Scan(&clienteID, &nombre, &email, &telefono, &fechaRegistro)
	if err != nil {
		if err == sql.ErrNoRows {
			return modelos.Cliente{}, fmt.Errorf("cliente no encontrado: %w", err)
		}
		return modelos.Cliente{}, fmt.Errorf("buscar cliente: %w", err)
	}

	cliente, err := modelos.ReconstruirCliente(clienteID, nombre, email, telefono, fechaRegistro)
	if err != nil {
		return modelos.Cliente{}, fmt.Errorf("reconstruir cliente: %w", err)
	}

	return cliente, nil
}

func (r *RepositorioSQLite) Listar(ctx context.Context) ([]modelos.Cliente, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, nombre, email, telefono, fecha_registro
		FROM clientes
		ORDER BY nombre ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("listar clientes: %w", err)
	}
	defer rows.Close()

	clientes := make([]modelos.Cliente, 0)
	for rows.Next() {
		var (
			id            string
			nombre        string
			email         string
			telefono      string
			fechaRegistro string
		)
		if err := rows.Scan(&id, &nombre, &email, &telefono, &fechaRegistro); err != nil {
			return nil, fmt.Errorf("leer cliente: %w", err)
		}
		cliente, err := modelos.ReconstruirCliente(id, nombre, email, telefono, fechaRegistro)
		if err != nil {
			return nil, fmt.Errorf("reconstruir cliente: %w", err)
		}
		clientes = append(clientes, cliente)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("recorrer clientes: %w", err)
	}

	return clientes, nil
}
