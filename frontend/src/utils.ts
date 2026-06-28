export function moneda(valor: number) {
  return new Intl.NumberFormat("es-EC", { style: "currency", currency: "USD" }).format(valor)
}

export function fechaLegible(valor: string) {
  if (!valor) return "—"
  const fecha = new Date(valor)
  if (Number.isNaN(fecha.getTime())) return valor
  return new Intl.DateTimeFormat("es-EC", { dateStyle: "medium", timeStyle: "short" }).format(fecha)
}

export function mensajeError(error: unknown) {
  if (error instanceof Error) return error.message
  return "Ocurrió un error inesperado"
}
