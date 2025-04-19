<template>
  <nav class="navbar" :class="{ 'mobile-menu-open': isMobileMenuOpen }" @click.stop>
    <!-- 品牌Logo/名称 -->
    <div class="navbar-brand">
      <router-link to="/" class="logo-link">
        <span class="logo-icon">✨</span>
        <span class="logo-text">诗词星图</span>
      </router-link>
      
      <!-- 移动端菜单按钮 -->
      <button 
        class="mobile-menu-button" 
        @click="toggleMobileMenu"
        aria-label="Toggle navigation"
      >
        <span class="bar" :class="{ 'bar-1': isMobileMenuOpen }"></span>
        <span class="bar" :class="{ 'bar-2': isMobileMenuOpen }"></span>
        <span class="bar" :class="{ 'bar-3': isMobileMenuOpen }"></span>
      </button>
    </div>

    <!-- 居中搜索框 (PC端显示) -->
    <div class="search-container">
      <input
        type="text"
        class="search-input"
        placeholder="支持 诗人/诗 模糊搜索 (部分数据为繁体，不支持简搜繁)"
        v-model="searchKeyword" 
        @keyup.enter="handleSearch"
      />
      <button class="search-button" @click="handleSearch">
        <svg viewBox="0 0 24 24" width="18" height="18">
          <path fill="currentColor" d="M15.5 14h-.79l-.28-.27a6.5 6.5 0 0 0 1.48-5.34c-.47-2.78-2.79-5-5.59-5.34a6.505 6.505 0 0 0-7.27 7.27c.34 2.8 2.56 5.12 5.34 5.59a6.5 6.5 0 0 0 5.34-1.48l.27.28v.79l4.25 4.25c.41.41 1.08.41 1.49 0 .41-.41.41-1.08 0-1.49L15.5 14zm-6 0C7.01 14 5 11.99 5 9.5S7.01 5 9.5 5 14 7.01 14 9.5 11.99 14 9.5 14z"/>
        </svg>
      </button>
    </div>
    
    <!-- 导航链接 -->
    <div class="navbar-links">
      <router-link 
        v-for="link in navLinks" 
        :key="link.path" 
        :to="link.path"
        class="nav-link"
        :class="{ 'active': isActive(link.path) }"
        @click="closeMobileMenu"
      >
        {{ link.text }}
        <span class="link-underline"></span>
      </router-link>
    </div>
  </nav>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue';
import { useRoute, useRouter } from 'vue-router';  // 确保两个都导入

const route = useRoute();
const router = useRouter();
const isMobileMenuOpen = ref(false);

const searchKeyword = ref('');

// 导航链接配置
const navLinks = [
  { path: '/', text: '首页' },
  { path: '/poet', text: '诗人星图' },
  { path: '/data', text: '数据一览' },
//   { path: '/reading', text: '诗词朗诵' }
];

// 检查当前激活路由
const isActive = (path) => {
  return route.path === path || 
         (path !== '/' && route.path.startsWith(path));
};

// 切换移动菜单
const toggleMobileMenu = () => {
  isMobileMenuOpen.value = !isMobileMenuOpen.value;
};

// 关闭移动菜单
const closeMobileMenu = () => {
  if (window.innerWidth <= 768) {
    isMobileMenuOpen.value = false;
  }
};

// 响应式处理
const handleResize = () => {
  if (window.innerWidth > 768) {
    isMobileMenuOpen.value = false;
  }
};

const handleSearch = () => {
    console.log('搜索关键词:', searchKeyword.value);
  if (searchKeyword.value.trim()) {
    router.push({
      path: '/search',
      query: { keyword: searchKeyword.value.trim(), page: 1 }
    });
    searchKeyword.value = '';
    closeMobileMenu();
  }
};

// 添加窗口大小监听
onMounted(() => {
  window.addEventListener('resize', handleResize);
});

onBeforeUnmount(() => {
  window.removeEventListener('resize', handleResize);
});
</script>

<style scoped>
/* 基础样式 */
.navbar {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  padding: 1rem 2rem;
  background: rgba(10, 15, 20, 0.9);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid rgba(255, 204, 0, 0.2);
  box-shadow: 0 2px 20px rgba(0, 0, 0, 0.5);
  z-index: 1000;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.navbar-brand {
  display: flex;
  align-items: center;
}

.logo-link {
  display: flex;
  align-items: center;
  text-decoration: none;
  color: #ffcc00;
  font-size: 1.5rem;
  font-weight: bold;
  transition: all 0.3s;
}

.logo-icon {
  margin-right: 0.5rem;
  font-size: 1.8rem;
}

.logo-text {
  position: relative;
}

.logo-link:hover {
  text-shadow: 0 0 10px rgba(255, 204, 0, 0.5);
}

/* 导航链接样式 */
.navbar-links {
  display: flex;
  gap: 2rem;
}

.nav-link {
  position: relative;
  color: #ddd;
  text-decoration: none;
  font-size: 1.1rem;
  padding: 0.5rem 0;
  transition: all 0.3s;
}

.nav-link:hover {
  color: #fff;
}

.link-underline {
  position: absolute;
  bottom: 0;
  left: 0;
  width: 0;
  height: 2px;
  background: #ffcc00;
  transition: width 0.3s ease;
}

.nav-link:hover .link-underline,
.nav-link.active .link-underline {
  width: 100%;
}

.nav-link.active {
  color: #ffcc00;
  font-weight: bold;
}

/* 搜索框样式 */
.search-container {
  display: flex;
  align-items: center;
  flex-grow: 1;
  max-width: 500px;
  margin: 0 2rem;
  position: relative;
}

.search-input {
  width: 100%;
  padding: 0.5rem 1rem;
  padding-right: 2.5rem;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 204, 0, 0.3);
  border-radius: 20px;
  color: white;
  font-size: 1rem;
  transition: all 0.3s;
}

.search-input:focus {
  outline: none;
  border-color: #ffcc00;
  box-shadow: 0 0 0 2px rgba(255, 204, 0, 0.2);
}

.search-input::placeholder {
  color: rgba(255, 255, 255, 0.5);
}

.search-button {
  position: absolute;
  right: 0.8rem;
  background: transparent;
  border: none;
  color: rgba(255, 204, 0, 0.7);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0.2rem;
}

.search-button:hover {
  color: #ffcc00;
}


/* 移动端菜单按钮 */
.mobile-menu-button {
  display: none;
  background: transparent;
  border: none;
  cursor: pointer;
  padding: 0.5rem;
  margin-left: 1rem;
}

.bar {
  display: block;
  width: 25px;
  height: 2px;
  margin: 5px 0;
  background: #ffcc00;
  transition: all 0.3s ease;
}

/* 移动端样式 */
@media (max-width: 768px) {
  .navbar {
    padding: 1rem;
    flex-direction: column;
    align-items: flex-start;
  }
  
  .navbar-brand {
    width: 100%;
    justify-content: space-between;
  }
  
  .mobile-menu-button {
    display: block;
  }
  
  .navbar-links {
    width: 100%;
    flex-direction: column;
    gap: 0;
    max-height: 0;
    overflow: hidden;
    transition: max-height 0.3s ease-out;
  }
  
  .mobile-menu-open .navbar-links {
    max-height: 500px;
    padding-top: 1rem;
  }
  
  .nav-link {
    padding: 1rem 0;
    border-bottom: 1px solid rgba(255, 204, 0, 0.1);
  }
  
  .search-container {
    display: none; /* 移动端隐藏搜索框 */
  }
  
  /* 汉堡菜单动画 */
  .bar-1 {
    transform: translateY(7px) rotate(45deg);
  }
  
  .bar-2 {
    opacity: 0;
  }
  
  .bar-3 {
    transform: translateY(-7px) rotate(-45deg);
  }
}
</style>