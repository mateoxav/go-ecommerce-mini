import { api, construirQuery } from "./clienteApi"
import type {
  Cliente,
  Pedido,
  Producto,
  RespuestaAgregarItem,
  RespuestaConcurrencia,
  RespuestaMensaje,
  RespuestaReporte,
  RespuestaSalud,
  RespuestaStock,
  RespuestaStockBajo,
  RespuestaTotal,
} from "../types/dominio"

export const ecommerceApi = {
  salud: () => api<RespuestaSalud>("/salud"),

  listarProductos: () => api<Producto[]>("/productos"),
  crearProducto: (datos: { nombre: string; precio: number; stock: number; categoria: string }) =>
    api<Producto>("/productos", { metodo: "POST", cuerpo: datos }),
  buscarProducto: (id: string) => api<Producto>(`/productos/${encodeURIComponent(id)}`),
  actualizarStockProducto: (id: string, cambio: number) =>
    api<Producto>(`/productos/${encodeURIComponent(id)}/stock`, { metodo: "PUT", cuerpo: { cambio } }),
  eliminarProducto: (id: string) => api<RespuestaMensaje>(`/productos/${encodeURIComponent(id)}`, { metodo: "DELETE" }),

  listarClientes: () => api<Cliente[]>("/clientes"),
  registrarCliente: (datos: { nombre: string; email: string; telefono: string }) =>
    api<Cliente>("/clientes", { metodo: "POST", cuerpo: datos }),
  buscarCliente: (id: string) => api<Cliente>(`/clientes/${encodeURIComponent(id)}`),
  eliminarCliente: (id: string) => api<RespuestaMensaje>(`/clientes/${encodeURIComponent(id)}`, { metodo: "DELETE" }),

  listarPedidos: (estado?: string) => api<Pedido[]>(`/pedidos${construirQuery({ estado })}`),
  crearPedido: (clienteID: string) => api<Pedido>("/pedidos", { metodo: "POST", cuerpo: { cliente_id: clienteID } }),
  buscarPedido: (id: string) => api<Pedido>(`/pedidos/${encodeURIComponent(id)}`),
  agregarItemPedido: (pedidoID: string, productoID: string, cantidad: number) =>
    api<RespuestaAgregarItem>(`/pedidos/${encodeURIComponent(pedidoID)}/items`, {
      metodo: "POST",
      cuerpo: { producto_id: productoID, cantidad },
    }),
  calcularTotalPedido: (pedidoID: string) => api<RespuestaTotal>(`/pedidos/${encodeURIComponent(pedidoID)}/total`),
  cambiarEstadoPedido: (pedidoID: string, estado: string) =>
    api<Pedido>(`/pedidos/${encodeURIComponent(pedidoID)}/estado`, { metodo: "PUT", cuerpo: { estado } }),

  verificarStock: (productoID: string, cantidad: number) =>
    api<RespuestaStock>(`/inventario/stock${construirQuery({ id: productoID, cantidad })}`),
  alertasStockBajo: (umbral: number) => api<RespuestaStockBajo>(`/inventario/stock-bajo${construirQuery({ umbral })}`),
  reponerStock: (productoID: string, cantidad: number) =>
    api<Producto>("/inventario/reponer", { metodo: "POST", cuerpo: { producto_id: productoID, cantidad } }),
  generarReporteInventario: (orden: string) => api<RespuestaReporte>(`/inventario/reporte${construirQuery({ orden })}`),

  resumenConcurrente: (umbral: number, orden: string) =>
    api<RespuestaConcurrencia>(`/concurrencia/resumen-inventario${construirQuery({ umbral, orden })}`),
}
