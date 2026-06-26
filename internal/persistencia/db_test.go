package persistencia

import "testing"

func TestAbrirDBInicializaTablasEnMemoria(t *testing.T) {
	db, err := AbrirDB(":memory:")
	if err != nil {
		t.Fatalf("no se pudo abrir la base de datos en memoria: %v", err)
	}
	defer db.Close()

	tablasEsperadas := map[string]bool{
		"productos":    false,
		"clientes":     false,
		"pedidos":      false,
		"items_pedido": false,
	}

	rows, err := db.Query(`SELECT name FROM sqlite_master WHERE type = 'table'`)
	if err != nil {
		t.Fatalf("no se pudo consultar sqlite_master: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var nombre string
		if err := rows.Scan(&nombre); err != nil {
			t.Fatalf("no se pudo leer nombre de tabla: %v", err)
		}
		if _, existe := tablasEsperadas[nombre]; existe {
			tablasEsperadas[nombre] = true
		}
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("error al recorrer tablas: %v", err)
	}

	for tabla, encontrada := range tablasEsperadas {
		if !encontrada {
			t.Fatalf("tabla esperada no encontrada: %s", tabla)
		}
	}
}
