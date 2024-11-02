import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  build: {
    outDir : "../sdk/meta/page/",
    commonjsOptions: {
     // transformMixedEsModules: true,
      ignoreTryCatch: false,
    },
  },
})
