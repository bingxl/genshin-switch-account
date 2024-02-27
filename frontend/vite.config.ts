
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// @ts-ignore
import path, { resolve } from 'path';
// @ts-ignore
import { fileURLToPath } from 'url';

const __dirname = path.dirname(fileURLToPath(import.meta.url))

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  resolve: {
    // @ts-ignore
    alias: {
      '@': resolve(__dirname, 'src'),
      '@go': resolve(__dirname, 'wailsjs/go'),
      '@runtime': resolve(__dirname, 'wailsjs/runtime'),
    }
  }
})
