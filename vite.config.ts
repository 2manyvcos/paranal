import react from '@vitejs/plugin-react';
import { dirname, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';
import { defineConfig } from 'vite';
import tsconfigPaths from 'vite-tsconfig-paths';

const __dirname = dirname(fileURLToPath(import.meta.url));

// https://vite.dev/config/
export default defineConfig({
  server: { host: '0.0.0.0', port: 8081, allowedHosts: true },
  plugins: [react(), tsconfigPaths()],
  build: {
    outDir: 'client/dist',
    rollupOptions: {
      input: {
        main: resolve(__dirname, 'index.html'),
        reference: resolve(__dirname, 'api/reference/index.html'),
      },
    },
  },
  publicDir: 'client/public',
  envPrefix: 'PARANAL_',
});
