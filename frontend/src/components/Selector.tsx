import { Box, Stack, Text } from "@chakra-ui/react"
import type { ChangeEvent } from "react"

type Opcion = {
  valor: string
  etiqueta: string
}

type SelectorProps = {
  etiqueta: string
  valor: string
  opciones: Opcion[]
  onChange: (valor: string) => void
}

export function Selector({ etiqueta, valor, opciones, onChange }: SelectorProps) {
  return (
    <Stack gap={1.5} align="stretch">
      <Text as="label" fontSize="sm" fontWeight="semibold" color={{ base: "gray.700", _dark: "gray.200" }}>
        {etiqueta}
      </Text>
      <Box
        as="select"
        value={valor}
        onChange={(evento: ChangeEvent<HTMLSelectElement>) => onChange(evento.target.value)}
        w="full"
        px={3}
        py={2}
        rounded="md"
        borderWidth="1px"
        borderColor={{ base: "gray.300", _dark: "whiteAlpha.300" }}
        bg={{ base: "white", _dark: "gray.950" }}
        color={{ base: "gray.900", _dark: "gray.100" }}
      >
        {opciones.map((opcion) => (
          <option key={opcion.valor} value={opcion.valor}>
            {opcion.etiqueta}
          </option>
        ))}
      </Box>
    </Stack>
  )
}
