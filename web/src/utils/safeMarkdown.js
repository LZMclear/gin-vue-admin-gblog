import DOMPurify from 'dompurify'
import { Marked } from 'marked'
import { markedHighlight } from 'marked-highlight'
import hljs from 'highlight.js'

const marked = new Marked(
  { gfm: true, breaks: true },
  markedHighlight({
    langPrefix: 'hljs language-',
    highlight(code, lang) {
      return hljs.highlight(code, { language: hljs.getLanguage(lang) ? lang : 'plaintext' }).value
    }
  })
)

const allowedTags = [
  'p', 'br', 'hr', 'h1', 'h2', 'h3', 'h4', 'h5', 'h6', 'blockquote',
  'ul', 'ol', 'li', 'pre', 'code', 'em', 'strong', 'del', 's', 'a', 'img',
  'table', 'thead', 'tbody', 'tfoot', 'tr', 'th', 'td', 'div', 'span',
  'details', 'summary', 'kbd', 'sup', 'sub', 'input'
]

DOMPurify.addHook('uponSanitizeAttribute', (node, data) => {
  if (data.attrName === 'href' || data.attrName === 'src') {
    try {
      const url = new URL(data.attrValue, 'https://markdown.invalid/')
      const protocols = data.attrName === 'href' ? ['http:', 'https:', 'mailto:', 'tel:'] : ['http:', 'https:']
      if (!protocols.includes(url.protocol)) data.keepAttr = false
    } catch (_) {
      data.keepAttr = false
    }
  }
  if (node.nodeName === 'INPUT' && data.attrName === 'type' && data.attrValue !== 'checkbox') {
    data.keepAttr = false
  }
})

DOMPurify.addHook('afterSanitizeAttributes', (node) => {
  if (node.nodeName === 'INPUT') {
    node.setAttribute('type', 'checkbox')
    node.setAttribute('disabled', '')
  }
})

/** 代码高亮完成后再清洗，返回值直接用于 v-html；禁止脚本、嵌入页面及危险 URL。 */
export function renderSafeMarkdown(text) {
  return DOMPurify.sanitize(marked.parse(String(text || '')), {
    ALLOWED_TAGS: allowedTags,
    ALLOWED_ATTR: ['href', 'src', 'alt', 'title', 'class', 'colspan', 'rowspan', 'align', 'start', 'type', 'checked', 'disabled'],
    ALLOW_DATA_ATTR: false,
    ALLOW_ARIA_ATTR: false
  })
}
