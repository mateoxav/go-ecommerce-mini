package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/mateoxav/go-ecommerce-mini/internal/cli"
	"github.com/mateoxav/go-ecommerce-mini/internal/clientes"
	"github.com/mateoxav/go-ecommerce-mini/internal/inventario"
	"github.com/mateoxav/go-ecommerce-mini/internal/pedidos"
	"github.com/mateoxav/go-ecommerce-mini/internal/persistencia"
	"github.com/mateoxav/go-ecommerce-mini/internal/productos"
)

func main() {
	ctx := context.Background()

	db, err := persistencia.AbrirDB("ecommerce.db")
	if err != nil {
		log.Fatalf("no se pudo abrir la base de datos: %v", err)
	}
	defer func() {
		if err := persistencia.CerrarDB(db); err != nil {
			fmt.Fprintf(os.Stderr, "advertencia al cerrar la base de datos: %v\n", err)
		}
	}()

	repoProductos := productos.NuevoRepositorioSQLite(db)
	repoClientes := clientes.NuevoRepositorioSQLite(db)
	repoPedidos := pedidos.NuevoRepositorioSQLite(db)

	servicioProductos := productos.NuevoServicio(repoProductos)
	servicioClientes := clientes.NuevoServicio(repoClientes)
	servicioPedidos := pedidos.NuevoServicio(repoPedidos, repoClientes)
	servicioInventario := inventario.NuevoServicio(repoProductos)

	app := cli.NuevaAplicacion(
		servicioProductos,
		servicioClientes,
		servicioPedidos,
		servicioInventario,
	)

	if err := app.Ejecutar(ctx); err != nil {
		log.Fatalf("error de ejecución: %v", err)
	}
}
