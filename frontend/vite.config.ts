import { fileURLToPath, URL } from 'node:url';

import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
      '@wailsjs': fileURLToPath(new URL('./wailsjs/go', import.meta.url)),
      '@runtime': fileURLToPath(new URL('./wailsjs/runtime', import.meta.url)),
    },
  },
  build: {
    rollupOptions: {
      external: 'chart.js',
      output: {
        entryFileNames: 'assets/[name].js',
        chunkFileNames: 'assets/[name].js',
        assetFileNames: 'assets/[name].[ext]',
      },
    },
  },
});
