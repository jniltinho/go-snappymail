import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite'

// The admin SPA is served under /admin/ on the ISOLATED admin listener (:7071),
// and its build output lands in web/admin-dist (go:embed'd into the binary).
export default defineConfig({
  base: '/admin/',
  plugins: [vue(), tailwindcss()],
  define: {
    API_BASE: JSON.stringify('/api/v1/admin'),
  },
  server: {
    port: 5273,
    proxy: {
      '/api': {
        target: 'http://localhost:7071',
        changeOrigin: true,
      },
    },
  },
  build: {
    outDir: '../web/admin-dist',
    emptyOutDir: true,
  },
})
