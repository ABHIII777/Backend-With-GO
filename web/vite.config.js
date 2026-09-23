import react from '@vitejs/plugin-react'
import { defineConfig } from 'vite'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  // Dev convenience: `fetch('/users')` on :5173 forwards to the Go server
  // on :8080. Adjust the target if your Go server runs elsewhere.
  // (UI still defaults to API_BASE in src/api.js for production builds.)
  server: {
    proxy: {
      '/users': 'http://localhost:8080',
    },
  },
})
