'use client'

import { ChakraProvider } from '@chakra-ui/react'
import { extendTheme, ThemeConfig } from '@chakra-ui/react'

const config: ThemeConfig = {
  initialColorMode: 'dark',
  useSystemColorMode: false,
}

const theme = extendTheme({ 
  config,
  styles: {
    global: {
      'html, body': {
        backgroundColor: 'gray.900',
        color: 'white',
      },
    },
  },
})

export function Providers({ children }: { children: React.ReactNode }) {
  return (
    <ChakraProvider theme={theme} cssVarsRoot="body">
      {children}
    </ChakraProvider>
  );
}