"""真实组件浏览器回归。先运行 npm run test:ai:serve；模型接口全部使用本地响应。"""
import argparse
import json
import sys
from playwright.sync_api import sync_playwright, expect

sys.stdout.reconfigure(encoding='utf-8')

parser = argparse.ArgumentParser()
parser.add_argument('--base-url', default='http://127.0.0.1:5191')
parser.add_argument('--channel', default='chrome')
parser.add_argument('--screenshot', help='可选：保存逐段对比的界面截图')
args = parser.parse_args()

with sync_playwright() as playwright:
    browser = playwright.chromium.launch(headless=True, channel=args.channel)
    page = browser.new_page(viewport={'width': 1600, 'height': 1050})
    errors = []
    page.on('pageerror', lambda error: errors.append(str(error)))
    response_text = ['AI 修改']
    defer_response = [False]
    pending = []

    def fulfill_chat(route, done=True):
        body = 'event: message\ndata: ' + json.dumps({'delta': response_text[0]}, ensure_ascii=False) + '\n\n'
        if done:
            body += 'event: done\ndata: {"finishReason":"stop"}\n\n'
        route.fulfill(status=200, content_type='text/event-stream', body=body)

    def chat(route):
        if defer_response[0]:
            pending.append(route)
        else:
            fulfill_chat(route)

    page.route('**/test-api/blog/ai/status', lambda route: route.fulfill(json={'code': 0, 'data': {'enabled': True}}))
    page.route('**/test-api/blog/ai/chat', chat)
    # 无效图片返回 404，避免 Vite 将图片路径回退到正式应用的 index.html。
    page.route('**/missing', lambda route: route.fulfill(status=404, body=''))

    def reset(content):
        page.goto(args.base_url + '/tests/ai-editor/index.html')
        page.wait_for_load_state('networkidle')
        page.wait_for_function('Boolean(window.aiEditorTest)')
        page.evaluate('text => window.aiEditorTest.setContent(text)', content)

    def select(text):
        page.locator('.markdown-textarea').evaluate('''(el, text) => {
          const start = el.value.indexOf(text);
          if (start < 0) throw new Error('selection not found');
          el.focus(); el.setSelectionRange(start, start + text.length);
          document.dispatchEvent(new Event('selectionchange'));
        }''', text.replace('\r\n', '\n').replace('\r', '\n'))

    def generate(text):
        select(text)
        page.get_by_role('button', name='润色', exact=True).click()
        expect(page.locator('.ai-diff-banner')).to_be_visible()

    def content():
        return page.evaluate('window.aiEditorTest.getContent()')

    original = '前文\r\n\r\n待修改😀\r\n\r\n后文\r\n'
    for selected in ['前文', '待修改😀', '后文', original]:
        reset(original)
        response_text[0] = 'AI 修改'
        generate(selected)
        assert content() == original, '生成不能直接改正文'
        if selected != original:
            expect(page.locator('.markdown-preview')).to_contain_text('AI 修改')
            outside = '后文' if selected != '后文' else '前文'
            expect(page.locator('.markdown-preview')).to_contain_text(outside)
        page.get_by_role('button', name='应用修改', exact=True).click()
        assert content() == original.replace(selected, 'AI 修改', 1)
        expect(page.get_by_text('已在正文编辑器中打开 diff 对比', exact=False)).to_have_count(0)
        page.get_by_role('button', name='撤销 AI 修改', exact=True).click()
        assert content() == original, '撤销必须精确恢复 CRLF 和选区外正文'
    print('PASS: 首部/中部/尾部/全文生成、全文预览、应用和精确撤销')

    reset('前文\n\n旧 A\n\n旧 B\n\n后文')
    response_text[0] = '新 A\n\n新 B'
    generate('旧 A\n\n旧 B')
    blocks = page.locator('.diff-block.is-modified')
    expect(blocks).to_have_count(2)
    blocks.nth(0).get_by_role('button', name='保留原文', exact=True).click()
    if args.screenshot:
        page.screenshot(path=args.screenshot, full_page=True)
    page.get_by_role('button', name='应用修改', exact=True).click()
    assert content() == '前文\n\n旧 A\n\n新 B\n\n后文'
    page.locator('.markdown-textarea').fill(content() + '人工新增')
    expect(page.get_by_role('button', name='撤销 AI 修改', exact=True)).to_be_disabled()
    assert content().endswith('人工新增')
    print('PASS: 连续多段部分采纳；后续人工编辑受到撤销冲突保护')

    reset(original)
    generate('待修改😀')
    page.get_by_role('button', name='全部保留原文', exact=True).click()
    page.get_by_role('button', name='应用修改', exact=True).click()
    assert content() == original
    expect(page.get_by_role('button', name='撤销 AI 修改', exact=True)).to_have_count(0)
    generate('待修改😀')
    page.get_by_role('button', name='取消', exact=True).click()
    assert content() == original
    print('PASS: 全部保留、取消对比均不改变原文')

    reset(original)
    generate('待修改😀')
    page.evaluate('text => window.aiEditorTest.setContent(text)', original + '更新')
    expect(page.get_by_role('button', name='应用修改', exact=True)).to_be_disabled()
    assert content() == original + '更新'
    print('PASS: 对比期间正文变化后禁止过期替换')

    for switch_document in [False, True]:
        reset(original)
        defer_response[0] = True
        select('待修改😀')
        page.get_by_role('button', name='润色', exact=True).click()
        expect(page.get_by_role('button', name='停止生成', exact=True)).to_be_visible()
        if switch_document:
            page.evaluate('window.aiEditorTest.changeDocument()')
        else:
            page.evaluate('text => window.aiEditorTest.setContent(text)', original + '更新')
        assert pending
        fulfill_chat(pending.pop())
        expect(page.get_by_role('button', name='停止生成', exact=True)).to_have_count(0)
        expect(page.locator('.ai-diff-banner')).to_have_count(0)
        assert content() == (original if switch_document else original + '更新')
        defer_response[0] = False
    print('PASS: 生成期间改文或切换文章，迟到结果无法进入替换流程')

    reset(original)
    defer_response[0] = True
    select('待修改😀')
    page.get_by_role('button', name='润色', exact=True).click()
    expect(page.get_by_role('button', name='停止生成', exact=True)).to_be_visible()
    fulfill_chat(pending.pop(), done=False)
    expect(page.get_by_role('button', name='停止生成', exact=True)).to_have_count(0)
    expect(page.locator('.ai-diff-banner')).to_have_count(0)
    assert content() == original
    defer_response[0] = False
    print('PASS: 不完整输出不允许进入对比替换流程')

    malicious = '''# 标题
<img src="/missing" onerror="window.__aiXss=1">
<a href="java&#x09;script:window.__aiXss=2">危险链接</a>
<iframe srcdoc="<script>parent.__aiXss=3</script>"></iframe>
<svg onload="window.__aiXss=4"></svg><math><mtext>x</mtext></math>
<style>body{display:none}</style><form id="app"><input name="x" autofocus></form>
<img src="data:image/svg+xml,test"><a href="data:text/html,test">data</a>

| A | B |
| --- | --- |
| 1 | 2 |

- [x] 完成

```js
const text = '<img onerror="bad">'
```
'''
    reset(malicious)

    def assert_safe(selector):
        result = page.locator(selector).evaluate('''el => ({
          forbidden: el.querySelectorAll('script,iframe,svg,math,style,form,object,embed').length,
          badAttributes: [...el.querySelectorAll('*')].flatMap(node => [...node.attributes])
            .filter(a => /^on/i.test(a.name) || ['style','id','name','srcdoc'].includes(a.name)).length,
          badUrls: [...el.querySelectorAll('[href],[src]')].filter(node =>
            /^(javascript|data|vbscript):/i.test(node.getAttribute('href') || node.getAttribute('src') || '')).length,
          table: !!el.querySelector('table'), code: !!el.querySelector('pre code.hljs'),
          checked: !!el.querySelector('input[type="checkbox"][disabled][checked]')
        })''')
        assert result == {'forbidden': 0, 'badAttributes': 0, 'badUrls': 0, 'table': True, 'code': True, 'checked': True}, result
        assert page.evaluate('window.__aiXss === undefined')

    assert_safe('.markdown-preview')
    response_text[0] = malicious
    page.locator('.custom-input textarea').fill('输出测试结果')
    page.get_by_role('button', name='发送', exact=True).click()
    expect(page.locator('.result-body')).to_contain_text('标题')
    assert_safe('.result-body')
    select('标题')
    page.get_by_role('button', name='润色', exact=True).click()
    expect(page.locator('.ai-diff-banner')).to_be_visible()
    assert_safe('.markdown-preview')
    print('PASS: 正文预览、助手结果、diff 预览过滤危险 HTML，保留表格、代码高亮与任务列表')
    assert not errors, errors
    browser.close()
    print('All browser regressions passed.')
