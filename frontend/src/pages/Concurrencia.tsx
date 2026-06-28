import { Badge, Box, Button, HStack, SimpleGrid, Stack, Text } from "@chakra-ui/react"
import { useState } from "react"
import { ecommerceApi } from "../api/ecommerce"
import { Campo } from "../components/Campo"
import { Mensaje } from "../components/Mensaje"
import { Panel } from "../components/Panel"
import { Selector } from "../components/Selector"
import { Tabla } from "../components/Tabla"
import { TarjetaEstadistica } from "../components/TarjetaEstadistica"
import type { RespuestaConcurrencia } from "../types/dominio"
import { mensajeError, moneda } from "../utils"

const opcionesOrden = [
  { valor: "nombre", etiqueta: "Nombre" },
  { valor: "precio", etiqueta: "Precio" },
  { valor: "stock", etiqueta: "Stock" },
]

export function Concurrencia() {
  const [umbral, setUmbral] = useState("5")
  const [orden, setOrden] = useState("stock")
  const [respuesta, setRespuesta] = useState<RespuestaConcurrencia | null>(null)
  const [mensaje, setMensaje] = useState("")
  const [error, setError] = useState("")
  const [cargando, setCargando] = useState(false)

  async function ejecutarResumen() {
    setCargando(true)
    setMensaje("")
    setError("")
    setRespuesta(null)
    try {
      const resultado = await ecommerceApi.resumenConcurrente(Number(umbral), orden)
      setRespuesta(resultado)
      setMensaje(resultado.mensaje)
    } catch (err) {
      setError(mensajeError(err))
    } finally {
      setCargando(false)
    }
  }

  return (
    <Stack gap={6}>
      <Panel
        titulo="Concurrencia"
        descripcion="Demostración práctica de goroutines y canales desde el backend."
      >
        <Stack gap={4}>
          {mensaje ? <Mensaje tipo="exito" texto={mensaje} /> : null}
          {error ? <Mensaje tipo="error" texto={error} /> : null}
          <SimpleGrid columns={{ base: 1, md: 3 }} gap={4} alignItems="end">
            <Campo etiqueta="Umbral de stock bajo" type="number" min="0" value={umbral} onChange={(e) => setUmbral(e.target.value)} />
            <Selector etiqueta="Orden del reporte" valor={orden} opciones={opcionesOrden} onChange={setOrden} />
            <Button colorPalette="blue" onClick={() => void ejecutarResumen()} disabled={cargando}>
              {cargando ? "Ejecutando..." : "Ejecutar resumen concurrente"}
            </Button>
          </SimpleGrid>
        </Stack>
      </Panel>

      {respuesta ? (
        <Stack gap={6}>
          <SimpleGrid columns={{ base: 1, md: 4 }} gap={4}>
            <TarjetaEstadistica titulo="Total productos" valor={respuesta.total_productos} detalle="Consulta concurrente" />
            <TarjetaEstadistica titulo="Stock bajo" valor={respuesta.productos_stock_bajo.length} detalle={`Umbral menor a ${umbral}`} />
            <TarjetaEstadistica titulo="Tiempo" valor={`${respuesta.tiempo_ms} ms`} detalle="Tiempo informado por el backend" />
            <TarjetaEstadistica titulo="Goroutines" valor={respuesta.concurrencia.goroutines_ejecutadas} detalle={respuesta.concurrencia.canal_usado} />
          </SimpleGrid>

          <Panel titulo="Resultado técnico" descripcion="Datos útiles para explicar la implementación de concurrencia en el video o informe.">
            <HStack gap={3} wrap="wrap">
              <Badge colorPalette="blue" px={3} py={1} rounded="full">goroutines: {respuesta.concurrencia.goroutines_ejecutadas}</Badge>
              <Badge colorPalette="purple" px={3} py={1} rounded="full">canal: {respuesta.concurrencia.canal_usado}</Badge>
              <Badge colorPalette={respuesta.errores?.length ? "red" : "green"} px={3} py={1} rounded="full">errores: {respuesta.errores?.length ?? 0}</Badge>
            </HStack>
          </Panel>

          <Panel titulo="Productos con stock bajo" descripcion="Resultado del cálculo concurrente de alertas.">
            <Tabla
              datos={respuesta.productos_stock_bajo}
              obtenerClave={(producto) => producto.id}
              columnas={[
                { cabecera: "Producto", celda: (producto) => producto.nombre },
                { cabecera: "ID", celda: (producto) => <Text fontFamily="mono" fontSize="xs">{producto.id}</Text> },
                { cabecera: "Precio", celda: (producto) => moneda(producto.precio) },
                { cabecera: "Stock", celda: (producto) => <Badge colorPalette="red">{producto.stock}</Badge> },
              ]}
            />
          </Panel>

          <Panel titulo="Reporte generado" descripcion="Texto producido por la goroutine de reporte de inventario.">
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
              {respuesta.reporte}
            </Box>
          </Panel>
        </Stack>
      ) : null}
    </Stack>
  )
}
