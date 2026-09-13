// 按 SSE 行规则解析，保留跨网络分片的 CRLF、UTF-8 和多行 data。
export async function consumeAiStream(response, handlers = {}) {
  if (!response.ok || !response.headers.get('content-type')?.toLowerCase().startsWith('text/event-stream')) {
    let message = `AI 请求失败（HTTP ${response.status}）`
    try {
      const data = await response.json()
      message = data?.data?.error || data?.msg || message
    } catch (_) { /* 非 JSON 错误页 */ }
    throw new Error(message)
  }
  if (!response.body) throw new Error('AI 响应没有可读取的内容')
  const reader = response.body.getReader()
  const decoder = new TextDecoder()
  let buffer = ''
  let event = 'message'
  let dataLines = []
  let eventSize = 0
  let completed = false
  const line = (value) => {
    if (value === '') {
      if (dataLines.length) {
        let data
        try { data = JSON.parse(dataLines.join('\n')) } catch (_) { throw new Error('AI 返回了无效的事件数据') }
        if (!data || typeof data !== 'object') throw new Error('AI 返回了无效的事件数据')
        if (data.code !== undefined && data.code !== 0) throw new Error(data.msg || 'AI 调用失败')
        if (event === 'error') throw new Error(data.message || '模型输出中断')
        if (event === 'message') {
          if (typeof data.delta !== 'string') throw new Error('AI 返回了无效的文本数据')
          handlers.onDelta?.(data.delta)
        } else if (event === 'context') handlers.onContext?.(data)
        else if (event === 'tool') handlers.onTool?.(data)
        else if (event === 'done') {
          if (data.finishReason !== 'stop') throw new Error('生成未完整结束，请重新生成后再应用')
          completed = true
          handlers.onDone?.(data)
        }
      }
      event = 'message'
      dataLines = []
      eventSize = 0
      return
    }
    if (value.startsWith(':')) return
    const colon = value.indexOf(':')
    const field = colon < 0 ? value : value.slice(0, colon)
    let data = colon < 0 ? '' : value.slice(colon + 1)
    if (data.startsWith(' ')) data = data.slice(1)
    if (field === 'event') event = data
    if (field === 'data') {
      eventSize += data.length
      if (eventSize > 1024 * 1024) throw new Error('AI 单条事件过大')
      dataLines.push(data)
    }
  }
  const drain = (eof = false) => {
    while (!completed) {
      const index = buffer.search(/[\r\n]/)
      if (index < 0) break
      if (buffer[index] === '\r' && index === buffer.length - 1 && !eof) break
      const length = buffer[index] === '\r' && buffer[index + 1] === '\n' ? 2 : 1
      const value = buffer.slice(0, index)
      buffer = buffer.slice(index + length)
      line(value)
    }
    if (buffer.length > 1024 * 1024) throw new Error('AI 单条事件过大')
  }
  try {
    while (!completed) {
      const { done, value } = await reader.read()
      buffer += done ? decoder.decode() : decoder.decode(value, { stream: true })
      drain(done)
      if (done) break
    }
    if (!completed) throw new Error('生成未完整结束，已保留结果供复制，请重新生成后再应用')
    return { completed: true }
  } finally {
    try { await reader.cancel() } catch (_) { /* 连接可能已经取消 */ }
    reader.releaseLock()
  }
}
