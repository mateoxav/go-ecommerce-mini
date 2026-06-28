import { useState } from "react"
import { Layout } from "./components/Layout"
import { Dashboard } from "./pages/Dashboard"
import { Productos } from "./pages/Productos"
import { Clientes } from "./pages/Clientes"
import { Pedidos } from "./pages/Pedidos"
import { Inventario } from "./pages/Inventario"
import { Concurrencia } from "./pages/Concurrencia"

export type Pagina = "dashboard" | "productos" | "clientes" | "pedidos" | "inventario" | "concurrencia"

export function App() {
  const [paginaActual, setPaginaActual] = useState<Pagina>("dashboard")

  return (
    <Layout paginaActual={paginaActual} cambiarPagina={setPaginaActual}>
      {paginaActual === "dashboard" && <Dashboard cambiarPagina={setPaginaActual} />}
      {paginaActual === "productos" && <Productos />}
      {paginaActual === "clientes" && <Clientes />}
      {paginaActual === "pedidos" && <Pedidos />}
      {paginaActual === "inventario" && <Inventario />}
      {paginaActual === "concurrencia" && <Concurrencia />}
    </Layout>
  )
}
