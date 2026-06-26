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
	"github.com/mateoxav/go-ecommerce-mini/internal/web"
)

func main() {
	ctx := context.Background()
	modo := obtenerModoEjecucion()
	rutaDB := obtenerVariableEntorno("DB_RUTA", "ecommerce.db")

	db, err := persistencia.AbrirDB(rutaDB)
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

	switch modo {
	case "cli":
		ejecutarCLI(ctx, servicioProductos, servicioClientes, servicioPedidos, servicioInventario)
	case "web", "api", "":
		ejecutarServidorWeb(servicioProductos, servicioClientes, servicioPedidos, servicioInventario)
	default:
		log.Fatalf("modo de ejecución no válido: %s. Use web o cli", modo)
	}
}

func ejecutarCLI(
	ctx context.Context,
	servicioProductos *productos.Servicio,
	servicioClientes *clientes.Servicio,
	servicioPedidos *pedidos.Servicio,
	servicioInventario *inventario.Servicio,
) {
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

func ejecutarServidorWeb(
	servicioProductos *productos.Servicio,
	servicioClientes *clientes.Servicio,
	servicioPedidos *pedidos.Servicio,
	servicioInventario *inventario.Servicio,
) {
	direccion := ":" + obtenerVariableEntorno("PUERTO", "8080")
	servidor := web.NuevoServidor(
		servicioProductos,
		servicioClientes,
		servicioPedidos,
		servicioInventario,
	)
	httpServer := web.NuevoHTTPServer(direccion, servidor.Rutas())

	log.Printf("Servidor REST JSON iniciado en http://localhost%s", direccion)
	log.Printf("Modo CLI disponible con: go run . cli")

	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatalf("no se pudo iniciar el servidor web: %v", err)
	}
}

func obtenerModoEjecucion() string {
	if len(os.Args) < 2 {
		return "web"
	}
	return os.Args[1]
}

func obtenerVariableEntorno(nombre string, valorDefecto string) string {
	valor := os.Getenv(nombre)
	if valor == "" {
		return valorDefecto
	}
	return valor
}
