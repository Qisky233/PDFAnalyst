<template>
  <div class="search-result-container">
    <!-- 星空背景画布 -->
    <div class="star-map" ref="starMapContainer"></div>

    <nav class="navigation">
      <router-link to="/poet" class="back-btn">
        返回星图
      </router-link>
    </nav>

    <!-- 搜索框 -->
    <div class="search-box">
      <input 
        v-model="searchKeyword" 
        @keyup.enter="performSearch"
        placeholder="输入诗人姓名或诗句..."
      />
      <button @click="performSearch">搜索</button>
    </div>
    
    <!-- 主要内容 -->
    <div class="search-content">
      <!-- 作者搜索结果 -->
      <div v-if="authorResults.length > 0" class="result-section">
        <h2 class="section-title">诗人</h2>
        <div class="author-grid">
          <div 
            v-for="author in authorResults" 
            :key="author.author_id"
            class="author-card"
            @click="viewAuthorDetail(author)"
          >
            <div class="author-avatar">
              <img 
                v-if="author.imgUrl" 
                :src="author.imgUrl" 
                :alt="author.name"
                @error="handleImageError(author)"
              />
              <div v-else class="no-image">
                <i class="icon-user"></i>
              </div>
            </div>
            <div class="author-name">{{ author.name }}</div>
          </div>
        </div>
      </div>
      
      <!-- 诗作搜索结果 -->
      <div v-if="poemResults.length > 0" class="result-section">
        <h2 class="section-title">诗作</h2>
        <div class="poem-list">
          <div 
            v-for="poem in poemResults" 
            :key="poem.poem_id"
            class="poem-item"
            @click="showPoemDetail(poem)"
          >
            <div class="poem-title">{{ poem.title }}</div>
            <!-- <div class="poem-author">{{ getAuthorName(poem.author_id) }}</div> -->
            <div class="poem-excerpt">{{ getFirstLine(poem.content) }}</div>
          </div>
        </div>


        <!-- 新增分页组件 -->
        <div class="pagination-container">
            <div class="pagination">
                <button 
                    class="page-button" 
                    :class="{ disabled: currentPage === 1 }"
                    @click="goToPage(1)"
                >
                    首页
                </button>
                <button 
                    class="page-button" 
                    :class="{ disabled: currentPage === 1 }"
                    @click="goToPage(currentPage - 1)"
                >
                    &lt; 上一页
                </button>
                
                <template v-for="page in visiblePages" :key="page">
                    <button 
                    class="page-button"
                    :class="{ active: page === currentPage }"
                    @click="goToPage(page)"
                    >
                    {{ page }}
                    </button>
                </template>
                
                <button 
                    class="page-button" 
                    :class="{ disabled: currentPage === totalPages }"
                    @click="goToPage(currentPage + 1)"
                >
                    下一页 &gt;
                </button>
                <button 
                    class="page-button" 
                    :class="{ disabled: currentPage === totalPages }"
                    @click="goToPage(totalPages)"
                >
                    尾页
                </button>
                </div>
                <div class="page-info">
                共 {{ totalItems }} 条结果，第 {{ currentPage }} 页 / 共 {{ totalPages }} 页
                </div>
            </div>
        </div>
      
      <!-- 无结果提示 -->
      <div v-if="authorResults.length === 0 && poemResults.length === 0 && searched" class="no-results">
        <i class="icon-search"></i>
        <p>没有找到与"{{ searchKeyword }}"相关的结果</p>
      </div>
    </div>
    
    <!-- 诗词弹窗 -->
    <div v-if="selectedPoem" class="poem-modal">
      <div class="modal-content">
        <h3>{{ selectedPoem.title }}</h3>
        <!-- <div class="poem-author">{{ getAuthorName(selectedPoem.author_id) }}</div> -->
        <pre class="poem-content">{{ selectedPoem.content }}</pre>
        <button class="close-btn" @click="selectedPoem = null">关闭</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import * as THREE from 'three';
import { searchAuthor, searchPoem } from '@/api/search';

// 路由控制
const route = useRoute();
const router = useRouter();

// 监听路由变化
watch(
  () => route.query.page,
  (newPage, oldPage) => {
    if (newPage !== oldPage) {
      performSearch();
    }
  }
);

// 搜索相关状态
const searchKeyword = ref(route.query.keyword);
const authorResults = ref([]);
const poemResults = ref([]);
const searched = ref(false);
const selectedPoem = ref(null);
const totalItems = ref(0);
const pageSize = ref(6); // 与API返回的page_size一致

// 计算属性
const currentPage = computed(() => route.query.page ? parseInt(route.query.page) : 1);
const totalPages = computed(() => Math.ceil(totalItems.value / pageSize.value));

// 可见页码范围 (最多显示5个页码)
const visiblePages = computed(() => {
  const pages = [];
  let start = Math.max(1, currentPage.value - 2);
  let end = Math.min(totalPages.value, start + 4);
  
  // 调整起始位置确保显示5个页码
  if (end - start < 4) {
    start = Math.max(1, end - 4);
  }
  
  for (let i = start; i <= end; i++) {
    pages.push(i);
  }
  return pages;
});

// 跳转到指定页码
const goToPage = (page) => {
  if (page < 1 || page > totalPages.value || page === currentPage.value) return;
  
  router.push({
    path: '/search',
    query: { 
      keyword: searchKeyword.value,
      page 
    }
  });
};

// Three.js 星空背景初始化
const starMapContainer = ref(null);
let scene, camera, renderer;

// 执行搜索
const performSearch = async () => {
  if (!searchKeyword.value.trim()) return;

  const currentPage = route.query.page ? parseInt(route.query.page) : 1;

  // 使用 router.push 更新 URL 参数
  router.push({
    path: '/search',
    query: { 
        keyword: searchKeyword.value,
        page: currentPage
    }
  });

  try {
    // 搜索作者
    const authorsResponse = await searchAuthor(searchKeyword.value);
    authorResults.value = authorsResponse.data || [];
    
    // 搜索诗作
    const poemsResponse = await searchPoem(searchKeyword.value, currentPage);
    poemResults.value = poemsResponse.data || [];
    
    totalItems.value = poemsResponse.total || 0; // 更新总结果数
    searched.value = true;
  } catch (error) {
    console.error('搜索失败:', error);
    authorResults.value = [];
    poemResults.value = [];
  }
};

// 查看作者详情
const viewAuthorDetail = (author) => {
  router.push({
    path: '/detail',
    query: {
      authorId: author.author_id,
      authorName: author.name
    }
  });
};

// 显示诗词详情
const showPoemDetail = (poem) => {
  selectedPoem.value = poem;
};

// 获取作者姓名
const getAuthorName = (authorId) => {
  const author = authorResults.value.find(a => a.author_id === authorId);
  return author ? author.name : '未知作者';
};

// 获取诗句首行
const getFirstLine = (content) => {
  return content.split('\n')[0].slice(0, 20) + (content.split('\n')[0].length > 20 ? '...' : '');
};

// 处理图片加载错误
const handleImageError = (author) => {
  author.imgUrl = '';
};

// 初始化星空背景
const initStarField = () => {
  scene = new THREE.Scene();
  camera = new THREE.PerspectiveCamera(75, window.innerWidth / window.innerHeight, 0.1, 1000);
  camera.position.z = 30;
  
  renderer = new THREE.WebGLRenderer({ alpha: true });
  renderer.setSize(window.innerWidth, window.innerHeight);
  starMapContainer.value.appendChild(renderer.domElement);
  
  const createStars = () => {
    const geometry = new THREE.BufferGeometry();
    const vertices = [];
    for (let i = 0; i < 2000; i++) {
      vertices.push(
        THREE.MathUtils.randFloatSpread(2000),
        THREE.MathUtils.randFloatSpread(2000),
        THREE.MathUtils.randFloatSpread(2000)
      );
    }
    geometry.setAttribute('position', new THREE.Float32BufferAttribute(vertices, 3));
    const material = new THREE.PointsMaterial({ 
      color: 0xaaaaaa,
      size: 0.5,
      transparent: true,
      opacity: 0.8
    });
    const stars = new THREE.Points(geometry, material);
    scene.add(stars);
  };
  
  createStars();
  
  const animate = () => {
    requestAnimationFrame(animate);
    scene.rotation.y += 0.001;
    renderer.render(scene, camera);
  };
  animate();
};

onMounted(() => {
  initStarField();

  performSearch(); // 执行初始搜索
});
</script>

<style scoped>
.search-result-container {
  position: relative;
  width: 100%;
  min-height: 100vh;
  overflow: hidden;
  background: radial-gradient(ellipse at bottom, #1B2735 0%, #090A0F 100%);
  color: white;
  padding-bottom: 50px;
}

.star-map {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 0;
}

.navigation {
  position: fixed;
  top: 2rem;
  left: 2rem;
  z-index: 10;
}

.back-btn {
  display: inline-block;
  padding: 0.8rem 1.5rem;
  background: rgba(255, 255, 255, 0.1);
  color: #ffffff;
  text-decoration: none;
  border-radius: 5px;
  backdrop-filter: blur(10px);
  transition: all 0.3s ease;
}

.back-btn:hover {
  background: rgba(255, 255, 255, 0.2);
}

.search-box {
  position: relative;
  z-index: 10;
  max-width: 800px;
  margin: 2rem auto 0;
  display: flex;
  gap: 10px;
}

.search-box input {
  flex: 1;
  padding: 12px 20px;
  border: none;
  border-radius: 30px;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
  color: white;
  font-size: 16px;
  border: 1px solid rgba(255, 204, 0, 0.3);
}

.search-box input::placeholder {
  color: rgba(255, 255, 255, 0.6);
}

.search-box input:focus {
  outline: none;
  border-color: #ffcc00;
  box-shadow: 0 0 10px rgba(255, 204, 0, 0.3);
}

.search-box button {
  padding: 0 25px;
  background: rgba(255, 204, 0, 0.2);
  border: 1px solid #ffcc00;
  color: #ffcc00;
  border-radius: 30px;
  cursor: pointer;
  transition: all 0.3s;
}

.search-box button:hover {
  background: rgba(255, 204, 0, 0.3);
}

.search-content {
  position: relative;
  z-index: 10;
  max-width: 1200px;
  margin: 2rem auto;
  padding: 0 20px;
}

.result-section {
  margin-bottom: 40px;
}

.section-title {
  color: #ffcc00;
  font-size: 1.5rem;
  margin-bottom: 20px;
  padding-bottom: 10px;
  border-bottom: 1px solid rgba(255, 204, 0, 0.3);
}

/* 作者网格样式 */
.author-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
  gap: 20px;
}

.author-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  cursor: pointer;
  transition: all 0.3s;
}

.author-card:hover {
  transform: translateY(-5px);
}

.author-avatar {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  background: rgba(255, 204, 0, 0.1);
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  border: 2px solid rgba(255, 204, 0, 0.3);
  margin-bottom: 10px;
}

.author-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.no-image {
  color: rgba(255, 204, 0, 0.5);
  font-size: 30px;
}

.author-name {
  text-align: center;
  font-size: 14px;
  color: #ddd;
}

/* 诗作列表样式 */
.poem-list {
  display: grid;
  gap: 15px;
}

.poem-item {
  background: rgba(0, 0, 0, 0.5);
  border-radius: 8px;
  padding: 15px;
  cursor: pointer;
  transition: all 0.3s;
  border-left: 3px solid transparent;
}

.poem-item:hover {
  background: rgba(255, 204, 0, 0.1);
  border-left-color: #ffcc00;
}

.poem-title {
  font-size: 18px;
  color: #ffcc00;
  margin-bottom: 5px;
}

.poem-author {
  font-size: 14px;
  color: #aaa;
  margin-bottom: 10px;
}

.poem-excerpt {
  color: #ddd;
  font-size: 15px;
  line-height: 1.6;
}

/* 无结果提示 */
.no-results {
  text-align: center;
  padding: 50px 0;
  color: rgba(255, 255, 255, 0.6);
}

.no-results i {
  font-size: 50px;
  color: rgba(255, 204, 0, 0.3);
  margin-bottom: 20px;
}

/* 诗词弹窗样式 */
.poem-modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.8);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 100;
}

.modal-content {
  background: rgba(0, 0, 0, 0.9);
  padding: 30px;
  border-radius: 8px;
  max-width: 600px;
  width: 90%;
  border: 1px solid #ffcc00;
  position: relative;
}

.modal-content h3 {
  color: #ffcc00;
  margin-top: 0;
  font-size: 24px;
}

.poem-author {
  color: #aaa;
  margin-bottom: 20px;
  font-size: 16px;
}

.poem-content {
  white-space: pre-wrap;
  font-family: '宋体', serif;
  line-height: 2;
  color: #eee;
  font-size: 18px;
  max-height: 80vh;
  overflow-y: auto;
}

.close-btn {
  position: absolute;
  top: 15px;
  right: 15px;
  background: transparent;
  color: #ffcc00;
  border: none;
  cursor: pointer;
  font-size: 20px;
}

/* 新增分页样式 */
.pagination-container {
  margin-top: 30px;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.pagination {
  display: flex;
  gap: 5px;
  flex-wrap: wrap;
  justify-content: center;
}

.page-button {
  min-width: 36px;
  height: 36px;
  padding: 0 12px;
  background: rgba(255, 204, 0, 0.1);
  border: 1px solid rgba(255, 204, 0, 0.3);
  color: #ffcc00;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.3s;
  font-size: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.page-button:hover:not(.disabled):not(.active) {
  background: rgba(255, 204, 0, 0.2);
  border-color: rgba(255, 204, 0, 0.5);
}

.page-button.active {
  background: rgba(255, 204, 0, 0.3);
  border-color: #ffcc00;
  color: #fff;
  font-weight: bold;
}

.page-button.disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.page-info {
  margin-top: 15px;
  color: rgba(255, 204, 0, 0.8);
  font-size: 14px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .author-grid {
    grid-template-columns: repeat(auto-fill, minmax(80px, 1fr));
    gap: 15px;
  }
  
  .author-avatar {
    width: 60px;
    height: 60px;
  }
  
  .search-box {
    margin: 1rem 15px 0;
  }
  
  .search-content {
    padding: 0 15px;
  }
  .page-button {
    min-width: 30px;
    height: 30px;
    padding: 0 8px;
    font-size: 12px;
  }
  
  .page-info {
    font-size: 12px;
  }
}
</style>