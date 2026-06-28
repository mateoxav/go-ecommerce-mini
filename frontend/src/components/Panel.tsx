import { Box, Heading, HStack, Text } from "@chakra-ui/react"
import type { PropsWithChildren, ReactNode } from "react"

export function Panel({
  titulo,
  descripcion,
  accion,
  children,
}: PropsWithChildren<{ titulo: string; descripcion?: string; accion?: ReactNode }>) {
  return (
    <Box
      borderWidth="1px"
      borderColor={{ base: "gray.200", _dark: "whiteAlpha.200" }}
      bg={{ base: "white", _dark: "gray.900" }}
      rounded="2xl"
      p={{ base: 4, md: 6 }}
      shadow={{ base: "sm", _dark: "none" }}
    >
      <HStack justify="space-between" align="start" gap={4} mb={5}>
        <Box>
          <Heading size="md" color={{ base: "gray.900", _dark: "whiteAlpha.950" }}>
            {titulo}
          </Heading>
          {descripcion ? (
            <Text mt={1} color={{ base: "gray.600", _dark: "gray.400" }} fontSize="sm">
              {descripcion}
            </Text>
          ) : null}
        </Box>
        {accion}
      </HStack>
      {children}
    </Box>
  )
}
