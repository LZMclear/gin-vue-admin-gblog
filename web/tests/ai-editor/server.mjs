import { fileURLToPath } from 'node:url'
import { createServer } from 'vite'
import vue from '@vitejs/plugin-vue'

// 独立夹具只加载实际编辑器/助手组件，不连接后台、不改写项目路由映射文件。
const server = await createServer({
  configFile: false,
  root: fileURLToPath(new URL('../../', import.meta.url)),
  cacheDir: 'node_modules/.vite-ai-editor-tests',
  plugins: [vue()],
  optimizeDeps: { entries: ['tests/ai-editor/index.html'] },
  resolve: { alias: { '@': fileURLToPath(new URL('../../src', import.meta.url)) } },
  define: { 'import.meta.env.VITE_BASE_API': JSON.stringify('/test-api') },
  server: { host: '127.0.0.1', port: Number(process.env.AI_TEST_PORT || 5191), strictPort: true },
  css: { preprocessorOptions: { scss: { api: 'modern-compiler' } } }
})
await server.listen()
server.printUrls()
