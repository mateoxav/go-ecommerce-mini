import { Badge, Box, Button, Grid, HStack, SimpleGrid, Stack, Text } from "@chakra-ui/react"
import { useEffect, useState } from "react"
import { ecommerceApi } from "../api/ecommerce"
import { Acciones } from "../components/Acciones"
import { Campo } from "../components/Campo"
import { Mensaje } from "../components/Mensaje"
import { Panel } from "../components/Panel"
import { Selector } from "../components/Selector"
import { Tabla } from "../components/Tabla"
import type { Pedido, RespuestaAgregarItem, RespuestaTotal } from "../types/dominio"
import { fechaLegible, mensajeError, moneda } from "../utils"

const estados = [
  { valor: "", etiqueta: "Todos" },
  { valor: "pendiente", etiqueta: "Pendiente" },
  { valor: "enviado", etiqueta: "Enviado" },
  { valor: "entregado", etiqueta: "Entregado" },
  { valor: "cancelado", etiqueta: "Cancelado" },
]

const colorEstado: Record<string, string> = {
  pendiente: "yellow",
  enviado: "blue",
  entregado: "green",
  cancelado: "red",
}

export function Pedidos() {
  const [pedidos, setPedidos] = useState<Pedido[]>([])
  const [pedidoEncontrado, setPedidoEncontrado] = useState<Pedido | null>(null)
  const [resultadoItem, setResultadoItem] = useState<RespuestaAgregarItem | null>(null)
  const [resultadoTotal, setResultadoTotal] = useState<RespuestaTotal | null>(null)
  const [estadoFiltro, setEstadoFiltro] = useState("")
  const [clienteID, setClienteID] = useState("")
  const [idBusqueda, setIdBusqueda] = useState("")
  const [pedidoItemID, setPedidoItemID] = useState("")
  const [productoItemID, setProductoItemID] = useState("")
  const [cantidadItem, setCantidadItem] = useState("1")
  const [pedidoEstadoID, setPedidoEstadoID] = useState("")
  const [nuevoEstado, setNuevoEstado] = useState("pendiente")
  const [pedidoTotalID, setPedidoTotalID] = useState("")
  const [mensaje, setMensaje] = useState("")
  const [error, setError] = useState("")
  const [cargando, setCargando] = useState(false)

  async function cargarPedidos() {
    setCargando(true)
    setError("")
    try {
      setPedidos(await ecommerceApi.listarPedidos(estadoFiltro))
    } catch (err) {
      setError(mensajeError(err))
    } finally {
      setCargando(false)
    }
  }

  useEffect(() => {
    void cargarPedidos()
  }, [estadoFiltro])

  async function crearPedido() {
    setMensaje("")
    setError("")
    try {
      const pedido = await ecommerceApi.crearPedido(clienteID.trim())
      setMensaje(`Pedido creado correctamente: ${pedido.id}`)
      setPedidoEncontrado(pedido)
      setPedidoItemID(pedido.id)
      setPedidoTotalID(pedido.id)
      setPedidoEstadoID(pedido.id)
      setClienteID("")
      await cargarPedidos()
    } catch (err) {
      setError(mensajeError(err))
    }
  }

  async function buscarPedido() {
    setMensaje("")
    setError("")
    try {
      const pedido = await ecommerceApi.buscarPedido(idBusqueda.trim())
      setPedidoEncontrado(pedido)
      setMensaje("Pedido encontrado correctamente")
    } catch (err) {
      setPedidoEncontrado(null)
      setError(mensajeError(err))
    }
  }

  async function agregarItem() {
    setMensaje("")
    setError("")
    setResultadoItem(null)
    try {
      const resultado = await ecommerceApi.agregarItemPedido(pedidoItemID.trim(), productoItemID.trim(), Number(cantidadItem))
      setResultadoItem(resultado)
      setMensaje(resultado.mensaje)
      await cargarPedidos()
    } catch (err) {
      setError(mensajeError(err))
    }
  }

  async function calcularTotal() {
    setMensaje("")
    setError("")
    setResultadoTotal(null)
    try {
      const resultado = await ecommerceApi.calcularTotalPedido(pedidoTotalID.trim())
      setResultadoTotal(resultado)
      setMensaje("Total calculado correctamente")
      await cargarPedidos()
    } catch (err) {
      setError(mensajeError(err))
    }
  }

  async function cambiarEstado() {
    setMensaje("")
    setError("")
    try {
      const pedido = await ecommerceApi.cambiarEstadoPedido(pedidoEstadoID.trim(), nuevoEstado)
      setPedidoEncontrado(pedido)
      setMensaje(`Estado actualizado a ${pedido.estado}`)
      await cargarPedidos()
    } catch (err) {
      setError(mensajeError(err))
    }
  }

  return (
    <Stack gap={6}>
      <Panel
        titulo="Pedidos"
        descripcion="Creación, búsqueda, filtros por estado, ítems, total y cambio de estado."
        accion={<Button size="sm" onClick={() => void cargarPedidos()} disabled={cargando}>Actualizar lista</Button>}
      >
        <Stack gap={4}>
          {mensaje ? <Mensaje tipo="exito" texto={mensaje} /> : null}
          {error ? <Mensaje tipo="error" texto={error} /> : null}
          <SimpleGrid columns={{ base: 1, md: 3 }} gap={4} alignItems="end">
            <Selector etiqueta="Filtrar por estado" valor={estadoFiltro} opciones={estados} onChange={setEstadoFiltro} />
          </SimpleGrid>
          <Tabla
            datos={pedidos}
            obtenerClave={(pedido) => pedido.id}
            columnas={[
              { cabecera: "ID", celda: (pedido) => <Text fontFamily="mono" fontSize="xs">{pedido.id}</Text> },
              { cabecera: "Cliente", celda: (pedido) => <Text fontFamily="mono" fontSize="xs">{pedido.cliente_id}</Text> },
              { cabecera: "Total", celda: (pedido) => moneda(pedido.total) },
              { cabecera: "Estado", celda: (pedido) => <Badge colorPalette={colorEstado[pedido.estado] ?? "gray"}>{pedido.estado}</Badge> },
              { cabecera: "Fecha", celda: (pedido) => fechaLegible(pedido.fecha) },
              {
                cabecera: "Acciones",
                celda: (pedido) => (
                  <Acciones>
                    <Button size="xs" variant="subtle" onClick={() => { setIdBusqueda(pedido.id); setPedidoEncontrado(pedido) }}>Ver</Button>
                    <Button size="xs" variant="outline" onClick={() => { setPedidoItemID(pedido.id); setPedidoTotalID(pedido.id); setPedidoEstadoID(pedido.id) }}>Usar ID</Button>
                  </Acciones>
                ),
              },
            ]}
          />
        </Stack>
      </Panel>

      <Grid templateColumns={{ base: "1fr", xl: "1fr 1fr" }} gap={6}>
        <Panel titulo="Crear pedido" descripcion="Crea un pedido vacío asociado a un cliente existente.">
          <Stack gap={4}>
            <Campo etiqueta="ID del cliente" value={clienteID} onChange={(e) => setClienteID(e.target.value)} placeholder="CLI-..." />
            <Button colorPalette="blue" onClick={() => void crearPedido()}>Crear pedido</Button>
          </Stack>
        </Panel>

        <Panel titulo="Buscar pedido" descripcion="Consulta directa por ID usando GET /api/pedidos/{id}.">
          <Stack gap={4}>
            <HStack align="end" gap={3}>
              <Campo etiqueta="ID del pedido" value={idBusqueda} onChange={(e) => setIdBusqueda(e.target.value)} placeholder="PED-..." />
              <Button colorPalette="blue" onClick={() => void buscarPedido()}>Buscar</Button>
            </HStack>
            {pedidoEncontrado ? (
              <Box rounded="xl" borderWidth="1px" borderColor={{ base: "gray.200", _dark: "whiteAlpha.200" }} p={4}>
                <Text fontWeight="bold">Pedido {pedidoEncontrado.estado}</Text>
                <Text fontFamily="mono" fontSize="xs">{pedidoEncontrado.id}</Text>
                <Text mt={2}>Cliente: {pedidoEncontrado.cliente_id}</Text>
                <Text>Total: {moneda(pedidoEncontrado.total)}</Text>
              </Box>
            ) : null}
          </Stack>
        </Panel>
      </Grid>

      <Grid templateColumns={{ base: "1fr", xl: "1fr 1fr" }} gap={6}>
        <Panel titulo="Agregar ítem al pedido" descripcion="Valida pedido, producto, cantidad y stock disponible.">
          <Stack gap={4}>
            <SimpleGrid columns={{ base: 1, md: 3 }} gap={4}>
              <Campo etiqueta="ID del pedido" value={pedidoItemID} onChange={(e) => setPedidoItemID(e.target.value)} placeholder="PED-..." />
              <Campo etiqueta="ID del producto" value={productoItemID} onChange={(e) => setProductoItemID(e.target.value)} placeholder="PROD-..." />
              <Campo etiqueta="Cantidad" type="number" min="1" value={cantidadItem} onChange={(e) => setCantidadItem(e.target.value)} />
            </SimpleGrid>
            <Button colorPalette="blue" onClick={() => void agregarItem()}>Agregar ítem</Button>
            {resultadoItem ? <Mensaje tipo="info" texto={`Pedido ${resultadoItem.pedido_id} · Total actual: ${moneda(resultadoItem.total)}`} /> : null}
          </Stack>
        </Panel>

        <Panel titulo="Total y estado" descripcion="Calcula el total actual y actualiza el estado del pedido.">
          <Stack gap={4}>
            <SimpleGrid columns={{ base: 1, md: 2 }} gap={4}>
              <Campo etiqueta="ID para calcular total" value={pedidoTotalID} onChange={(e) => setPedidoTotalID(e.target.value)} placeholder="PED-..." />
              <Button alignSelf="end" variant="subtle" onClick={() => void calcularTotal()}>Calcular total</Button>
            </SimpleGrid>
            {resultadoTotal ? <Mensaje tipo="info" texto={`Total del pedido: ${moneda(resultadoTotal.total)}`} /> : null}
            <SimpleGrid columns={{ base: 1, md: 3 }} gap={4} alignItems="end">
              <Campo etiqueta="ID del pedido" value={pedidoEstadoID} onChange={(e) => setPedidoEstadoID(e.target.value)} placeholder="PED-..." />
              <Selector etiqueta="Nuevo estado" valor={nuevoEstado} opciones={estados.filter((estado) => estado.valor !== "")} onChange={setNuevoEstado} />
              <Button colorPalette="blue" onClick={() => void cambiarEstado()}>Cambiar estado</Button>
            </SimpleGrid>
          </Stack>
        </Panel>
      </Grid>
    </Stack>
  )
}
