const API_URL = (import.meta.env.VITE_API_URL ?? "http://localhost:8080/api").replace(/\/$/, "")

export class ApiError extends Error {
  estado: number

  constructor(mensaje: string, estado: number) {
    super(mensaje)
    this.name = "ApiError"
    this.estado = estado
  }
}

type Metodo = "GET" | "POST" | "PUT" | "DELETE"

type Opciones = {
  metodo?: Metodo
  cuerpo?: unknown
  signal?: AbortSignal
}

export async function api<T>(ruta: string, opciones: Opciones = {}): Promise<T> {
  const respuesta = await fetch(`${API_URL}${ruta}`, {
    method: opciones.metodo ?? "GET",
    headers: opciones.cuerpo ? { "Content-Type": "application/json" } : undefined,
    body: opciones.cuerpo ? JSON.stringify(opciones.cuerpo) : undefined,
    signal: opciones.signal,
  })

  const texto = await respuesta.text()
  const datos = texto ? JSON.parse(texto) : null

  if (!respuesta.ok) {
    const mensaje = datos?.error ?? `Error HTTP ${respuesta.status}`
    throw new ApiError(mensaje, respuesta.status)
  }

  return datos as T
}

export function construirQuery(parametros: Record<string, string | number | undefined>) {
  const query = new URLSearchParams()
  Object.entries(parametros).forEach(([clave, valor]) => {
    if (valor !== undefined && String(valor).trim() !== "") {
      query.set(clave, String(valor))
    }
  })
  const texto = query.toString()
  return texto ? `?${texto}` : ""
}
