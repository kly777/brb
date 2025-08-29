import { defineConfig } from 'vite'
import { dirname, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'

const __dirname = dirname(fileURLToPath(import.meta.url))

export default defineConfig({
  plugins: [],
  server: {
    cors: {
      // 通过浏览器访问的源
      origin: 'http://localhost:8080',
      credentials: true,
    },
    // 新增HMR配置
    hmr: {
      host: 'localhost',
      port: 5173,
      protocol: 'ws', // 使用WebSocket协议
      clientPort: 5173 // 告诉Vite客户端使用8080端口
    }
  },
  build: {
    // 在 outDir 中生成 .vite/manifest.json
    manifest: true,
    rollupOptions: {
      input: {
        index: resolve(__dirname, 'index.html'),
        login: resolve(__dirname, 'partials/login.html'),
      }
    }
  }
})