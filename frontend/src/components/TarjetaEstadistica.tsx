import { Box, Heading, Text } from "@chakra-ui/react"

export function TarjetaEstadistica({ titulo, valor, detalle }: { titulo: string; valor: string | number; detalle?: string }) {
  return (
    <Box
      borderWidth="1px"
      borderColor={{ base: "gray.200", _dark: "whiteAlpha.200" }}
      rounded="2xl"
      p={5}
      bg={{ base: "white", _dark: "gray.900" }}
    >
      <Text fontSize="sm" color={{ base: "gray.600", _dark: "gray.400" }}>
        {titulo}
      </Text>
      <Heading mt={2} size="xl" color={{ base: "gray.900", _dark: "whiteAlpha.950" }}>
        {valor}
      </Heading>
      {detalle ? (
        <Text mt={2} fontSize="sm" color={{ base: "gray.500", _dark: "gray.500" }}>
          {detalle}
        </Text>
      ) : null}
    </Box>
  )
}
