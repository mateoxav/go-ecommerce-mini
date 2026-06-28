import { Badge, Box, Button, Grid, Heading, HStack, SimpleGrid, Stack, Text } from "@chakra-ui/react"
import { useEffect, useState } from "react"
import { ecommerceApi } from "../api/ecommerce"
import { Mensaje } from "../components/Mensaje"
import { Panel } from "../components/Panel"
import { TarjetaEstadistica } from "../components/TarjetaEstadistica"
import type { Pagina } from "../App"
import type { Cliente, Pedido, Producto, RespuestaSalud } from "../types/dominio"
import { mensajeError } from "../utils"

export function Dashboard({ cambiarPagina }: { cambiarPagina: (pagina: Pagina) => void }) {
  const [productos, setProductos] = useState<Producto[]>([])
  const [clientes, setClientes] = useState<Cliente[]>([])
  const [pedidos, setPedidos] = useState<Pedido[]>([])
  const [salud, setSalud] = useState<RespuestaSalud | null>(null)
  const [error, setError] = useState("")
  const [cargando, setCargando] = useState(false)

  async function cargarDatos() {
    setCargando(true)
    setError("")
    try {
      const [respuestaSalud, listaProductos, listaClientes, listaPedidos] = await Promise.all([
        ecommerceApi.salud(),
        ecommerceApi.listarProductos(),
        ecommerceApi.listarClientes(),
        ecommerceApi.listarPedidos(),
      ])
      setSalud(respuestaSalud)
      setProductos(listaProductos)
      setClientes(listaClientes)
      setPedidos(listaPedidos)
    } catch (err) {
      setError(mensajeError(err))
    } finally {
      setCargando(false)
    }
  }

  useEffect(() => {
    void cargarDatos()
  }, [])

  const stockBajo = productos.filter((producto) => producto.stock < 5).length
  const pedidosPendientes = pedidos.filter((pedido) => pedido.estado === "pendiente").length

  return (
    <Stack gap={6}>
      <Panel
        titulo="Panel principal"
        descripcion="Resumen del backend REST JSON, orientado a probar rápidamente el estado del sistema."
        accion={
          <Button size="sm" onClick={() => void cargarDatos()} disabled={cargando} colorPalette="blue">
            {cargando ? "Actualizando..." : "Actualizar"}
          </Button>
        }
      >
        <Stack gap={4}>
          {error ? <Mensaje tipo="error" texto={error} /> : null}
          <HStack gap={3} wrap="wrap">
            <Badge colorPalette={salud?.estado === "ok" ? "green" : "red"} px={3} py={1} rounded="full">
              Backend: {salud?.estado ?? "sin conexión"}
            </Badge>
            <Badge colorPalette="blue" px={3} py={1} rounded="full">
              API: {import.meta.env.VITE_API_URL ?? "http://localhost:8080/api"}
            </Badge>
            <Badge colorPalette="purple" px={3} py={1} rounded="full">
              Serialización JSON
            </Badge>
          </HStack>
        </Stack>
      </Panel>

      <SimpleGrid columns={{ base: 1, md: 2, xl: 4 }} gap={4}>
        <TarjetaEstadistica titulo="Productos activos" valor={productos.length} detalle="Catálogo disponible" />
        <TarjetaEstadistica titulo="Clientes registrados" valor={clientes.length} detalle="Clientes activos" />
        <TarjetaEstadistica titulo="Pedidos" valor={pedidos.length} detalle={`${pedidosPendientes} pendientes`} />
        <TarjetaEstadistica titulo="Alertas" valor={stockBajo} detalle="Productos con stock menor a 5" />
      </SimpleGrid>

     
    </Stack>
  )
}
