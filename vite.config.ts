import { rmSync } from 'fs';
import path from 'path';
import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';
import electron from 'vite-plugin-electron';
import WindiCss from 'vite-plugin-windicss';
import pkg from './package.json';

rmSync('dist', { recursive: true, force: true }); // v14.14.0

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    WindiCss(),
    electron({
      main: {
        entry: 'src/main/index.ts',
        vite: {
          build: {
            // For Debug
            sourcemap: true,
            outDir: 'dist/electron/main',
          },
          resolve: {
            alias: {
              '@': `${path.resolve(__dirname, './src/')}`,
            },
          },
        },
      },
      preload: {
        input: {
          // You can configure multiple preload here
          index: path.join(__dirname, 'src/preload/index.ts'),
        },

        vite: {
          build: {
            // For Debug
            sourcemap: 'inline',
            outDir: 'dist/electron/preload',
          },
          resolve: {
            alias: {
              '@': `${path.resolve(__dirname, './src/')}`,
            },
          },
        },
      },
      // Enables use of Node.js API in the Renderer-process
      // https://github.com/electron-vite/vite-plugin-electron/tree/main/packages/electron-renderer#electron-renderervite-serve
      renderer: {

      },
    }),
  ],
  resolve: {
    alias: {
      '@': `${path.resolve(__dirname, './src/')}`,
      '@components': `${path.resolve(__dirname, './src/renderer/components/')}`,
    },
  },
  server: process.env.VSCODE_DEBUG ? {
    host: pkg.debug.env.VITE_DEV_SERVER_HOSTNAME,
    port: pkg.debug.env.VITE_DEV_SERVER_PORT,
  } : undefined,
});
