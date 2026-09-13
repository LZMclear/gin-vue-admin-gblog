import { test } from 'node:test'
import assert from 'node:assert/strict'
import { consumeAiStream } from '../../src/api/blog/aiStream.js'

function response(text, step = 1) {
  const bytes = new TextEncoder().encode(text)
  let offset = 0
  return new Response(new ReadableStream({ pull(controller) {
    if (offset >= bytes.length) return controller.close()
    controller.enqueue(bytes.slice(offset, offset += step))
  } }), { headers: { 'content-type': 'text/event-stream; charset=utf-8' } })
}

for (const eol of ['\n', '\r\n', '\r']) {
  test(`SSE 分片、UTF-8、多行 data 和 ${JSON.stringify(eol)}`, async () => {
    let text = ''
    let done = 0
    const body = [': heartbeat', '', 'event: message', 'data: {', 'data: "delta":"中文😀"}', '',
      'event: done', 'data: {"finishReason":"stop"}', '', ''].join(eol)
    const result = await consumeAiStream(response(body), { onDelta: d => { text += d }, onDone: () => { done++ } })
    assert.equal(text, '中文😀')
    assert.equal(done, 1)
    assert.equal(result.completed, true)
  })
}

test('HTTP 200 JSON 业务失败保留服务端提示', async () => {
  await assert.rejects(consumeAiStream(Response.json({ code: 7, msg: '今日额度已用完' })), /今日额度已用完/)
})
test('HTTP 403 和非 SSE 页面不会当作成功', async () => {
  await assert.rejects(consumeAiStream(new Response('Forbidden', { status: 403 })), /403/)
  await assert.rejects(consumeAiStream(new Response('<html/>')), /请求失败/)
})
test('有内容但无 done 或残缺 done 都不成功', async () => {
  for (const suffix of ['', 'event: done\ndata: {"finishReason":"stop"}']) {
    await assert.rejects(consumeAiStream(response('data: {"delta":"partial"}\n\n' + suffix)), /未完整结束/)
  }
})
test('错误是终止事件，不能继续接收后续文本或 done', async () => {
  let done = false
  let text = ''
  await assert.rejects(consumeAiStream(response('event: error\ndata: {"message":"上游失败"}\n\ndata: {"delta":"BAD"}\n\nevent: done\ndata: {"finishReason":"stop"}\n\n', 9999), {
    onDelta: d => { text += d }, onDone: () => { done = true }
  }), /上游失败/)
  assert.equal(text, '')
  assert.equal(done, false)
})
test('无效 JSON、字段类型、截断结束原因拒绝应用', async () => {
  for (const text of ['data: nope\n\n', 'data: {"delta":42}\n\n', 'event: done\ndata: {"finishReason":"length"}\n\n']) {
    await assert.rejects(consumeAiStream(response(text)))
  }
})
test('done 后主动释放仍未关闭的连接，忽略后续事件', async () => {
  let cancelled = false
  const stream = new ReadableStream({ start(c) { c.enqueue(new TextEncoder().encode('event: done\ndata: {"finishReason":"stop"}\n\n')) }, cancel() { cancelled = true } })
  await consumeAiStream(new Response(stream, { headers: { 'content-type': 'text/event-stream' } }))
  assert.equal(cancelled, true)
})
