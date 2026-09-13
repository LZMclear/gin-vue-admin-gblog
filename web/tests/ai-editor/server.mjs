import { fileURLToPath } from 'node:url'
import { createServer } from 'vite'
import vue from '@vitejs/plugin-vue'

// 独立夹具只加载实际编辑器/助手组件，不连接后台、不改写项目路由映射文件。
const server = await createServer({
  configFile: false,
  root: fileURLToPath(new URL('../../', import.meta.url)),
  cacheDir: 'node_modules/.vite-ai-editor-tests',
  plugins: [vue(), {
    name: 'optional-local-selection-mock',
    configureServer(server) {
      // AI_TEST_MOCK=1 用于手工验证真实抽屉首次挂载和选区交接，不调用模型。
      if (process.env.AI_TEST_MOCK !== '1') return
      server.middlewares.use(async (req, res, next) => {
        if (req.url === '/test-api/blog/ai/status') {
          res.setHeader('Content-Type', 'application/json')
          res.end(JSON.stringify({ code: 0, data: { enabled: true } }))
        } else if (req.url === '/test-api/blog/ai/chat') {
          try {
            let body = ''
            for await (const chunk of req) body += chunk
            const payload = JSON.parse(body)
            res.setHeader('Content-Type', 'text/event-stream')
            res.end(`event: message\ndata: ${JSON.stringify({ delta: `${payload.action}：${payload.selection}` })}\n\nevent: done\ndata: {"finishReason":"stop"}\n\n`)
          } catch {
            res.statusCode = 400
            res.end('Invalid test request')
          }
        } else next()
      })
    }
  }],
  optimizeDeps: { entries: ['tests/ai-editor/index.html'] },
  resolve: { alias: { '@': fileURLToPath(new URL('../../src', import.meta.url)) } },
  define: { 'import.meta.env.VITE_BASE_API': JSON.stringify('/test-api') },
  server: { host: '127.0.0.1', port: Number(process.env.AI_TEST_PORT || 5191), strictPort: true },
  css: { preprocessorOptions: { scss: { api: 'modern-compiler' } } }
})
await server.listen()
server.printUrls()
