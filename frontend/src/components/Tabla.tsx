import { Box, Text } from "@chakra-ui/react"
import type { ReactNode } from "react"

type Columna<T> = {
  cabecera: string
  celda: (item: T) => ReactNode
  ancho?: string
}

export function Tabla<T>({
  columnas,
  datos,
  obtenerClave,
  vacio = "No hay datos disponibles.",
}: {
  columnas: Columna<T>[]
  datos: T[]
  obtenerClave: (item: T) => string
  vacio?: string
}) {
  if (datos.length === 0) {
    return (
      <Box
        borderWidth="1px"
        borderStyle="dashed"
        borderColor={{ base: "gray.300", _dark: "whiteAlpha.300" }}
        rounded="xl"
        p={5}
        textAlign="center"
      >
        <Text color={{ base: "gray.600", _dark: "gray.400" }}>{vacio}</Text>
      </Box>
    )
  }

  return (
    <Box overflowX="auto" borderWidth="1px" borderColor={{ base: "gray.200", _dark: "whiteAlpha.200" }} rounded="xl">
      <Box as="table" w="full" borderCollapse="collapse" fontSize="sm">
        <Box as="thead" bg={{ base: "gray.50", _dark: "whiteAlpha.100" }}>
          <Box as="tr">
            {columnas.map((columna) => (
              <Box
                as="th"
                key={columna.cabecera}
                w={columna.ancho}
                px={4}
                py={3}
                textAlign="left"
                fontWeight="bold"
                color={{ base: "gray.700", _dark: "gray.200" }}
                borderBottomWidth="1px"
                borderColor={{ base: "gray.200", _dark: "whiteAlpha.200" }}
              >
                {columna.cabecera}
              </Box>
            ))}
          </Box>
        </Box>
        <Box as="tbody">
          {datos.map((item) => (
            <Box as="tr" key={obtenerClave(item)} _hover={{ bg: { base: "gray.50", _dark: "whiteAlpha.50" } }}>
              {columnas.map((columna) => (
                <Box
                  as="td"
                  key={columna.cabecera}
                  px={4}
                  py={3}
                  borderBottomWidth="1px"
                  borderColor={{ base: "gray.100", _dark: "whiteAlpha.100" }}
                  color={{ base: "gray.800", _dark: "gray.100" }}
                  verticalAlign="top"
                >
                  {columna.celda(item)}
                </Box>
              ))}
            </Box>
          ))}
        </Box>
      </Box>
    </Box>
  )
}
