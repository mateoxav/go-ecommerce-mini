import { Button, Grid, HStack, SimpleGrid, Stack, Text } from "@chakra-ui/react"
import { useEffect, useState } from "react"
import { ecommerceApi } from "../api/ecommerce"
import { Acciones } from "../components/Acciones"
import { Campo } from "../components/Campo"
import { Mensaje } from "../components/Mensaje"
import { Panel } from "../components/Panel"
import { Tabla } from "../components/Tabla"
import type { Cliente } from "../types/dominio"
import { fechaLegible, mensajeError } from "../utils"

export function Clientes() {
  const [clientes, setClientes] = useState<Cliente[]>([])
  const [clienteEncontrado, setClienteEncontrado] = useState<Cliente | null>(null)
  const [idBusqueda, setIdBusqueda] = useState("")
  const [mensaje, setMensaje] = useState("")
  const [error, setError] = useState("")
  const [cargando, setCargando] = useState(false)
  const [formulario, setFormulario] = useState({ nombre: "", email: "", telefono: "" })

  async function cargarClientes() {
    setCargando(true)
    setError("")
    try {
      setClientes(await ecommerceApi.listarClientes())
    } catch (err) {
      setError(mensajeError(err))
    } finally {
      setCargando(false)
    }
  }

  useEffect(() => {
    void cargarClientes()
  }, [])

  async function registrarCliente() {
    setMensaje("")
    setError("")
    try {
      const cliente = await ecommerceApi.registrarCliente(formulario)
      setMensaje(`Cliente registrado correctamente: ${cliente.id}`)
      setFormulario({ nombre: "", email: "", telefono: "" })
      await cargarClientes()
    } catch (err) {
      setError(mensajeError(err))
    }
  }

  async function buscarCliente() {
    setMensaje("")
    setError("")
    try {
      const cliente = await ecommerceApi.buscarCliente(idBusqueda.trim())
      setClienteEncontrado(cliente)
      setMensaje("Cliente encontrado correctamente")
    } catch (err) {
      setClienteEncontrado(null)
      setError(mensajeError(err))
    }
  }

  async function eliminarCliente(id: string) {
    setMensaje("")
    setError("")
    try {
      await ecommerceApi.eliminarCliente(id)
      setMensaje("Cliente eliminado de forma lógica")
      await cargarClientes()
    } catch (err) {
      setError(mensajeError(err))
    }
  }

  return (
    <Stack gap={6}>
      <Panel
        titulo="Clientes"
        descripcion="Registro, búsqueda, listado y eliminación lógica de clientes."
        accion={<Button size="sm" onClick={() => void cargarClientes()} disabled={cargando}>Actualizar lista</Button>}
      >
        <Stack gap={4}>
          {mensaje ? <Mensaje tipo="exito" texto={mensaje} /> : null}
          {error ? <Mensaje tipo="error" texto={error} /> : null}
          <Tabla
            datos={clientes}
            obtenerClave={(cliente) => cliente.id}
            columnas={[
              { cabecera: "ID", celda: (cliente) => <Text fontFamily="mono" fontSize="xs">{cliente.id}</Text> },
              { cabecera: "Nombre", celda: (cliente) => cliente.nombre },
              { cabecera: "Email", celda: (cliente) => cliente.email },
              { cabecera: "Teléfono", celda: (cliente) => cliente.telefono || "—" },
              { cabecera: "Registro", celda: (cliente) => fechaLegible(cliente.fecha_registro) },
              {
                cabecera: "Acciones",
                celda: (cliente) => (
                  <Acciones>
                    <Button size="xs" variant="subtle" onClick={() => { setIdBusqueda(cliente.id); setClienteEncontrado(cliente) }}>Ver</Button>
                    <Button size="xs" variant="outline" colorPalette="red" onClick={() => void eliminarCliente(cliente.id)}>Eliminar</Button>
                  </Acciones>
                ),
              },
            ]}
          />
        </Stack>
      </Panel>

      <Grid templateColumns={{ base: "1fr", xl: "1fr 1fr" }} gap={6}>
        <Panel titulo="Registrar cliente" descripcion="El backend valida nombre, email y teléfono.">
          <Stack gap={4}>
            <SimpleGrid columns={{ base: 1, md: 2 }} gap={4}>
              <Campo etiqueta="Nombre" ayuda="Sin números." value={formulario.nombre} onChange={(e) => setFormulario({ ...formulario, nombre: e.target.value })} placeholder="Ana Pérez" />
              <Campo etiqueta="Email" type="email" value={formulario.email} onChange={(e) => setFormulario({ ...formulario, email: e.target.value })} placeholder="ana@example.com" />
              <Campo etiqueta="Teléfono" ayuda="Máximo 10 dígitos." value={formulario.telefono} onChange={(e) => setFormulario({ ...formulario, telefono: e.target.value })} placeholder="0991234567" />
            </SimpleGrid>
            <Button colorPalette="blue" onClick={() => void registrarCliente()}>Registrar cliente</Button>
          </Stack>
        </Panel>

        <Panel titulo="Buscar cliente" descripcion="Consulta directa por ID usando GET /api/clientes/{id}.">
          <Stack gap={4}>
            <HStack align="end" gap={3}>
              <Campo etiqueta="ID del cliente" value={idBusqueda} onChange={(e) => setIdBusqueda(e.target.value)} placeholder="CLI-..." />
              <Button colorPalette="blue" onClick={() => void buscarCliente()}>Buscar</Button>
            </HStack>
            {clienteEncontrado ? (
              <Stack rounded="xl" borderWidth="1px" borderColor={{ base: "gray.200", _dark: "whiteAlpha.200" }} p={4} gap={1}>
                <Text fontWeight="bold">{clienteEncontrado.nombre}</Text>
                <Text fontSize="sm">{clienteEncontrado.email}</Text>
                <Text fontSize="sm">Teléfono: {clienteEncontrado.telefono || "—"}</Text>
                <Text fontFamily="mono" fontSize="xs" color={{ base: "gray.500", _dark: "gray.400" }}>{clienteEncontrado.id}</Text>
              </Stack>
            ) : null}
          </Stack>
        </Panel>
      </Grid>
    </Stack>
  )
}
