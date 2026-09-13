# AI 编辑器回归（第一批、第二批）

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

第二批增加：SSE 跨分片 CRLF/UTF-8、多行 data、业务错误与完成事件校验；摘要/标签取消后立即重试、旧请求隔离、错误结果禁止回填；文章切换清理历史；模型管理 ID 对接。

后端包含本地模拟供应商的 API 流式完成原因测试（正常结束、输出截断、缺少完成原因），不会调用外部模型。

后端在 `server` 目录运行 `go test ./service/blog ./service/ai ./api/v1/blog ./api/v1/ai`。
配额集成用例需要专用临时 Redis：设置 `AI_TEST_REDIS_ADDR=127.0.0.1:16391` 后运行相同命令并加 `-count=1`；未设置时该集成用例跳过。测试使用 DB 15 的专用用户键，不能指向业务 Redis。

配额约定：按服务器本地日期计算，次日零点过期；状态查询不计次，参数校验失败或无可用配置不计次；通过校验且配置可用的生成请求计一次，模型失败或用户取消不退回。启用限额时 Redis 不可用会拒绝生成并提示。

文章隔离约定：切换文档、编辑器或离开编辑页面时，取消请求并清空助手结果、指令和对话历史；再次返回文章开始新对话。关闭 AI 抽屉保留当前文章状态。
