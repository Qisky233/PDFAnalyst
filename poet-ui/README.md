# 品墨千秋 | 数据里的大唐诗卷

一个优雅的唐诗可视化Web应用，展现唐代诗人之间的关系网络和他们的诗作。

## 特性

- 3D星图展示诗人关系网络
- 沉浸式的诗人详情页面
- 响应式设计，支持移动端
- 优雅的动画和交互效果

## 技术栈

- Vue 3
- Vite
- Three.js
- GSAP
- Vue Router
- @vueuse/core

## 开发环境要求

- Node.js >= 18.0.0
- npm >= 9.0.0

## 安装

```bash
# 克隆项目
git clone [项目地址]

# 进入项目目录
cd poet-ui

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

## 构建

```bash
# 构建生产版本
npm run build

# 预览生产构建
npm run preview
```

## 项目结构

```
poet-ui/
├── src/
│   ├── assets/        # 静态资源
│   ├── components/    # 组件
│   ├── views/         # 页面
│   ├── router/        # 路由配置
│   ├── App.vue        # 根组件
│   └── main.js        # 入口文件
├── public/            # 公共资源
├── index.html         # HTML 模板
└── vite.config.js     # Vite 配置
```

## 功能

- 诗人星图：以3D星图的形式展示唐代诗人之间的关系
- 诗人详情：展示诗人的生平、代表作品和相关诗人
- 响应式设计：完美支持移动端浏览
- 沉浸式体验：优雅的动画和交互效果

## 贡献

欢迎提交 Issue 和 Pull Request。

## 许可

MIT License
