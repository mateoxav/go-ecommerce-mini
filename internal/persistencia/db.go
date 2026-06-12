package persistencia

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func AbrirDB(rutaArchivo string) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s?_foreign_keys=on", rutaArchivo)

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("abrir conexión sqlite: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("verificar conexión sqlite: %w", err)
	}

	if err := InicializarTablas(db); err != nil {
		return nil, fmt.Errorf("inicializar tablas: %w", err)
	}

	return db, nil
}

func InicializarTablas(db *sql.DB) error {
	sentencias := []string{
		`CREATE TABLE IF NOT EXISTS productos (
			id TEXT PRIMARY KEY,
			nombre TEXT NOT NULL,
			precio REAL NOT NULL,
			stock INTEGER NOT NULL DEFAULT 0,
			categoria TEXT NOT NULL,
			activo INTEGER NOT NULL DEFAULT 1
		);`,
		`CREATE TABLE IF NOT EXISTS clientes (
			id TEXT PRIMARY KEY,
			nombre TEXT NOT NULL,
			email TEXT NOT NULL,
			telefono TEXT,
			fecha_registro TEXT NOT NULL,
			activo INTEGER NOT NULL DEFAULT 1
		);`,
		`CREATE TABLE IF NOT EXISTS pedidos (
			id TEXT PRIMARY KEY,
			cliente_id TEXT NOT NULL,
			total REAL NOT NULL DEFAULT 0,
			estado TEXT NOT NULL DEFAULT 'pendiente',
			fecha TEXT NOT NULL,
			FOREIGN KEY (cliente_id) REFERENCES clientes(id)
		);`,
		`CREATE TABLE IF NOT EXISTS items_pedido (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			pedido_id TEXT NOT NULL,
			producto_id TEXT NOT NULL,
			cantidad INTEGER NOT NULL,
			precio_unitario REAL NOT NULL,
			FOREIGN KEY (pedido_id) REFERENCES pedidos(id),
			FOREIGN KEY (producto_id) REFERENCES productos(id)
		);`,
	}

	for _, sentencia := range sentencias {
		if _, err := db.Exec(sentencia); err != nil {
			return fmt.Errorf("ejecutar sentencia de creación: %w", err)
		}
	}

	if err := asegurarColumnaClientesActivo(db); err != nil {
		return err
	}

	return nil
}

func asegurarColumnaClientesActivo(db *sql.DB) error {
	columnas, err := columnasDeTabla(db, "clientes")
	if err != nil {
		return fmt.Errorf("verificar columnas de clientes: %w", err)
	}
	if columnas["activo"] {
		return nil
	}

	if _, err := db.Exec(`ALTER TABLE clientes ADD COLUMN activo INTEGER NOT NULL DEFAULT 1`); err != nil {
		return fmt.Errorf("agregar columna activo a clientes: %w", err)
	}
	return nil
}

func columnasDeTabla(db *sql.DB, tabla string) (map[string]bool, error) {
	rows, err := db.Query(fmt.Sprintf("PRAGMA table_info(%s)", tabla))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columnas := make(map[string]bool)
	for rows.Next() {
		var (
			cid          int
			nombre       string
			tipo         string
			noNulo       int
			valorDefault any
			pk           int
		)
		if err := rows.Scan(&cid, &nombre, &tipo, &noNulo, &valorDefault, &pk); err != nil {
			return nil, err
		}
		columnas[nombre] = true
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return columnas, nil
}

func CerrarDB(db *sql.DB) error {
	if db == nil {
		return nil
	}
	return db.Close()
}
