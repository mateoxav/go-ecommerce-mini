package web

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/mateoxav/go-ecommerce-mini/internal/modelos"
)

func TestCrearProductoHTTP(t *testing.T) {
	servicios := nuevoServiciosMemoria(t)
	manejador := NuevoServidor(servicios, servicios, servicios, servicios).Rutas()

	req := nuevaPeticionJSON(t, http.MethodPost, "/api/productos", `{
		"nombre":"Mouse Gamer",
		"precio":25.50,
		"stock":12,
		"categoria":"Periféricos"
	}`)
	res := httptest.NewRecorder()

	manejador.ServeHTTP(res, req)

	if res.Code != http.StatusCreated {
		t.Fatalf("estado esperado 201, obtenido %d. cuerpo: %s", res.Code, res.Body.String())
	}

	var producto ProductoRespuesta
	decodificarRespuesta(t, res, &producto)
	if producto.ID == "" || producto.Nombre != "Mouse Gamer" || producto.Stock != 12 {
		t.Fatalf("producto inesperado: %+v", producto)
	}
}

func TestBuscarProductoNoEncontradoHTTP(t *testing.T) {
	servicios := nuevoServiciosMemoria(t)
	manejador := NuevoServidor(servicios, servicios, servicios, servicios).Rutas()

	req := httptest.NewRequest(http.MethodGet, "/api/productos/PROD-no-existe", nil)
	res := httptest.NewRecorder()

	manejador.ServeHTTP(res, req)

	if res.Code != http.StatusNotFound {
		t.Fatalf("estado esperado 404, obtenido %d. cuerpo: %s", res.Code, res.Body.String())
	}
}

func TestFlujoAceptacionPedidoHTTP(t *testing.T) {
	servicios := nuevoServiciosMemoria(t)
	manejador := NuevoServidor(servicios, servicios, servicios, servicios).Rutas()

	productoID := crearProductoPorHTTP(t, manejador)
	clienteID := crearClientePorHTTP(t, manejador)
	pedidoID := crearPedidoPorHTTP(t, manejador, clienteID)

	req := nuevaPeticionJSON(t, http.MethodPost, "/api/pedidos/"+pedidoID+"/items", fmt.Sprintf(`{
		"producto_id":%q,
		"cantidad":2
	}`, productoID))
	res := httptest.NewRecorder()
	manejador.ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("estado esperado 200 al agregar item, obtenido %d. cuerpo: %s", res.Code, res.Body.String())
	}

	reqStock := httptest.NewRequest(http.MethodGet, "/api/inventario/stock?id="+productoID+"&cantidad=1", nil)
	resStock := httptest.NewRecorder()
	manejador.ServeHTTP(resStock, reqStock)
	if resStock.Code != http.StatusOK {
		t.Fatalf("estado esperado 200 al verificar stock, obtenido %d. cuerpo: %s", resStock.Code, resStock.Body.String())
	}

	var stock stockRespuesta
	decodificarRespuesta(t, resStock, &stock)
	if stock.StockDisponibleActual != 8 || !stock.Disponible {
		t.Fatalf("stock inesperado: %+v", stock)
	}
}

func TestResumenConcurrenteHTTP(t *testing.T) {
	servicios := nuevoServiciosMemoria(t)
	_, err := servicios.CrearProducto(context.Background(), "Teclado", 30, 2, "Periféricos")
	if err != nil {
		t.Fatal(err)
	}

	manejador := NuevoServidor(servicios, servicios, servicios, servicios).Rutas()
	req := httptest.NewRequest(http.MethodGet, "/api/concurrencia/resumen-inventario?umbral=5&orden=stock", nil)
	res := httptest.NewRecorder()

	manejador.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("estado esperado 200, obtenido %d. cuerpo: %s", res.Code, res.Body.String())
	}

	var resumen resumenConcurrenteRespuesta
	decodificarRespuesta(t, res, &resumen)
	if resumen.Concurrencia.GoroutinesEjecutadas != 3 || resumen.Concurrencia.CanalUsado == "" {
		t.Fatalf("no se evidenció concurrencia en la respuesta: %+v", resumen.Concurrencia)
	}
	if resumen.TotalProductos != 1 || len(resumen.ProductosStockBajo) != 1 {
		t.Fatalf("resumen inesperado: %+v", resumen)
	}
}

func nuevaPeticionJSON(t *testing.T, metodo string, ruta string, cuerpo string) *http.Request {
	t.Helper()
	req := httptest.NewRequest(metodo, ruta, bytes.NewBufferString(cuerpo))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func decodificarRespuesta(t *testing.T, res *httptest.ResponseRecorder, destino any) {
	t.Helper()
	if err := json.NewDecoder(res.Body).Decode(destino); err != nil {
		t.Fatalf("no se pudo decodificar JSON: %v. cuerpo: %s", err, res.Body.String())
	}
}

func crearProductoPorHTTP(t *testing.T, manejador http.Handler) string {
	t.Helper()
	req := nuevaPeticionJSON(t, http.MethodPost, "/api/productos", `{
		"nombre":"Laptop",
		"precio":800,
		"stock":10,
		"categoria":"Tecnología"
	}`)
	res := httptest.NewRecorder()
	manejador.ServeHTTP(res, req)
	if res.Code != http.StatusCreated {
		t.Fatalf("no se pudo crear producto: %d %s", res.Code, res.Body.String())
	}
	var producto ProductoRespuesta
	decodificarRespuesta(t, res, &producto)
	return producto.ID
}

func crearClientePorHTTP(t *testing.T, manejador http.Handler) string {
	t.Helper()
	req := nuevaPeticionJSON(t, http.MethodPost, "/api/clientes", `{
		"nombre":"Ana Pérez",
		"email":"ana@example.com",
		"telefono":"0991234567"
	}`)
	res := httptest.NewRecorder()
	manejador.ServeHTTP(res, req)
	if res.Code != http.StatusCreated {
		t.Fatalf("no se pudo crear cliente: %d %s", res.Code, res.Body.String())
	}
	var cliente ClienteRespuesta
	decodificarRespuesta(t, res, &cliente)
	return cliente.ID
}

func crearPedidoPorHTTP(t *testing.T, manejador http.Handler, clienteID string) string {
	t.Helper()
	req := nuevaPeticionJSON(t, http.MethodPost, "/api/pedidos", fmt.Sprintf(`{"cliente_id":%q}`, clienteID))
	res := httptest.NewRecorder()
	manejador.ServeHTTP(res, req)
	if res.Code != http.StatusCreated {
		t.Fatalf("no se pudo crear pedido: %d %s", res.Code, res.Body.String())
	}
	var pedido PedidoRespuesta
	decodificarRespuesta(t, res, &pedido)
	return pedido.ID
}

type serviciosMemoria struct {
	mu        sync.Mutex
	productos map[string]modelos.Producto
	clientes  map[string]modelos.Cliente
	pedidos   map[string]modelos.Pedido
}

func nuevoServiciosMemoria(t *testing.T) *serviciosMemoria {
	t.Helper()
	return &serviciosMemoria{
		productos: make(map[string]modelos.Producto),
		clientes:  make(map[string]modelos.Cliente),
		pedidos:   make(map[string]modelos.Pedido),
	}
}

func (s *serviciosMemoria) CrearProducto(ctx context.Context, nombre string, precio float64, stock int, categoria string) (modelos.Producto, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	producto, err := modelos.NuevoProducto(nombre, precio, stock, categoria)
	if err != nil {
		return modelos.Producto{}, err
	}
	s.productos[producto.ID()] = producto
	return producto, nil
}

func (s *serviciosMemoria) BuscarProducto(ctx context.Context, id string) (modelos.Producto, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	producto, ok := s.productos[id]
	if !ok {
		return modelos.Producto{}, fmt.Errorf("producto no encontrado")
	}
	return producto, nil
}

func (s *serviciosMemoria) ListarProductos(ctx context.Context) ([]modelos.Producto, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	productos := make([]modelos.Producto, 0, len(s.productos))
	for _, producto := range s.productos {
		productos = append(productos, producto)
	}
	return productos, nil
}

func (s *serviciosMemoria) ActualizarStock(ctx context.Context, id string, cambio int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	producto, ok := s.productos[id]
	if !ok {
		return fmt.Errorf("producto no encontrado")
	}
	productoActualizado, err := producto.ConStock(producto.Stock() + cambio)
	if err != nil {
		return err
	}
	s.productos[id] = productoActualizado
	return nil
}

func (s *serviciosMemoria) EliminarProducto(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.productos[id]; !ok {
		return fmt.Errorf("producto no encontrado")
	}
	delete(s.productos, id)
	return nil
}

func (s *serviciosMemoria) RegistrarCliente(ctx context.Context, nombre string, email string, telefono string) (modelos.Cliente, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	cliente, err := modelos.NuevoCliente(nombre, email, telefono)
	if err != nil {
		return modelos.Cliente{}, err
	}
	s.clientes[cliente.ID()] = cliente
	return cliente, nil
}

func (s *serviciosMemoria) BuscarCliente(ctx context.Context, id string) (modelos.Cliente, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	cliente, ok := s.clientes[id]
	if !ok {
		return modelos.Cliente{}, fmt.Errorf("cliente no encontrado")
	}
	return cliente, nil
}

func (s *serviciosMemoria) ListarClientes(ctx context.Context) ([]modelos.Cliente, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	clientes := make([]modelos.Cliente, 0, len(s.clientes))
	for _, cliente := range s.clientes {
		clientes = append(clientes, cliente)
	}
	return clientes, nil
}

func (s *serviciosMemoria) EliminarCliente(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.clientes[id]; !ok {
		return fmt.Errorf("cliente no encontrado")
	}
	delete(s.clientes, id)
	return nil
}

func (s *serviciosMemoria) CrearPedido(ctx context.Context, clienteID string) (modelos.Pedido, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.clientes[clienteID]; !ok {
		return modelos.Pedido{}, fmt.Errorf("cliente no encontrado")
	}
	pedido, err := modelos.NuevoPedido(clienteID)
	if err != nil {
		return modelos.Pedido{}, err
	}
	s.pedidos[pedido.ID()] = pedido
	return pedido, nil
}

func (s *serviciosMemoria) BuscarPedido(ctx context.Context, id string) (modelos.Pedido, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	pedido, ok := s.pedidos[id]
	if !ok {
		return modelos.Pedido{}, fmt.Errorf("pedido no encontrado")
	}
	return pedido, nil
}

func (s *serviciosMemoria) AgregarItem(ctx context.Context, pedidoID string, productoID string, cantidad int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	pedido, ok := s.pedidos[pedidoID]
	if !ok {
		return fmt.Errorf("pedido no encontrado")
	}
	producto, ok := s.productos[productoID]
	if !ok {
		return fmt.Errorf("producto no encontrado")
	}
	if cantidad <= 0 {
		return fmt.Errorf("la cantidad debe ser mayor que cero")
	}
	if producto.Stock() < cantidad {
		return fmt.Errorf("stock insuficiente")
	}

	productoActualizado, err := producto.ConStock(producto.Stock() - cantidad)
	if err != nil {
		return err
	}
	nuevoTotal := pedido.Total() + float64(cantidad)*producto.Precio()
	pedidoActualizado, err := modelos.ReconstruirPedido(pedido.ID(), pedido.ClienteID(), nuevoTotal, pedido.Estado(), pedido.Fecha())
	if err != nil {
		return err
	}
	s.productos[productoID] = productoActualizado
	s.pedidos[pedidoID] = pedidoActualizado
	return nil
}

func (s *serviciosMemoria) CalcularTotal(ctx context.Context, pedidoID string) (float64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	pedido, ok := s.pedidos[pedidoID]
	if !ok {
		return 0, fmt.Errorf("pedido no encontrado")
	}
	return pedido.Total(), nil
}

func (s *serviciosMemoria) CambiarEstado(ctx context.Context, pedidoID string, estado string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	pedido, ok := s.pedidos[pedidoID]
	if !ok {
		return fmt.Errorf("pedido no encontrado")
	}
	pedidoActualizado, err := modelos.ReconstruirPedido(pedido.ID(), pedido.ClienteID(), pedido.Total(), estado, pedido.Fecha())
	if err != nil {
		return err
	}
	s.pedidos[pedidoID] = pedidoActualizado
	return nil
}

func (s *serviciosMemoria) ListarPedidos(ctx context.Context, filtroEstado string) ([]modelos.Pedido, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	pedidos := make([]modelos.Pedido, 0, len(s.pedidos))
	for _, pedido := range s.pedidos {
		if filtroEstado == "" || strings.EqualFold(pedido.Estado(), filtroEstado) {
			pedidos = append(pedidos, pedido)
		}
	}
	return pedidos, nil
}

func (s *serviciosMemoria) VerificarStock(ctx context.Context, productoID string, cantidad int) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	producto, ok := s.productos[productoID]
	if !ok {
		return false, fmt.Errorf("producto no encontrado")
	}
	return producto.Stock() >= cantidad, nil
}

func (s *serviciosMemoria) AlertasStockBajo(ctx context.Context, umbral int) ([]modelos.Producto, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	productos := make([]modelos.Producto, 0)
	for _, producto := range s.productos {
		if producto.Stock() < umbral {
			productos = append(productos, producto)
		}
	}
	return productos, nil
}

func (s *serviciosMemoria) ReponerStock(ctx context.Context, productoID string, cantidad int) error {
	return s.ActualizarStock(ctx, productoID, cantidad)
}

func (s *serviciosMemoria) GenerarReporte(ctx context.Context, ordenarPor string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var b strings.Builder
	b.WriteString("REPORTE DE INVENTARIO\n")
	for _, producto := range s.productos {
		b.WriteString(fmt.Sprintf("%s | Stock: %d\n", producto.Nombre(), producto.Stock()))
	}
	return b.String(), nil
}
