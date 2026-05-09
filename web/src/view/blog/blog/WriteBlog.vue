<template>
  <div class="write-blog-page">
    <el-form :model="form" :rules="formRules" ref="formRef" label-position="top" class="writer-form">
      <div class="writer-shell">
        <main class="writer-main">
          <div class="writer-header">
            <div>
              <div class="writer-title">{{ pageTitle }}</div>
              <div class="writer-subtitle">
                {{ contentStats.words }} 字 · {{ form.readTime || 0 }} 分钟阅读
              </div>
            </div>
            <div class="writer-actions">
              <el-button icon="Back" @click="$router.back()">返回</el-button>
              <el-button type="primary" icon="Check" @click="openPublishDialog">保存</el-button>
            </div>
          </div>

          <section class="section-block">
            <el-row :gutter="16">
              <el-col :xs="24" :sm="24" :md="14">
                <el-form-item label="文章标题" prop="title">
                  <el-input v-model="form.title" placeholder="请输入标题" maxlength="120" show-word-limit />
                </el-form-item>
              </el-col>
              <el-col :xs="24" :sm="24" :md="10">
                <el-form-item label="文章首图 URL" prop="firstPicture">
                  <el-input v-model="form.firstPicture" placeholder="文章首图，用于随机文章展示" clearable />
                </el-form-item>
              </el-col>
            </el-row>

            <el-form-item label="文章摘要" prop="description">
              <mavon-editor
                v-model="form.description"
                height="260px"
                placeholder="请输入文章摘要，支持 Markdown"
              />
            </el-form-item>
          </section>

          <section class="section-block content-section">
            <div class="section-title">
              <span>文章正文</span>
              <el-button size="small" icon="Refresh" @click="recalculateWords">重新统计</el-button>
            </div>
            <el-form-item prop="content">
              <mavon-editor
                ref="contentEditorRef"
                v-model="form.content"
                height="680px"
                placeholder="从这里开始写作，支持标题、引用、列表、代码块、表格、图片和实时预览"
              />
            </el-form-item>
          </section>
        </main>

        <aside class="writer-side">
          <section class="side-panel">
            <div class="side-title">发布设置</div>
            <el-radio-group v-model="radio" class="visibility-group">
              <el-radio :label="1">公开</el-radio>
              <el-radio :label="2">私密</el-radio>
              <el-radio :label="3">密码保护</el-radio>
            </el-radio-group>
            <el-form-item v-if="radio === 3" label="访问密码">
              <el-input v-model="form.password" show-password clearable />
            </el-form-item>
            <div v-if="radio !== 2" class="switch-grid">
              <el-checkbox v-model="form.appreciation">赞赏</el-checkbox>
              <el-checkbox v-model="form.recommend">推荐</el-checkbox>
              <el-checkbox v-model="form.commentEnabled">评论</el-checkbox>
              <el-checkbox v-model="form.top">置顶</el-checkbox>
            </div>
          </section>

          <section class="side-panel">
            <div class="side-title">分类与标签</div>
            <el-form-item label="分类" prop="cate">
              <el-select
                v-model="form.cate"
                placeholder="请选择分类（输入可添加）"
                allow-create
                filterable
                clearable
                class="full-width"
              >
                <el-option
                  v-for="item in categoryList"
                  :key="item.id"
                  :label="item.categoryName"
                  :value="item.id"
                />
              </el-select>
            </el-form-item>
            <el-form-item label="标签" prop="tagList">
              <el-select
                v-model="form.tagList"
                placeholder="请选择标签（输入可添加）"
                allow-create
                filterable
                multiple
                class="full-width"
              >
                <el-option
                  v-for="item in tagList"
                  :key="item.id"
                  :label="item.tagName"
                  :value="item.id"
                />
              </el-select>
            </el-form-item>
          </section>

          <section class="side-panel">
            <div class="side-title">文章数据</div>
            <el-form-item label="字数" prop="words">
              <el-input-number
                v-model="form.words"
                :min="0"
                :controls="false"
                class="full-width"
                @change="autoWords = false"
              />
            </el-form-item>
            <el-form-item label="阅读时长(分钟)" prop="readTime">
              <el-input-number v-model="form.readTime" :min="0" :controls="false" class="full-width" />
            </el-form-item>
            <el-form-item label="浏览次数" prop="views">
              <el-input-number v-model="form.views" :min="0" :controls="false" class="full-width" />
            </el-form-item>
          </section>
        </aside>
      </div>
    </el-form>

    <el-dialog title="确认保存" width="460px" v-model="dialogVisible">
      <div class="publish-summary">
        <div>
          <span>标题</span>
          <strong>{{ form.title || '未填写' }}</strong>
        </div>
        <div>
          <span>可见性</span>
          <strong>{{ visibilityText }}</strong>
        </div>
        <div>
          <span>统计</span>
          <strong>{{ form.words || 0 }} 字 · {{ form.readTime || 0 }} 分钟</strong>
        </div>
      </div>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" icon="Check" @click="submit">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script>
  import {
    getCategoryAndTag,
    saveBlog as createBlog,
    getBlogById,
    updateBlog
  } from '@/api/blog/article'

  const countMarkdownWords = (content = '') => {
    const words = content
      .replace(/```[\s\S]*?```/g, ' ')
      .replace(/[#>*_`~\-[\]()!|]/g, ' ')
      .match(/[\u4e00-\u9fa5]|[a-zA-Z0-9]+/g)
    return words ? words.length : 0
  }

  export default {
    name: 'BlogWriteArticle',
    data() {
      return {
        categoryList: [],
        tagList: [],
        dialogVisible: false,
        radio: 1,
        autoWords: true,
        form: {
          title: '',
          firstPicture: '',
          description: '',
          content: '',
          cate: null,
          tagList: [],
          words: 0,
          readTime: 0,
          views: 0,
          appreciation: false,
          recommend: false,
          commentEnabled: false,
          top: false,
          published: false,
          password: '',
        },
        formRules: {
          title: [{ required: true, message: '请输入标题', trigger: 'change' }],
          firstPicture: [{ required: true, message: '请输入首图链接', trigger: 'change' }],
          description: [{ required: true, message: '请输入文章摘要', trigger: 'change' }],
          content: [{ required: true, message: '请输入文章正文', trigger: 'change' }],
          cate: [{ required: true, message: '请选择分类', trigger: 'change' }],
          tagList: [{ required: true, message: '请选择标签', trigger: 'change' }],
          words: [{ required: true, message: '请输入文章字数', trigger: 'change' }],
        },
      }
    },
    computed: {
      pageTitle() {
        return this.$route.params.id ? '编辑文章' : '写新文章'
      },
      contentStats() {
        const content = this.form.content || ''
        return {
          words: countMarkdownWords(content),
          characters: content.length
        }
      },
      visibilityText() {
        if (this.radio === 1) return '公开'
        if (this.radio === 2) return '私密'
        return '密码保护'
      }
    },
    watch: {
      'form.content'() {
        if (this.autoWords) {
          this.recalculateWords()
        }
      },
      'form.words'(newValue) {
        this.form.readTime = newValue ? Math.max(1, Math.round(newValue / 200)) : 0
      },
    },
    created() {
      this.getData()
      if (this.$route.params.id) {
        this.getBlog(this.$route.params.id)
      }
    },
    methods: {
      getData() {
        getCategoryAndTag().then(res => {
          this.categoryList = res.data.categories
          this.tagList = res.data.tags
        })
      },
      getBlog(id) {
        getBlogById(id).then(res => {
          this.computeCategoryAndTag(res.data)
          this.form = res.data
          this.autoWords = false
          this.radio = this.form.published ? (this.form.password !== '' ? 3 : 1) : 2
        })
      },
      computeCategoryAndTag(blog) {
        blog.cate = blog.category.id
        blog.tagList = []
        blog.tags.forEach(item => {
          blog.tagList.push(item.id)
        })
      },
      recalculateWords() {
        this.autoWords = true
        this.form.words = this.contentStats.words
      },
      openPublishDialog() {
        if (this.autoWords || !this.form.words) {
          this.recalculateWords()
        }
        this.dialogVisible = true
      },
      submit() {
        if (this.radio === 3 && (this.form.password === '' || this.form.password === null)) {
          return this.msgError('密码保护模式必须填写密码！')
        }
        this.$refs.formRef.validate(valid => {
          if (valid) {
            if (this.radio === 2) {
              this.form.appreciation = false
              this.form.recommend = false
              this.form.commentEnabled = false
              this.form.top = false
              this.form.published = false
            } else {
              this.form.published = true
            }
            if (this.radio !== 3) {
              this.form.password = ''
            }
            if (this.$route.params.id) {
              this.form.category = null
              this.form.tags = null
              updateBlog(this.form).then(res => {
                this.msgSuccess(res.msg)
                this.$router.push('/layout/gblog/list')
              })
            } else {
              createBlog(this.form).then(res => {
                this.msgSuccess(res.msg)
                this.$router.push('/layout/gblog/list')
              })
            }
          } else {
            this.dialogVisible = false
            return this.msgError('请填写必要的表单项')
          }
        })
      }
    }
  }
</script>

<style scoped lang="scss">
.write-blog-page {
  min-height: calc(100vh - 100px);
  padding: 16px;
  background: #f5f7fa;
}

.writer-shell {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 320px;
  gap: 16px;
  align-items: start;
}

.writer-main,
.writer-side {
  min-width: 0;
}

.writer-header,
.section-block,
.side-panel {
  border: 1px solid #ebeef5;
  border-radius: 8px;
  background: #fff;
}

.writer-header {
  position: sticky;
  top: 0;
  z-index: 10;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 16px;
  padding: 14px 16px;
}

.writer-title {
  color: #303133;
  font-size: 20px;
  font-weight: 700;
}

.writer-subtitle {
  margin-top: 4px;
  color: #909399;
  font-size: 13px;
}

.writer-actions {
  display: flex;
  flex-shrink: 0;
  gap: 8px;
}

.section-block {
  margin-bottom: 16px;
  padding: 16px;
}

.content-section {
  padding-bottom: 8px;
}

.section-title,
.side-title {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
  color: #303133;
  font-size: 15px;
  font-weight: 700;
}

.writer-side {
  position: sticky;
  top: 16px;
}

.side-panel {
  margin-bottom: 16px;
  padding: 16px;
}

.visibility-group {
  display: grid;
  gap: 8px;
  margin-bottom: 14px;
}

.switch-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8px 12px;
}

.full-width {
  width: 100%;
}

.publish-summary {
  display: grid;
  gap: 12px;

  div {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 16px;
  }

  span {
    color: #909399;
  }

  strong {
    max-width: 280px;
    overflow: hidden;
    color: #303133;
    font-weight: 600;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
}

:deep(.el-form-item:last-child) {
  margin-bottom: 0;
}

:deep(.el-form-item__content > .markdown-editor) {
  width: 100%;
}

@media (max-width: 1200px) {
  .writer-shell {
    grid-template-columns: 1fr;
  }

  .writer-side {
    position: static;
  }
}

@media (max-width: 768px) {
  .write-blog-page {
    padding: 10px;
  }

  .writer-header {
    position: static;
    align-items: flex-start;
    flex-direction: column;
  }

  .writer-actions {
    width: 100%;

    .el-button {
      flex: 1;
    }
  }
}
</style>
