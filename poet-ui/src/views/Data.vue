<template>
  <div class="data-overview-container">
    <!-- 星空背景画布 -->
    <div class="star-map" ref="starMapContainer"></div>
    <navBar />
    <!-- 主内容区域 -->
    <div class="data-content">
      <!-- 标题 -->
      <h1 class="page-title">诗词数据一览<small>（仅收录）</small></h1>

            <!-- 新增：数据概览卡片 -->
            <div class="stats-cards">
        <div class="stat-card">
          <div class="stat-icon">
            <el-icon><User /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-number">{{ formatNumber(stats.poets) }}</div>
            <div class="stat-label">收录诗人</div>
          </div>
        </div>
        
        <div class="stat-card">
          <div class="stat-icon">
            <el-icon><Notebook /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-number">{{ formatNumber(stats.poems) }}</div>
            <div class="stat-label">收录古诗</div>
          </div>
        </div>
        
        <div class="stat-card">
          <div class="stat-icon">
            <el-icon><EditPen /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-number">{{ formatNumber(stats.words) }}</div>
            <div class="stat-label">总字数</div>
          </div>
        </div>
      </div>
      
      <!-- 双图表部分 -->
      <div class="double-charts">
        <div class="section-card">
          <h2 class="section-title">
            <el-icon><PieChart /></el-icon>
            <span>句式分布</span>
          </h2>
          <div class="chart-container" ref="dynastyChart"></div>
        </div>
        
        <div class="section-card">
          <h2 class="section-title">
            <el-icon><Histogram /></el-icon>
            <span>作品数量TOP10</span>
          </h2>
          <div class="chart-container" ref="worksChart"></div>
        </div>
      </div>
      
      <!-- 词云部分 -->
      <div class="section-card">
        <h2 class="section-title">
          <el-icon><Collection /></el-icon>
          <span>诗人词云</span>
        </h2>
        <div class="word-cloud" ref="wordCloudChart"></div>
      </div>
      
      <!-- 表格部分 -->
      <div class="section-card">
        <h2 class="section-title">
          <el-icon><Grid /></el-icon>
          <span>诗人数据总览</span>
        </h2>
        <el-table
          :data="tableData"
          style="width: 100%"
          class="transparent-table"
          :header-cell-style="tableHeaderStyle"
          :row-style="tableRowStyle"
          @row-click="handleRowClick"
        >
          <el-table-column prop="author_id" label="id" />
          <el-table-column prop="author_name" label="诗人" />
          <el-table-column prop="dynasty" label="朝代" />
          <el-table-column prop="poem_count" label="作品数" />
          <el-table-column prop="word_count" label="总字数" />
          <!-- <el-table-column prop="popularity" label="传唱度">
            <template #default="{ row }">
              <el-progress 
                :percentage="row.popularity" 
                :color="customColors" 
                :show-text="false"
              />
            </template>
          </el-table-column> -->
          <el-table-column label="操作" width="120">
            <template #default="{ row }">
              <el-button 
                link
                @click.stop="viewPoetDetail(row)"
                class="detail-btn"
              >
                查看详情
              </el-button>
            </template>
          </el-table-column>
        </el-table>
        
        <div class="pagination-container">
          <el-pagination
            v-model:current-page="currentPage"
            :page-size="pageSize"
            layout="prev, pager, next, jumper"
            :total="total"
            @current-change="handlePageChange"
            class="custom-pagination"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, nextTick  } from 'vue';
import { useRouter } from 'vue-router';
import * as echarts from 'echarts';
import 'echarts-wordcloud';
import * as THREE from 'three';

import { getDataStats, getEchartData, getTableData } from '@/api/data.js';

import { 
  User, 
  Notebook, 
  EditPen, 
  Collection, 
  PieChart, 
  Histogram, 
  Grid 
} from '@element-plus/icons-vue';


import navBar from '@/components/navBar.vue';

const router = useRouter();

// 新增：统计数据
const stats = ref({
  poets: 0,
  poems: 0,
  words: 0
});

// 获取统计数据
const fetchStats = async () => {
  try {
    const response = await getDataStats();
    console.log('获取统计数据:', response.data);
    stats.value = response.data;
  } catch (error) {
    console.error('获取统计数据失败:', error);
  }
};

// 数字格式化
const formatNumber = (num) => {
  return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",");
};

// 星空背景初始化
const starMapContainer = ref(null);
let scene, camera, renderer;

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

// 图表相关
const wordCloudChart = ref(null);
const dynastyChart = ref(null);
const worksChart = ref(null);

// 表格数据
const tableData = ref([
  { rank: 1, name: '李白', dynasty: '唐', poemCount: 1010, wordCount: 24560, popularity: 95 },
  { rank: 2, name: '杜甫', dynasty: '唐', poemCount: 980, wordCount: 23890, popularity: 92 },
  { rank: 3, name: '白居易', dynasty: '唐', poemCount: 890, wordCount: 21560, popularity: 88 },
  { rank: 4, name: '苏轼', dynasty: '宋', poemCount: 780, wordCount: 19870, popularity: 85 },
  { rank: 5, name: '李清照', dynasty: '宋', poemCount: 720, wordCount: 18540, popularity: 82 },
  { rank: 6, name: '辛弃疾', dynasty: '宋', poemCount: 680, wordCount: 17650, popularity: 80 },
  { rank: 7, name: '王维', dynasty: '唐', poemCount: 650, wordCount: 16580, popularity: 78 },
  { rank: 8, name: '李商隐', dynasty: '唐', poemCount: 620, wordCount: 15890, popularity: 75 },
  { rank: 9, name: '杜牧', dynasty: '唐', poemCount: 590, wordCount: 15230, popularity: 72 },
  { rank: 10, name: '陆游', dynasty: '宋', poemCount: 560, wordCount: 14870, popularity: 70 }
]);

// 分页
const currentPage = ref(1);
const pageSize = ref(10);
const total = ref(100);

// 自定义颜色
const customColors = [
  { color: '#ffcc00', percentage: 20 },
  { color: '#ff9900', percentage: 40 },
  { color: '#ff6600', percentage: 60 },
  { color: '#ff3300', percentage: 80 },
  { color: '#ff0000', percentage: 100 }
];

// 表格样式
const tableHeaderStyle = {
  backgroundColor: 'rgba(255, 204, 0, 0.1)',
  color: '#ffcc00',
  fontSize: '14px'
}

const tableRowStyle = {
  backgroundColor: 'rgba(0, 0, 0, 0.3)',
  color: '#ddd',
  cursor: 'pointer',
  transition: 'all 0.3s'
}

// 初始化图表
const initCharts = () => {
  // 词云图
  const wordCloud = echarts.init(wordCloudChart.value);
  wordCloud.setOption({
    series: [{
      type: 'wordCloud',
      shape: 'circle',
      left: 'center',
      top: 'center',
      width: '100%',
      height: '100%',
      right: null,
      bottom: null,
      sizeRange: [12, 50],
      rotationRange: [-45, 45],
      rotationStep: 15,
      gridSize: 8,
      drawOutOfBound: false,
      textStyle: {
        fontFamily: 'sans-serif',
        fontWeight: 'bold',
        color: function () {
          return `rgb(${[
            Math.round(Math.random() * 155 + 100),
            Math.round(Math.random() * 155 + 100),
            Math.round(Math.random() * 55)
          ].join(',')})`;
        }
      },
      emphasis: {
        textStyle: {
          shadowBlur: 10,
          shadowColor: '#ffcc00'
        }
      },
      data: []
    }]
  });

  // 句式分布图
  const dynasty = echarts.init(dynastyChart.value);
  dynasty.setOption({
    tooltip: {
      trigger: 'item',
      formatter: '{a} <br/>{b}: {c} ({d}%)'
    },
    legend: {
      orient: 'vertical',
      right: 10,
      top: 'center',
      textStyle: {
        color: '#ddd'
      }
    },
    series: [
      {
        name: '句式分布',
        type: 'pie',
        radius: ['40%', '70%'],
        avoidLabelOverlap: false,
        itemStyle: {
          borderRadius: 10,
          borderColor: '#1B2735',
          borderWidth: 2
        },
        label: {
          show: false,
          position: 'center'
        },
        emphasis: {
          label: {
            show: true,
            fontSize: '18',
            fontWeight: 'bold',
            color: '#ffcc00'
          }
        },
        labelLine: {
          show: false
        },
        data: [
          { value: 46086, name: '五言' },
          { value: 46086, name: '七言' },
          { value: 5761, name: '四言' },
          { value: 9217, name: '杂言' },
          { value: 3456, name: '六言' },
          { value: 3456, name: '排律' },
          { value: 1152, name: '三言' }
        ],
        color: ['#ffcc00', '#ff9900', '#ff6600', '#ff3300', '#ff0000', '#cc0000']
      }
    ]
  });

  // 作品数量TOP10
  const works = echarts.init(worksChart.value);
  works.setOption({
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      }
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'value',
      axisLine: {
        lineStyle: {
          color: '#666'
        }
      },
      axisLabel: {
        color: '#999'
      },
      splitLine: {
        lineStyle: {
          color: 'rgba(255, 255, 255, 0.1)'
        }
      }
    },
    yAxis: {
      type: 'category',
      data: [],
      axisLine: {
        lineStyle: {
          color: '#666'
        }
      },
      axisLabel: {
        color: '#ddd'
      },
      inverse: true
    },
    series: [
      {
        name: '作品数量',
        type: 'bar',
        data: [],
        itemStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 1, 0, [
            { offset: 0, color: '#ffcc00' },
            { offset: 1, color: '#ff6600' }
          ]),
          borderRadius: [0, 4, 4, 0]
        },
        label: {
          show: true,
          position: 'right',
          color: '#ffcc00'
        }
      }
    ]
  });

  // 窗口大小变化时重绘图表
  const handleResize = () => {
    wordCloud.resize();
    dynasty.resize();
    works.resize();
    renderer.setSize(window.innerWidth, window.innerHeight);
    camera.aspect = window.innerWidth / window.innerHeight;
    camera.updateProjectionMatrix();
  };

  window.addEventListener('resize', handleResize);
  onBeforeUnmount(() => {
    window.removeEventListener('resize', handleResize);
    wordCloud.dispose();
    dynasty.dispose();
    works.dispose();
    renderer.dispose();
  });
};

// 表格行点击
const handleRowClick = (row) => {
  console.log('点击行:', row);
};

// 查看诗人详情
const viewPoetDetail = (row) => {
  // router.push({
  //   path: '/poet/detail',
  //   query: {
  //     authorId: row.id,
  //     authorName: row.name
  //   }
  // });
  router.push({
    path: '/detail',
    query: {
      authorId: row.author_id,
      authorName: row.author_name
    }
  });
};

// 分页变化
const handlePageChange = (val) => {
  currentPage.value = val;
  fetchTableData();
};

// 获取作品数量TOP10数据
const fetchTop10Data = async () => {
  try {
    const response = await getEchartData("two");
    const sortedData = response.data.authors.sort((a, b) => b.poem_count - a.poem_count).slice(0, 10);
    const names = sortedData.map(item => item.author_name);
    const counts = sortedData.map(item => item.poem_count);

    const works = echarts.init(worksChart.value);
    works.setOption({
      yAxis: {
        data: names
      },
      series: [
        {
          data: counts
        }
      ]
    });
  } catch (error) {
    console.error('获取作品数量TOP10数据失败:', error);
  }
};

// 获取诗人数据总览
const fetchTableData = async () => {
  try {
    const response = await getTableData();
    const list = response.data.list;
    total.value = list.length;

    // 分页逻辑
    const start = (currentPage.value - 1) * pageSize.value;
    const end = start + pageSize.value;
    tableData.value = list.slice(start, end);
  } catch (error) {
    console.error('获取诗人数据总览失败:', error);
  }
};

// 生成词云
const fetchWordCloudData = async () => {
  try {
    const response = await getEchartData("two");
    const sortedData = response.data.authors.sort((a, b) => b.poem_count - a.poem_count).slice(0, 300);
    const wordCloudData = sortedData.map(item => ({
      name: item.author_name,
      value: item.poem_count
    }));

    const wordCloud = echarts.init(wordCloudChart.value);
    wordCloud.setOption({
      series: [{
        type: 'wordCloud',
        shape: 'circle',
        left: 'center',
        top: 'center',
        width: '100%',
        height: '100%',
        right: null,
        bottom: null,
        sizeRange: [12, 50],
        rotationRange: [-45, 45],
        rotationStep: 15,
        gridSize: 8,
        drawOutOfBound: false,
        textStyle: {
          fontFamily: 'sans-serif',
          fontWeight: 'bold',
          color: function () {
            return `rgb(${[
              Math.round(Math.random() * 155 + 100),
              Math.round(Math.random() * 155 + 100),
              Math.round(Math.random() * 55)
            ].join(',')})`;
          }
        },
        emphasis: {
          textStyle: {
            shadowBlur: 10,
            shadowColor: '#ffcc00'
          }
        },
        data: wordCloudData
      }]
    });
  } catch (error) {
    console.error('获取词云数据失败:', error);
  }
};

onMounted(() => {
  initStarField();
  initCharts();

  fetchStats();
  fetchTop10Data();
  fetchTableData();
  fetchWordCloudData();
});
</script>

<style scoped>
.data-overview-container {
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

.data-content {
  position: relative;
  z-index: 10;
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
  margin-top: 60px;
}

.page-title {
  color: #ffcc00;
  text-align: center;
  margin-bottom: 30px;
  font-size: 2.2rem;
  text-shadow: 0 0 10px rgba(255, 204, 0, 0.3);
}

.section-card {
  background: rgba(10, 15, 20, 0.7);
  backdrop-filter: blur(10px);
  border-radius: 12px;
  padding: 20px;
  margin-bottom: 30px;
  border: 1px solid rgba(255, 204, 0, 0.2);
  box-shadow: 0 0 20px rgba(0, 0, 0, 0.5);
}

.section-title {
  color: #ffcc00;
  font-size: 1.3rem;
  margin-bottom: 15px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.word-cloud {
  width: 100%;
  height: 400px;
}

.double-charts {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}

.chart-container {
  width: 100%;
  height: 350px;
}

.transparent-table {
  background: transparent;
  --el-table-border-color: rgba(255, 204, 0, 0.2);
  --el-table-header-bg-color: rgba(255, 204, 0, 0.1);
  --el-table-tr-bg-color: rgba(0, 0, 0, 0.3);
}

.transparent-table::before {
  background-color: rgba(255, 204, 0, 0.2);
}

.transparent-table :deep(.el-table__body tr:hover td) {
  background-color: rgba(255, 204, 0, 0.15) !important;
}

.detail-btn {
  color: #ffcc00;
}

.detail-btn:hover {
  color: #ff9900;
}

.custom-pagination {
  margin-top: 20px;
  justify-content: center;
}

.custom-pagination :deep(.btn-prev),
.custom-pagination :deep(.btn-next),
.custom-pagination :deep(.number) {
  background: rgba(255, 204, 0, 0.1);
  color: #ffcc00;
  border: 1px solid rgba(255, 204, 0, 0.3);
  margin: 0 5px;
}

.custom-pagination :deep(.number.active) {
  background: rgba(255, 204, 0, 0.3);
  border-color: #ffcc00;
  color: #fff;
  font-weight: bold;
}

.custom-pagination :deep(.number:hover:not(.active)) {
  background: rgba(255, 204, 0, 0.2);
  border-color: rgba(255, 204, 0, 0.5);
}

/* 新增样式 */
.stats-cards {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
  margin-bottom: 30px;
}

.stat-card {
  background: rgba(255, 204, 0, 0.1);
  border: 1px solid rgba(255, 204, 0, 0.3);
  border-radius: 12px;
  padding: 20px;
  display: flex;
  align-items: center;
  transition: all 0.3s;
  backdrop-filter: blur(5px);
}

.stat-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 10px 20px rgba(0, 0, 0, 0.3);
  border-color: #ffcc00;
}

.stat-icon {
  width: 60px;
  height: 60px;
  background: rgba(255, 204, 0, 0.2);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 20px;
  flex-shrink: 0;
}

.stat-icon .el-icon {
  font-size: 30px;
  color: #ffcc00;
}

.stat-content {
  flex: 1;
}

.stat-number {
  font-size: 2rem;
  font-weight: bold;
  color: #ffcc00;
  margin-bottom: 5px;
}

.stat-label {
  font-size: 1rem;
  color: #aaa;
}

/* 响应式设计 */
@media (max-width: 992px) {
  .double-charts {
    grid-template-columns: 1fr;
  }
  
  .word-cloud {
    height: 300px;
  }
  
  .chart-container {
    height: 300px;
  }
}

@media (max-width: 768px) {
  .data-content {
    padding: 15px;
  }
  
  .page-title {
    font-size: 1.8rem;
  }
  
  .section-card {
    padding: 15px;
  }
  
  .word-cloud {
    height: 250px;
  }

  .stats-cards {
    grid-template-columns: 1fr;
  }
  
  .stat-card {
    padding: 15px;
  }
  
  .stat-icon {
    width: 50px;
    height: 50px;
    margin-right: 15px;
  }
  
  .stat-icon .el-icon {
    font-size: 24px;
  }
  
  .stat-number {
    font-size: 1.5rem;
  }
}
</style>