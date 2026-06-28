export type Producto = {
  id: string
  nombre: string
  precio: number
  stock: number
  categoria: string
  activo: boolean
}

export type Cliente = {
  id: string
  nombre: string
  email: string
  telefono: string
  fecha_registro: string
}

export type EstadoPedido = "pendiente" | "enviado" | "entregado" | "cancelado"

export type Pedido = {
  id: string
  cliente_id: string
  total: number
  estado: EstadoPedido | string
  fecha: string
}

export type RespuestaError = {
  error: string
}

export type RespuestaMensaje = {
  mensaje: string
}

export type RespuestaSalud = {
  estado: string
  servicio: string
}

export type RespuestaStock = {
  producto_id: string
  cantidad_solicitada: number
  stock_disponible_actual: number
  disponible: boolean
}

export type RespuestaStockBajo = {
  umbral: number
  productos: Producto[]
}

export type RespuestaReporte = {
  orden: string
  reporte: string
}

export type RespuestaTotal = {
  pedido_id: string
  total: number
}

export type RespuestaAgregarItem = {
  mensaje: string
  pedido_id: string
  total: number
}

export type RespuestaConcurrencia = {
  mensaje: string
  total_productos: number
  productos_stock_bajo: Producto[]
  reporte: string
  errores?: string[]
  tiempo_ms: number
  concurrencia: {
    goroutines_ejecutadas: number
    canal_usado: string
  }
}
