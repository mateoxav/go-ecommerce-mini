import { Box, Text } from "@chakra-ui/react"

type TipoMensaje = "exito" | "error" | "info"

const estilos: Record<TipoMensaje, { bg: string; border: string; color: string }> = {
  exito: { bg: "green.50", border: "green.300", color: "green.800" },
  error: { bg: "red.50", border: "red.300", color: "red.800" },
  info: { bg: "blue.50", border: "blue.300", color: "blue.800" },
}

export function Mensaje({ tipo = "info", texto }: { tipo?: TipoMensaje; texto: string }) {
  if (!texto) return null
  const estilo = estilos[tipo]
  return (
    <Box
      role="status"
      borderWidth="1px"
      borderColor={{ base: estilo.border, _dark: "whiteAlpha.300" }}
      bg={{ base: estilo.bg, _dark: tipo === "error" ? "red.950" : tipo === "exito" ? "green.950" : "blue.950" }}
      color={{ base: estilo.color, _dark: "whiteAlpha.900" }}
      rounded="xl"
      p={3}
      fontSize="sm"
    >
      <Text>{texto}</Text>
    </Box>
  )
}
