import { Box, Button, Container, Flex, Heading, HStack, Stack, Text } from "@chakra-ui/react"
import { useTheme } from "next-themes"
import type { PropsWithChildren } from "react"
import type { Pagina } from "../App"
import { ColorModeButton } from "./color-mode"

const paginas: Array<{ id: Pagina; etiqueta: string }> = [
  { id: "dashboard", etiqueta: "Dashboard" },
  { id: "productos", etiqueta: "Productos" },
  { id: "clientes", etiqueta: "Clientes" },
  { id: "pedidos", etiqueta: "Pedidos" },
  { id: "inventario", etiqueta: "Inventario" },
  { id: "concurrencia", etiqueta: "Concurrencia" },
]

export function Layout({ paginaActual, cambiarPagina, children }: PropsWithChildren<{ paginaActual: Pagina; cambiarPagina: (pagina: Pagina) => void }>) {
 

  return (
    <Box minH="100vh" bg={{ base: "gray.50", _dark: "#080808" }} color={{ base: "gray.900", _dark: "gray.50" }}>
      <Box
        as="header"
        position="sticky"
        top={0}
        zIndex={10}
        bg={{ base: "whiteAlpha.900", _dark: "blackAlpha.700" }}
        backdropFilter="blur(16px)"
        borderBottomWidth="1px"
        borderColor={{ base: "gray.200", _dark: "whiteAlpha.200" }}
      >
        <Container maxW="7xl" py={4}>
          <Flex direction={{ base: "column", lg: "row" }} align={{ base: "stretch", lg: "center" }} justify="space-between" gap={4}>
            <Box>
              <Heading size="lg">Go E-Commerce Mini</Heading>
              <Text color={{ base: "gray.600", _dark: "gray.400" }} fontSize="sm">
                Frontend en React para probar el backend
              </Text>
            </Box>
            <HStack gap={2} wrap="wrap">
              {paginas.map((pagina) => (
                <Button
                  key={pagina.id}
                  size="sm"
                  variant={paginaActual === pagina.id ? "solid" : "subtle"}
                  colorPalette={paginaActual === pagina.id ? "blue" : "gray"}
                  onClick={() => cambiarPagina(pagina.id)}
                >
                  {pagina.etiqueta}
                </Button>
              ))}
              <ColorModeButton />
            </HStack>
          </Flex>
        </Container>
      </Box>

      <Container maxW="7xl" py={{ base: 6, md: 8 }}>
        <Stack gap={6}>{children}</Stack>
      </Container>
    </Box>
  )
}
