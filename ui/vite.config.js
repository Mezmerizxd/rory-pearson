import react from '@vitejs/plugin-react';
import viteTsconfigPaths from 'vite-tsconfig-paths';
import { defineConfig } from 'vite';

export default defineConfig({
  base: '/',
  plugins: [react(), viteTsconfigPaths()],
  server: {
    open: true,
    port: 8080,
  },
  build: {
    outDir: './build',
    emptyOutDir: true,
    sourcemap: false,
    minify: true,
    target: 'esnext',
    polyfillDynamicImport: true,
    brotliSize: false,
    chunkSizeWarningLimit: 1500,
  },
});
