import service from '@/utils/request'
import { useUserStore } from '@/pinia/modules/user'

// 普通接口走 axios 实例。
// AI 调用耗时较长，全部关闭全局 loading 遮罩，由面板自身的"正在思考"状态反馈进度。
export function getAiStatus() {
  return service({
    url: '/blog/ai/status',
    method: 'GET',
    donNotShowLoading: true
  })
}

export function generateSummary(payload) {
  return service({
    url: '/blog/ai/summary',
    method: 'POST',
    data: payload,
    donNotShowLoading: true
  })
}

export function suggestTags(payload) {
  return service({
    url: '/blog/ai/suggest-tags',
    method: 'POST',
    data: payload,
    donNotShowLoading: true
  })
}

/**
 * SSE 流式对话（axios 不支持流式，用 fetch + ReadableStream 实现）。
 * @param {Object} payload /blog/ai/chat 请求体
 * @param {Object} handlers { onDelta, onTool, onDone, onError }
 * @returns {{ abort: () => void, promise: Promise<{aborted:boolean, errorMessage:string|null}> }}
 */
export function streamAiChat(payload, handlers = {}) {
  const controller = new AbortController()
  const userStore = useUserStore()

  const run = (async () => {
    let errorMessage = null
    try {
      const baseURL = import.meta.env.VITE_BASE_API || ''
      const response = await fetch(`${baseURL}/blog/ai/chat`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'x-token': userStore.token
        },
        body: JSON.stringify(payload),
        signal: controller.signal
      })

      if (!response.ok) {
        let message = `请求失败(${response.status})`
        try {
          const data = await response.json()
          message = data?.data?.error || data?.msg || message
        } catch (_) {
          /* 非 JSON 响应，保持默认提示 */
        }
        errorMessage = message
        handlers.onError?.(message)
        return { aborted: false, errorMessage }
      }

      const reader = response.body.getReader()
      const decoder = new TextDecoder('utf-8')
      let buffer = ''
      let received = false

      // 解析 SSE 块：event: xxx \n data: {...} \n\n
      const processBlock = (block) => {
        const lines = block.split('\n')
        let event = 'message'
        let dataText = ''
        for (const line of lines) {
          if (line.startsWith('event:')) {
            event = line.slice(6).trim()
          } else if (line.startsWith('data:')) {
            dataText += line.slice(5).trim()
          }
        }
        if (!dataText) return
        let data
        try {
          data = JSON.parse(dataText)
        } catch (_) {
          data = { raw: dataText }
        }
        // 兼容：SSE 通道里混入 gin-vue-admin 统一错误响应（code != 0）
        if (typeof data.code !== 'undefined' && data.code !== 0) {
          errorMessage = data.msg || 'AI 调用失败'
          handlers.onError?.(errorMessage)
          return
        }
        if (event === 'message') {
          if (data.delta) received = true
          handlers.onDelta?.(data.delta || '')
        } else if (event === 'tool') {
          handlers.onTool?.(data)
        } else if (event === 'done') {
          received = true
          handlers.onDone?.(data)
        } else if (event === 'error') {
          errorMessage = data.message || '模型输出中断'
          handlers.onError?.(errorMessage)
        }
      }

      for (;;) {
        const { done, value } = await reader.read()
        if (done) break
        buffer += decoder.decode(value, { stream: true })
        let idx
        while ((idx = buffer.indexOf('\n\n')) >= 0) {
          const block = buffer.slice(0, idx)
          buffer = buffer.slice(idx + 2)
          if (block.trim()) processBlock(block)
        }
      }
      if (buffer.trim()) processBlock(buffer)

      // 流正常结束但一个有效事件都没有（如代理截断、返回了非 SSE 内容）
      if (!received && !errorMessage) {
        errorMessage = 'AI 未返回任何内容，请重试'
        handlers.onError?.(errorMessage)
      }
    } catch (error) {
      if (error?.name === 'AbortError') {
        return { aborted: true, errorMessage: null }
      }
      errorMessage = error?.message || '网络错误'
      handlers.onError?.(errorMessage)
    }
    return { aborted: false, errorMessage }
  })()

  return {
    abort: () => controller.abort(),
    promise: run
  }
}
