# AI 编辑器第一批回归

在 `web` 目录执行 `npm run test:ai`，验证无损 diff、选区快照、撤销和 Store 接线。

浏览器测试使用 Python Playwright 和本机 Chrome（可通过 `--channel msedge` 使用 Edge）：

```sh
python -m pip install playwright
npm run test:ai:serve
# 在另一个终端执行：
python tests/ai-editor/browser.py
```

独立 Vite 夹具使用真实 MarkdownEditor、WritingAssistantPanel 和 Pinia Store，所有 AI 接口均由测试拦截，不需要登录或模型密钥。默认端口为 5191；可通过 `AI_TEST_PORT` 和测试的 `--base-url` 调整。

覆盖：局部与全文替换、全文预览、精确撤销、CRLF/emoji、逐段采纳、全部保留、取消、编辑冲突、切换文章、残缺流，以及三处 HTML 渲染的清洗和正常 Markdown 保留。
