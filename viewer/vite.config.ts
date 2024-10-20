import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  build: {
    outDir : "../ellyn_agent/page/",
    commonjsOptions: {
     // transformMixedEsModules: true,
      ignoreTryCatch: false,
    },
  },
})
