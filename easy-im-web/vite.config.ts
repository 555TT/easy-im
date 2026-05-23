import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import UnoCSS from 'unocss/vite'
import path from 'node:path'

const BACKEND_HOST = '192.168.200.129'

export default defineConfig({
  plugins: [vue(), UnoCSS()],
  resolve: {
    alias: { '@': path.resolve(__dirname, 'src') },
  },
  server: {
    port: 3000,
    proxy: {
      '/api/user': {
        target: `http://${BACKEND_HOST}:8880`,
        changeOrigin: true,
        rewrite: (p) => p.replace(/^\/api\/user/, '/v1/user'),
      },
      '/api/social': {
        target: `http://${BACKEND_HOST}:8881`,
        changeOrigin: true,
        rewrite: (p) => p.replace(/^\/api\/social/, '/v1/social'),
      },
      '/api/im': {
        target: `http://${BACKEND_HOST}:8882`,
        changeOrigin: true,
        rewrite: (p) => p.replace(/^\/api\/im/, '/v1/im'),
      },
      '/ws': {
        target: `ws://${BACKEND_HOST}:10090`,
        ws: true,
        changeOrigin: true,
      },
    },
  },
})
