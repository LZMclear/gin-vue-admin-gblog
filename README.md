# gin-vue-admin-gblog

基于 [gin-vue-admin](https://github.com/flipped-aurora/gin-vue-admin) 改造的个人博客系统。项目将原博客后端能力整合进 GVA 后端，将博客后台管理整合进 GVA 管理端，同时保留独立的博客前台展示端。

## 项目组成

```text
.
├── server      # 后端服务，基于 Gin + GORM + gin-vue-admin
├── web         # 后台管理端，基于 Vue 3 + Vite + Element Plus
├── blog-view   # 博客前台展示端，基于 Vue 2 + Vue CLI + Element UI
├── blog-api    # 原博客后端实现，当前主要作为迁移参考
├── blog-cms    # 原博客后台实现，当前主要作为迁移参考
├── docs        # 项目文档与静态说明资源
└── deploy      # Docker、docker-compose、Kubernetes 部署文件
```

## 功能概览

### 博客前台

- 首页文章列表、置顶、推荐、阅读量展示
- 文章详情、Markdown 内容渲染、密码文章访问校验
- 分类、标签、归档、搜索
- 关于、友链、动态页面
- 评论发布、回复、评论树展示
- QQ 评论信息补全与默认随机头像
- 访问记录、访客记录、文章阅读量统计

### 后台管理

- GVA 原有用户、角色、菜单、API、权限、日志等能力
- 文章管理：新增、编辑、删除、发布、推荐、置顶、密码、评论开关、赞赏开关
- 分类管理、标签管理
- 评论管理：筛选、审核状态、通知状态、编辑、删除
- 友链管理、动态管理、关于页配置、站点配置
- 访问日志、访客记录、异常日志、操作日志
- 博客数据仪表盘

### 后端能力

- 博客业务代码位于 `server/api/v1/blog`、`server/service/blog`、`server/router/blog`
- 公开接口与后台管理接口分离，后台接口复用 GVA 的 JWT 与 Casbin 权限体系
- 博客业务表通过 GORM 自动迁移维护
- Markdown 内容在后端渲染后返回给前台展示
- 访问、操作、异常等日志由博客中间件与业务服务统一记录

## 技术栈

| 模块 | 技术 |
| --- | --- |
| 后端 | Go 1.24、Gin、GORM、JWT、Casbin、Zap、Viper |
| 数据库 | MySQL，默认库名 `gva` |
| 缓存 | Redis，用于多点登录、缓存等 GVA 能力 |
| 后台管理端 | Vue 3、Vite、Element Plus、Pinia、Axios、ECharts |
| 博客前台 | Vue 2、Vue CLI、Vuex、Vue Router、Element UI、Semantic UI |

## 环境要求

- Go 1.24+
- Node.js 18+，推荐 Node.js 20
- MySQL 5.7+ 或 8.x
- Redis

## 后端配置

后端配置文件位于：

```text
server/config.yaml
```

开发前重点检查以下配置：

```yaml
mysql:
  path: 127.0.0.1
  port: "3306"
  db-name: gva
  username: root
  password: "你的数据库密码"
  config: charset=utf8mb4&parseTime=True&loc=Local

redis:
  addr: 127.0.0.1:6379
  password: ""
  db: 0

system:
  db-type: mysql
  addr: 8888
  router-prefix: ""
  use-redis: true
  disable-auto-migrate: false
```

当 `disable-auto-migrate: false` 时，后端启动会自动迁移 GVA 系统表和博客业务表。博客业务模型在 `server/initialize/gorm_biz.go` 中注册。

## 本地启动

### 1. 启动后端

```bash
cd server
go mod tidy
go run .
```

默认服务地址：

```text
http://127.0.0.1:8888
```

### 2. 启动后台管理端

```bash
cd web
npm install
npm run dev
```

`npm run serve` 也会以开发模式启动 Vite。默认访问地址：

```text
http://127.0.0.1:8080
```

后台管理端代理配置位于 `web/.env.development`：

```env
VITE_CLI_PORT = 8080
VITE_SERVER_PORT = 8888
VITE_BASE_API = /api
VITE_BASE_PATH = http://127.0.0.1
```

### 3. 启动博客前台

```bash
cd blog-view
npm install
npm run serve
```

`blog-view/vue.config.js` 已将 `/api` 代理到后端：

```js
proxy: {
  '/api': {
    target: 'http://127.0.0.1:8888/',
    changeOrigin: true,
    pathRewrite: {
      '^/api': ''
    }
  }
}
```

如果使用较新的 Node.js 构建 Vue CLI 旧项目遇到 OpenSSL 相关错误，可临时设置：

```powershell
$env:NODE_OPTIONS="--openssl-legacy-provider"
npm run build
```

## 构建与测试

后台管理端构建：

```bash
cd web
npm run build
```

博客前台构建：

```bash
cd blog-view
npm run build
```

后端测试：

```bash
cd server
go test ./...
```

Makefile 也保留了容器化构建、镜像构建、Swagger 生成、插件打包等命令，例如：

```bash
make build
make image
make doc
make plugin PLUGIN=email
```

## 接口说明

博客公开接口直接挂载在后端根路由下，例如：

```text
GET  /site
GET  /blogs
GET  /blog
GET  /searchBlog
GET  /archives
GET  /category
GET  /category/blogs
GET  /tag
GET  /tag/blogs
GET  /moments
POST /moment/like/:id
GET  /friends
POST /friend
GET  /about
GET  /comments
POST /comment
POST /checkBlogPassword
```

博客后台接口主要挂载在 `/admin` 下，并受 GVA 登录与权限控制，例如：

```text
GET    /admin/dashboard

GET    /admin/blogs
GET    /admin/blog
GET    /admin/categoryAndTag
POST   /admin/blog
PUT    /admin/blog
DELETE /admin/blog
PUT    /admin/blog/top
PUT    /admin/blog/recommend
PUT    /admin/blog/:id/visibility

GET    /admin/categories
POST   /admin/category
PUT    /admin/category
DELETE /admin/category

GET    /admin/tags
POST   /admin/tag
PUT    /admin/tag
DELETE /admin/tag

GET    /admin/comments
PUT    /admin/comment
DELETE /admin/comment

GET    /admin/friends
GET    /admin/moments
GET    /admin/siteSettings
GET    /admin/visitLogs
GET    /admin/visitors
GET    /admin/operationLogs
GET    /admin/exceptionLogs
```

实际接口以 `server/router/blog` 下的路由定义为准。

## 数据表

博客业务表由 GORM 自动迁移，主要模型包括：

- `about`
- `blog`
- `blog_tag`
- `category`
- `city_visitor`
- `comment`
- `exception_log`
- `friend`
- `moment`
- `operation_log`
- `schedule_job_log`
- `site_setting`
- `tag`
- `visit_log`
- `visit_record`
- `visitor`

用户、角色、菜单、权限等能力统一由 GVA 系统模块负责，不再使用旧博客项目的用户表。

## 开发约定

- 后端业务优先写在 `server/service/blog`，路由与 API 层只做参数绑定、调用和响应。
- 后台管理端博客接口位于 `web/src/api/blog`，页面位于 `web/src/view/blog*`、`web/src/view/blogPage`、`web/src/view/blogLog`、`web/src/view/blogStatistics`。
- 博客前台接口位于 `blog-view/src/api`，页面位于 `blog-view/src/views`。
- 前后端字段尽量保持一致，减少前端兼容映射。
- Markdown 内容由后端渲染为 HTML 后返回，前台直接展示渲染结果。
- 分类文章列表和标签文章列表分别通过 `/category/blogs`、`/tag/blogs` 获取。

## 常见问题

### 前端接口访问失败

优先检查：

- 后端是否启动在 `8888` 端口
- `server/config.yaml` 中数据库、Redis 是否可连接
- `web/.env.development` 中 `VITE_BASE_PATH`、`VITE_SERVER_PORT` 是否正确
- `blog-view/vue.config.js` 中 `/api` 代理地址是否正确
- 浏览器 Network 中实际请求是否经过 `/api` 代理
- 后端控制台是否有权限、参数绑定或数据库错误

### 后台出现 `/sysError/createSysError` 请求

这是 GVA 前端捕获运行时错误后上报系统错误日志的请求。需要先查看浏览器控制台中的真实前端错误，再判断是字段为空、组件异常，还是后端返回结构不符合页面预期。

### `web/jsconfig.json` 报错但项目仍可运行

`jsconfig.json` 主要服务于 IDE 路径提示和分析，不直接参与 Vite 编译。只要 Vite alias 配置正确，项目通常仍可正常运行。

## 迁移说明

当前主要运行链路已经迁移到：

- 后端：`server`
- 后台管理端：`web`
- 博客前台：`blog-view`

`blog-api` 和 `blog-cms` 仍保留为原始实现参考。后续开发建议优先修改 `server`、`web`、`blog-view`，避免继续在旧项目中新增业务逻辑。

## License

本项目基于 gin-vue-admin 改造，请遵守原项目 Apache 2.0 协议及相关版权声明。
