import { Badge, Box, Button, Grid, HStack, SimpleGrid, Stack, Text } from "@chakra-ui/react"
import { useState } from "react"
import { ecommerceApi } from "../api/ecommerce"
import { Campo } from "../components/Campo"
import { Mensaje } from "../components/Mensaje"
import { Panel } from "../components/Panel"
import { Selector } from "../components/Selector"
import { Tabla } from "../components/Tabla"
import type { Producto, RespuestaReporte, RespuestaStock } from "../types/dominio"
import { mensajeError, moneda } from "../utils"

const opcionesOrden = [
  { valor: "nombre", etiqueta: "Nombre" },
  { valor: "precio", etiqueta: "Precio" },
  { valor: "stock", etiqueta: "Stock" },
]

export function Inventario() {
  const [productoIDStock, setProductoIDStock] = useState("")
  const [cantidadStock, setCantidadStock] = useState("1")
  const [resultadoStock, setResultadoStock] = useState<RespuestaStock | null>(null)
  const [umbral, setUmbral] = useState("5")
  const [productosStockBajo, setProductosStockBajo] = useState<Producto[]>([])
  const [productoIDReponer, setProductoIDReponer] = useState("")
  const [cantidadReponer, setCantidadReponer] = useState("1")
  const [productoRepuesto, setProductoRepuesto] = useState<Producto | null>(null)
  const [orden, setOrden] = useState("nombre")
  const [reporte, setReporte] = useState<RespuestaReporte | null>(null)
  const [mensaje, setMensaje] = useState("")
  const [error, setError] = useState("")

  async function verificarStock() {
    setMensaje("")
    setError("")
    setResultadoStock(null)
    try {
      const resultado = await ecommerceApi.verificarStock(productoIDStock.trim(), Number(cantidadStock))
      setResultadoStock(resultado)
      setMensaje(resultado.disponible ? "Hay stock suficiente" : "No hay stock suficiente")
    } catch (err) {
      setError(mensajeError(err))
    }
  }

  async function consultarAlertas() {
    setMensaje("")
    setError("")
    try {
      const respuesta = await ecommerceApi.alertasStockBajo(Number(umbral))
      setProductosStockBajo(respuesta.productos)
      setMensaje(`Se encontraron ${respuesta.productos.length} productos con stock menor a ${respuesta.umbral}`)
    } catch (err) {
      setProductosStockBajo([])
      setError(mensajeError(err))
    }
  }

  async function reponerStock() {
    setMensaje("")
    setError("")
    setProductoRepuesto(null)
    try {
      const producto = await ecommerceApi.reponerStock(productoIDReponer.trim(), Number(cantidadReponer))
      setProductoRepuesto(producto)
      setMensaje(`Stock repuesto correctamente. Stock actual: ${producto.stock}`)
    } catch (err) {
      setError(mensajeError(err))
    }
  }

  async function generarReporte() {
    setMensaje("")
    setError("")
    setReporte(null)
    try {
      const resultado = await ecommerceApi.generarReporteInventario(orden)
      setReporte(resultado)
      setMensaje(`Reporte generado por ${resultado.orden}`)
    } catch (err) {
      setError(mensajeError(err))
    }
  }

  return (
    <Stack gap={6}>
      <Panel titulo="Inventario" descripcion="Verificación de stock, alertas, reposición y reporte ordenado.">
        <Stack gap={4}>
          {mensaje ? <Mensaje tipo="exito" texto={mensaje} /> : null}
          {error ? <Mensaje tipo="error" texto={error} /> : null}
        </Stack>
      </Panel>

      <Grid templateColumns={{ base: "1fr", xl: "1fr 1fr" }} gap={6}>
        <Panel titulo="Verificar stock" descripcion="Consulta si existe stock suficiente e informa el stock disponible actual.">
          <Stack gap={4}>
            <SimpleGrid columns={{ base: 1, md: 2 }} gap={4}>
              <Campo etiqueta="ID del producto" value={productoIDStock} onChange={(e) => setProductoIDStock(e.target.value)} placeholder="PROD-..." />
              <Campo etiqueta="Cantidad a verificar" type="number" min="1" value={cantidadStock} onChange={(e) => setCantidadStock(e.target.value)} />
            </SimpleGrid>
            <Button colorPalette="blue" onClick={() => void verificarStock()}>Verificar stock</Button>
            {resultadoStock ? (
              <Box rounded="xl" borderWidth="1px" borderColor={{ base: "gray.200", _dark: "whiteAlpha.200" }} p={4}>
                <HStack justify="space-between" gap={4} wrap="wrap">
                  <Text fontWeight="bold">{resultadoStock.disponible ? "Stock suficiente" : "Stock insuficiente"}</Text>
                  <Badge colorPalette={resultadoStock.disponible ? "green" : "red"}>{resultadoStock.stock_disponible_actual} unidades disponibles</Badge>
                </HStack>
                <Text mt={2} fontSize="sm" color={{ base: "gray.600", _dark: "gray.400" }}>
                  Producto: {resultadoStock.producto_id} · Cantidad solicitada: {resultadoStock.cantidad_solicitada}
                </Text>
              </Box>
            ) : null}
          </Stack>
        </Panel>

        <Panel titulo="Reponer stock" descripcion="Incrementa el stock de un producto específico.">
          <Stack gap={4}>
            <SimpleGrid columns={{ base: 1, md: 2 }} gap={4}>
              <Campo etiqueta="ID del producto" value={productoIDReponer} onChange={(e) => setProductoIDReponer(e.target.value)} placeholder="PROD-..." />
              <Campo etiqueta="Cantidad a reponer" type="number" min="1" value={cantidadReponer} onChange={(e) => setCantidadReponer(e.target.value)} />
            </SimpleGrid>
            <Button colorPalette="blue" onClick={() => void reponerStock()}>Reponer stock</Button>
            {productoRepuesto ? (
              <Mensaje tipo="info" texto={`${productoRepuesto.nombre}: ${productoRepuesto.stock} unidades disponibles`} />
            ) : null}
          </Stack>
        </Panel>
      </Grid>

      <Grid templateColumns={{ base: "1fr", xl: "1fr 1fr" }} gap={6}>
        <Panel titulo="Alertas de stock bajo" descripcion="Lista productos con stock menor al umbral indicado.">
          <Stack gap={4}>
            <HStack align="end" gap={3}>
              <Campo etiqueta="Umbral" type="number" min="0" value={umbral} onChange={(e) => setUmbral(e.target.value)} />
              <Button colorPalette="blue" onClick={() => void consultarAlertas()}>Consultar alertas</Button>
            </HStack>
            <Tabla
              datos={productosStockBajo}
              obtenerClave={(producto) => producto.id}
              vacio="No hay alertas cargadas."
              columnas={[
                { cabecera: "Producto", celda: (producto) => producto.nombre },
                { cabecera: "ID", celda: (producto) => <Text fontFamily="mono" fontSize="xs">{producto.id}</Text> },
                { cabecera: "Precio", celda: (producto) => moneda(producto.precio) },
                { cabecera: "Stock", celda: (producto) => <Badge colorPalette="red">{producto.stock}</Badge> },
              ]}
            />
          </Stack>
        </Panel>

        <Panel titulo="Reporte de inventario" descripcion="Genera un reporte textual ordenado por nombre, precio o stock.">
          <Stack gap={4}>
            <SimpleGrid columns={{ base: 1, md: 2 }} gap={4} alignItems="end">
              <Selector etiqueta="Ordenar por" valor={orden} opciones={opcionesOrden} onChange={setOrden} />
              <Button colorPalette="blue" onClick={() => void generarReporte()}>Generar reporte</Button>
            </SimpleGrid>
            {reporte ? (
              <Box
                as="pre"
                whiteSpace="pre-wrap"
                overflowX="auto"
                bg={{ base: "gray.50", _dark: "gray.950" }}
                borderWidth="1px"
                borderColor={{ base: "gray.200", _dark: "whiteAlpha.200" }}
                rounded="xl"
                p={4}
                fontSize="sm"
              >
                {reporte.reporte}
              </Box>
            ) : null}
          </Stack>
        </Panel>
      </Grid>
    </Stack>
  )
}
