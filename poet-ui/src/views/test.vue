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
import { ref, onMounted, computed, onUnmounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import * as THREE from 'three';
import { getAuthorById } from '@/api/poet';
import { searchAllPoems } from '@/api/search'

// 路由控制
const route = useRoute();
const router = useRouter();
const goBack = () => router.push('/poet');

// 数据获取
const poet = ref({ 
  name: '',
  description: '',
  poems: [] 
});
const selectedPoem = ref(null);

const fetchPoetData = async () => {
  try {
    const response = await getAuthorById(route.query.authorId);
    poet.value = response;

    const poemsResponse = await searchAllPoems(route.query.authorId);
    poet.value.poems = poemsResponse.data || [];

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
}

.poem-list {
  margin-top: 1rem;
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

/* 响应式设计 */
@media (max-width: 768px) {
  .poet-card {
    margin: 1rem;
    padding: 1.5rem;
  }
  
  .poet-name {
    font-size: 1.8rem;
  }
}
</style>