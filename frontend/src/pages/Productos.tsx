import { Badge, Box, Button, Grid, Heading, HStack, SimpleGrid, Stack, Text } from "@chakra-ui/react"
import { useEffect, useState } from "react"
import { ecommerceApi } from "../api/ecommerce"
import { Acciones } from "../components/Acciones"
import { Campo } from "../components/Campo"
import { Mensaje } from "../components/Mensaje"
import { Panel } from "../components/Panel"
import { Tabla } from "../components/Tabla"
import type { Producto } from "../types/dominio"
import { mensajeError, moneda } from "../utils"

export function Productos() {
  const [productos, setProductos] = useState<Producto[]>([])
  const [productoEncontrado, setProductoEncontrado] = useState<Producto | null>(null)
  const [idBusqueda, setIdBusqueda] = useState("")
  const [idStock, setIdStock] = useState("")
  const [cambioStock, setCambioStock] = useState("1")
  const [mensaje, setMensaje] = useState("")
  const [error, setError] = useState("")
  const [cargando, setCargando] = useState(false)
  const [formulario, setFormulario] = useState({ nombre: "", precio: "", stock: "", categoria: "" })

  async function cargarProductos() {
    setCargando(true)
    setError("")
    try {
      setProductos(await ecommerceApi.listarProductos())
    } catch (err) {
      setError(mensajeError(err))
    } finally {
      setCargando(false)
    }
  }

  useEffect(() => {
    void cargarProductos()
  }, [])

  async function crearProducto() {
    setMensaje("")
    setError("")
    try {
      const producto = await ecommerceApi.crearProducto({
        nombre: formulario.nombre,
        precio: Number(formulario.precio),
        stock: Number(formulario.stock),
        categoria: formulario.categoria,
      })
      setMensaje(`Producto creado correctamente: ${producto.id}`)
      setFormulario({ nombre: "", precio: "", stock: "", categoria: "" })
      await cargarProductos()
    } catch (err) {
      setError(mensajeError(err))
    }
  }

  async function buscarProducto() {
    setMensaje("")
    setError("")
    try {
      const producto = await ecommerceApi.buscarProducto(idBusqueda.trim())
      setProductoEncontrado(producto)
      setMensaje("Producto encontrado correctamente")
    } catch (err) {
      setProductoEncontrado(null)
      setError(mensajeError(err))
    }
  }

  async function actualizarStock() {
    setMensaje("")
    setError("")
    try {
      const producto = await ecommerceApi.actualizarStockProducto(idStock.trim(), Number(cambioStock))
      setMensaje(`Stock actualizado. Stock actual de ${producto.nombre}: ${producto.stock}`)
      await cargarProductos()
    } catch (err) {
      setError(mensajeError(err))
    }
  }

  async function eliminarProducto(id: string) {
    setMensaje("")
    setError("")
    try {
      await ecommerceApi.eliminarProducto(id)
      setMensaje("Producto eliminado de forma lógica")
      await cargarProductos()
    } catch (err) {
      setError(mensajeError(err))
    }
  }

  return (
    <Stack gap={6}>
      <Panel
        titulo="Productos"
        descripcion="Gestión del catálogo: crear, listar, buscar, actualizar stock y eliminar lógicamente."
        accion={<Button size="sm" onClick={() => void cargarProductos()} disabled={cargando}>Actualizar lista</Button>}
      >
        <Stack gap={4}>
          {mensaje ? <Mensaje tipo="exito" texto={mensaje} /> : null}
          {error ? <Mensaje tipo="error" texto={error} /> : null}
          <Tabla
            datos={productos}
            obtenerClave={(producto) => producto.id}
            columnas={[
              { cabecera: "ID", celda: (producto) => <Text fontFamily="mono" fontSize="xs">{producto.id}</Text> },
              { cabecera: "Nombre", celda: (producto) => producto.nombre },
              { cabecera: "Categoría", celda: (producto) => <Badge colorPalette="purple">{producto.categoria}</Badge> },
              { cabecera: "Precio", celda: (producto) => moneda(producto.precio) },
              { cabecera: "Stock", celda: (producto) => <Badge colorPalette={producto.stock < 5 ? "red" : "green"}>{producto.stock}</Badge> },
              {
                cabecera: "Acciones",
                celda: (producto) => (
                  <Acciones>
                    <Button size="xs" variant="subtle" onClick={() => { setIdBusqueda(producto.id); setProductoEncontrado(producto) }}>Ver</Button>
                    <Button size="xs" variant="outline" colorPalette="red" onClick={() => void eliminarProducto(producto.id)}>Eliminar</Button>
                  </Acciones>
                ),
              },
            ]}
          />
        </Stack>
      </Panel>

      <Grid templateColumns={{ base: "1fr", xl: "1fr 1fr" }} gap={6}>
        <Panel titulo="Crear producto" descripcion="Los datos se envían por JSON a POST /api/productos.">
          <Stack gap={4}>
            <SimpleGrid columns={{ base: 1, md: 2 }} gap={4}>
              <Campo etiqueta="Nombre" value={formulario.nombre} onChange={(e) => setFormulario({ ...formulario, nombre: e.target.value })} placeholder="Laptop" />
              <Campo etiqueta="Categoría" value={formulario.categoria} onChange={(e) => setFormulario({ ...formulario, categoria: e.target.value })} placeholder="Tecnología" />
              <Campo etiqueta="Precio" type="number" min="0" step="0.01" value={formulario.precio} onChange={(e) => setFormulario({ ...formulario, precio: e.target.value })} placeholder="850" />
              <Campo etiqueta="Stock inicial" type="number" min="0" value={formulario.stock} onChange={(e) => setFormulario({ ...formulario, stock: e.target.value })} placeholder="10" />
            </SimpleGrid>
            <Button colorPalette="blue" onClick={() => void crearProducto()}>Crear producto</Button>
          </Stack>
        </Panel>

        <Panel titulo="Buscar producto" descripcion="Consulta directa por ID usando GET /api/productos/{id}.">
          <Stack gap={4}>
            <HStack align="end" gap={3}>
              <Campo etiqueta="ID del producto" value={idBusqueda} onChange={(e) => setIdBusqueda(e.target.value)} placeholder="PROD-..." />
              <Button colorPalette="blue" onClick={() => void buscarProducto()}>Buscar</Button>
            </HStack>
            {productoEncontrado ? (
              <Box rounded="xl" borderWidth="1px" borderColor={{ base: "gray.200", _dark: "whiteAlpha.200" }} p={4}>
                <Heading size="sm">{productoEncontrado.nombre}</Heading>
                <Text fontFamily="mono" fontSize="xs" color={{ base: "gray.500", _dark: "gray.400" }}>{productoEncontrado.id}</Text>
                <Text mt={2}>Precio: {moneda(productoEncontrado.precio)} · Stock: {productoEncontrado.stock}</Text>
              </Box>
            ) : null}
          </Stack>
        </Panel>
      </Grid>

      <Panel titulo="Actualizar stock" descripcion="El backend recibe un cambio positivo o negativo, no el stock final.">
      <SimpleGrid columns={{ base: 1, md: 3 }} gap={4} alignItems="start">
        <Campo
          etiqueta="ID del producto"
          value={idStock}
          onChange={(e) => setIdStock(e.target.value)}
          placeholder="PROD-..."
        />

        <Campo
          etiqueta="Cambio de stock"
          ayuda="Ejemplo: 5 para sumar, -2 para restar."
          type="number"
          value={cambioStock}
          onChange={(e) => setCambioStock(e.target.value)}
        />

        <Box pt={{ base: 0, md: "28px" }}>
          <Button
            w="full"
            h="40px"
            colorPalette="blue"
            onClick={() => void actualizarStock()}
          >
            Actualizar stock
          </Button>
        </Box>
      </SimpleGrid>
    </Panel>
    </Stack>
  )
}
