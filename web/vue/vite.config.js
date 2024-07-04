import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

import fs from 'node:fs'
import path from 'node:path'
const secretFile = path.join(fileURLToPath(new URL('.', import.meta.url)), '.secret.config.js')
let c = {}
if (fs.existsSync(secretFile)) {
  c = (await import(secretFile)).default
}
// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    vue(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  server: {
    proxy: {
      '/api': {
        target: c.server_proxy_target || 'http://localhost:5920',
        changeOrigin: true,
        secure: false,
      },
    },
  },
})
