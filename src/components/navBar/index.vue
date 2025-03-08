<template>
  <nav class="navbar">
    <div class="navbar-container">
      <!-- Logo -->
      <span class="navbar-logo">PDFAnalyst</span>

      <!-- 导航菜单 -->
      <ul class="navbar-menu">
        <li><a href="/">首页</a></li>
        <li><a href="/papers">论文全览</a></li>
        <li><a href="/analysis">数据分析</a></li>
        <li><a href="/wordCloud">云词统计</a></li>
        <li><a href="/rcharts">生成R图</a></li>
      </ul>

      <!-- 移动端导航按钮 -->
      <el-icon class="navbar-toggle" @click="toggleDrawer">
        <Menu style="width: 24px; height: 24px" />
      </el-icon>
    </div>
  </nav>
</template>

<script setup>
import { Menu } from '@element-plus/icons-vue';

// 定义事件
const emit = defineEmits(['toggle-drawer']);

// 点击按钮时触发事件
const toggleDrawer = () => {
  emit('toggle-drawer');
};
</script>

<style scoped>
/* 导航栏样式 */
.navbar {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  background: rgba(16, 24, 48, 0.8);
  backdrop-filter: blur(10px);
  z-index: 2;
}

.navbar-container {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 2rem;
  max-width: 1400px;
  margin: 0 auto;
}

.navbar-logo {
  font-size: 1.5rem;
  font-weight: bold;
  color: #6ab7ff;
}

.navbar-menu {
  display: flex;
  list-style: none;
  gap: 1.5rem;
}

.navbar-menu a {
  color: #fff;
  text-decoration: none;
  font-size: 1rem;
  transition: color 0.3s ease;
}

.navbar-menu a:hover {
  color: #6ab7ff;
}

.navbar-toggle {
  display: none; /* 默认隐藏移动端按钮 */
  font-size: 1.5rem;
  color: #fff;
  cursor: pointer;
}

/* 移动端样式 */
@media (max-width: 768px) {
  .navbar-container {
    padding: 1rem;
  }

  .navbar-menu {
    display: none; /* 移动端隐藏菜单 */
  }

  .navbar-toggle {
    display: block; /* 显示移动端按钮 */
  }
}
</style>