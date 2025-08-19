// vite.config.ts
import { fileURLToPath, URL } from 'node:url'

import { defineConfig, mergeConfig } from 'vite' // <-- Importa mergeConfig
import vue from '@vitejs/plugin-vue'
import { configDefaults, defineConfig as defineVitestConfig } from 'vitest/config' // <-- Importa de vitest/config

export default mergeConfig(
  defineConfig({ // Configuración principal de Vite
    plugins: [
      vue(),
    ],
    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./src', import.meta.url))
      }
    },
  }),
  defineVitestConfig({ // Configuración específica de Vitest
    test: {
      environment: 'jsdom',
      globals: true,
      include: ['**/*.{test,spec}.{js,mjs,cjs,ts,mts,cts,jsx,tsx}'],
      // Si usas TS, puedes excluir node_modules y otros directorios
      exclude: [...configDefaults.exclude, 'e2e/*'],
      root: fileURLToPath(new URL('./', import.meta.url)) // Asegura que Vitest encuentre los archivos correctamente
    }
  })
);