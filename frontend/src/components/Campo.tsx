import { Input, Stack, Text } from "@chakra-ui/react"
import type { ComponentProps } from "react"

type CampoProps = ComponentProps<typeof Input> & {
  etiqueta: string
  ayuda?: string
}

export function Campo({ etiqueta, ayuda, ...props }: CampoProps) {
  return (
    <Stack gap={1.5} align="stretch">
      <Text as="label" fontSize="sm" fontWeight="semibold" color={{ base: "gray.700", _dark: "gray.200" }}>
        {etiqueta}
      </Text>
      <Input
        {...props}
        bg={{ base: "white", _dark: "gray.950" }}
        borderColor={{ base: "gray.300", _dark: "whiteAlpha.300" }}
        _focusVisible={{ outlineColor: "blue.400", borderColor: "blue.400" }}
      />
      {ayuda ? (
        <Text fontSize="xs" color={{ base: "gray.500", _dark: "gray.500" }}>
          {ayuda}
        </Text>
      ) : null}
    </Stack>
  )
}
