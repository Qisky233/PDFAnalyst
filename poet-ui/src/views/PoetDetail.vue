<template>
  <div class="poet-detail-container">
    <!-- 星空背景画布 -->
    <div class="star-map" ref="starMapContainer"></div>
    
    <!-- 诗人信息主卡 -->
    <div class="poet-card">
      <!-- 标题区 -->
      <div class="poet-header">
        <h1 class="poet-name">{{ route.query.authorName }}</h1>
        <button class="back-btn" @click="goBack">返回星图</button>
      </div>
      
      <!-- 元数据 -->
      <div class="poet-meta">
        <span class="meta-item">朝代：{{ poet.dynasty || '唐朝' }}</span>
        <span class="meta-item">诗作：{{ poet.total_poems || 0 }}首</span>
      </div>
      
      <!-- 传记 -->
      <div class="poet-bio">
        <h2 class="section-title">诗人简介</h2>
        <p class="bio-content">{{ poet.description || '暂无详细传记资料' }}</p>
      </div>
      
      <!-- 代表作 -->
      <div class="poet-works">
        <h2 class="section-title">代表诗作</h2>
        <div class="poem-list">
          <div 
            v-for="(poem, index) in featuredPoems" 
            :key="poem.poem_id" 
            class="poem-item"
            @click="showPoemDetail(poem)"
          >
            <span class="poem-index">{{ index + 1 }}.</span>
            <span class="poem-title">{{ poem.title }}</span>
            <span class="poem-excerpt">{{ getFirstLine(poem.content) }}</span>
          </div>
        </div>
      </div>

      <!-- 分页组件 -->
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
          第 {{ currentPage }} 页 / 共 {{ totalPages }} 页
        </div>
      </div>


    </div>
    
    <!-- 诗词弹窗 -->
    <div v-if="selectedPoem" class="poem-modal">
      <div class="modal-content">
        <h3>{{ selectedPoem.title }}</h3>
        <pre class="poem-content">{{ selectedPoem.content }}</pre>
        <button class="close-btn" @click="selectedPoem = null">关闭</button>
      </div>
    </div>
  </div>
</template>
<script setup>
import { ref, onMounted, computed, onUnmounted, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import * as THREE from 'three';
import { getAuthorById } from '@/api/poet';
import { searchAllPoems } from '@/api/search'

// 路由控制
const route = useRoute();
const router = useRouter();

// 监听路由变化
watch(
  () => route.query.page,
  (newPage, oldPage) => {
    if (newPage !== oldPage) {
      fetchPoetData();
    }
  }
);
// 跳转到指定页码
const goToPage = (page) => {
  if (page < 1 || page > totalPages.value || page === currentPage.value) return;
  
  router.push({
    path: '/detail',
    query: { 
      authorId: route.query.authorId,
      authorName: route.query.authorName,
      page 
    }
  });
};
const goBack = () => router.push('/poet');

// 分页相关状态
const pageSize = ref(6); // 与API返回的page_size一致
const totalPoems = ref(0); // 总诗词数量

// 数据获取
const poet = ref({ 
  name: '',
  description: '',
  poems: [] 
});
const selectedPoem = ref(null);

// 计算属性
const poems = computed(() => poet.value.poems || []);
const currentPage = computed(() => route.query.page ? parseInt(route.query.page) : 1);
const totalPages = computed(() => Math.ceil(totalPoems.value / pageSize.value));

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

const fetchPoetData = async () => {
  try {
    const response = await getAuthorById(route.query.authorId);
    poet.value = response;

    const poemsResponse = await searchAllPoems(route.query.authorId, route.query.page);
    poet.value.poems = poemsResponse.data || [];
    totalPoems.value = poemsResponse.total || 0;

    poet.value.total_poems = poemsResponse.total || 0;
  } catch (error) {
    console.error('获取诗人详情失败:', error);
  }
};

// 计算属性
const featuredPoems = computed(() => {
  return poet.value.poems?.slice(0, 5) || [];
});

// 诗词处理
const getFirstLine = (content) => {
  return content.split('\n')[0].slice(0, 20) + '...';
};

const showPoemDetail = (poem) => {
  selectedPoem.value = poem;
};

// Three.js 星空背景初始化
const starMapContainer = ref(null);
let scene, camera, renderer;

const initStarField = () => {
  // 场景设置
  scene = new THREE.Scene();
  camera = new THREE.PerspectiveCamera(75, window.innerWidth / window.innerHeight, 0.1, 1000);
  camera.position.z = 30;
  
  // 渲染器
  renderer = new THREE.WebGLRenderer({ alpha: true });
  renderer.setSize(window.innerWidth, window.innerHeight);
  starMapContainer.value.appendChild(renderer.domElement);
  
  // 创建星星
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
  
  // 动画循环
  const animate = () => {
    requestAnimationFrame(animate);
    scene.rotation.y += 0.001;
    renderer.render(scene, camera);
  };
  animate();
};

onMounted(() => {
  fetchPoetData();
  initStarField();
});

onUnmounted(() => {
  if (renderer) {
    renderer.dispose();
  }
});
</script>
<style scoped>
.poet-detail-container {
  position: relative;
  width: 100%;
  height: 100vh;
  overflow: hidden;
  background: radial-gradient(ellipse at bottom, #1B2735 0%, #090A0F 100%);
}

.star-map {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 0;
}

.poet-card {
  position: relative;
  z-index: 10;
  max-width: 800px;
  margin: 2rem auto;
  padding: 2rem;
  background: rgba(0, 0, 0, 0.7);
  border-radius: 16px;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 204, 0, 0.2);
  box-shadow: 0 0 20px rgba(255, 204, 0, 0.1);
}

.poet-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.poet-name {
  color: #ffcc00;
  font-size: 2.2rem;
  margin: 0;
  text-shadow: 0 0 10px rgba(255, 204, 0, 0.3);
}

.back-btn {
  background: transparent;
  color: #ffcc00;
  border: 1px solid #ffcc00;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.3s;
}

.back-btn:hover {
  background: rgba(255, 204, 0, 0.2);
}

.poet-meta {
  display: flex;
  gap: 1.5rem;
  margin-bottom: 2rem;
  color: #aaa;
}

.section-title {
  color: #ffcc00;
  border-bottom: 1px solid rgba(255, 204, 0, 0.3);
  padding-bottom: 0.5rem;
  font-size: 1.3rem;
}

.bio-content {
  line-height: 1.8;
  color: #ddd;
  max-height: 120px;
  overflow-y: auto;
}

.poem-list {
  margin-top: 1rem;
  max-height: 42vh;
  overflow-y: auto;
}

.poem-item {
  padding: 1rem;
  margin: 0.5rem 0;
  background: rgba(255, 255, 255, 0.05);
  border-left: 3px solid transparent;
  cursor: pointer;
  transition: all 0.3s;
  display: grid;
  grid-template-columns: 30px 1fr;
  gap: 1rem;
}

.poem-item:hover {
  background: rgba(255, 204, 0, 0.1);
  border-left-color: #ffcc00;
}

.poem-index {
  color: #ffcc00;
}

.poem-title {
  font-weight: bold;
  color: white;
}

.poem-excerpt {
  grid-column: 2;
  color: #aaa;
  font-size: 0.9rem;
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
  padding: 2rem;
  border-radius: 8px;
  max-width: 600px;
  width: 90%;
  border: 1px solid #ffcc00;
  position: relative;
}

.modal-content h3 {
  color: #ffcc00;
  margin-top: 0;
}

.poem-content {
  white-space: pre-wrap;
  font-family: '宋体', serif;
  line-height: 2;
  color: #eee;
  max-height: 80vh;
  overflow-y: auto;
}

.close-btn {
  position: absolute;
  top: 1rem;
  right: 1rem;
  background: transparent;
  color: #ffcc00;
  border: none;
  cursor: pointer;
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
  .poet-card {
    margin: 1rem;
    padding: 1.5rem;
  }
  
  .poet-name {
    font-size: 1.8rem;
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