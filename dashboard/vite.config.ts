import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'
import path from 'path'

// https://vite.dev/config/
// API backend port - must match IMS_API_PORT in /.env (default: 8100)
const API_PORT = process.env.VITE_API_PORT || '8100'

export default defineConfig({
  plugins: [react(), tailwindcss()],
  base: '/',
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  server: {
    port: 5173,
    proxy: {
      // Proxy /api requests to IMS API backend
      '/api': {
        target: `http://localhost:${API_PORT}`,
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, '/v1'),
      },
      // Proxy /v1 requests directly to IMS API
      '/v1': {
        target: `http://localhost:${API_PORT}`,
        changeOrigin: true,
      },
      '/api/v1': {
        target: `http://localhost:${API_PORT}`,
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api\/v1/, '/v1'),
      },
      // WebSocket proxy for real-time features
      '/ws': {
        target: `ws://localhost:${API_PORT}`,
        ws: true,
      },
      // Proxy /hr-api requests to HR SSOT service for CRUD operations
      '/hr-api': {
        target: 'http://localhost:8300',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/hr-api/, '/v1'),
      },
    },
  },
  build: {
    outDir: 'dist',
    sourcemap: true,
    rollupOptions: {
      output: {
        manualChunks: {
          vendor: ['react', 'react-dom'],
        },
      },
    },
  },
})
