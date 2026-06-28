import { ChakraProvider, defaultSystem } from "@chakra-ui/react"
import { ThemeProvider } from "next-themes"
import type { PropsWithChildren } from "react"

export function Provider({ children }: PropsWithChildren) {
  return (
    <ChakraProvider value={defaultSystem}>
      <ThemeProvider attribute="class" defaultTheme="dark" enableSystem={false} disableTransitionOnChange>
        {children}
      </ThemeProvider>
    </ChakraProvider>
  )
}
