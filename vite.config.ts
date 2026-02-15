import { defineConfig } from 'vite';
import tsconfigPaths from 'vite-tsconfig-paths';
import react from '@vitejs/plugin-react';

// https://vite.dev/config/
export default defineConfig({
  server: { host: '0.0.0.0', port: 8081, allowedHosts: true },
  plugins: [react(), tsconfigPaths()],
  build: { outDir: 'client/dist' },
  publicDir: 'client/public',
  envPrefix: 'PARANAL_',
});
