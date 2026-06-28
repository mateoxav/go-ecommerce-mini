import { HStack } from "@chakra-ui/react"
import type { PropsWithChildren } from "react"

export function Acciones({ children }: PropsWithChildren) {
  return (
    <HStack gap={2} wrap="wrap" align="center">
      {children}
    </HStack>
  )
}
