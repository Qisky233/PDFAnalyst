<template>
  <div class="dashboard">
    <canvas ref="canvasRef" class="particle-canvas"></canvas>

    <div class="container">
      <div id="overview" class="intro">
        <h2>基于PDF的论文分析平台<small><br>（中医药方向）</small></h2>
        <div class="intro-content">
          <p>简单、便捷的中医药论文助手</p>
          <p>
            面向中医药领域复杂的术语体系，如多义词、药性、方剂配伍等，本平台集成常用的中医药数据库，帮助您快速解析论文核心内容，挖掘研究热点与趋势，提升科研效率。
          </p>
        </div>
      </div>

      <el-button @click="go('/papers')" type="primary" class="use-btn">立即使用</el-button>

      <!-- 数据卡片容器 -->
      <div id="charts" class="grid-container">
        <!-- 论文统计 -->
        <div class="card glow echart1" @mouseenter="handleCardHover">
          <h3>分析次数</h3>
          <div ref="chart1" class="chart"></div>
        </div>

        <!-- 词云 -->
        <div class="card glow echart2" @mouseenter="handleCardHover">
          <h3>高频中药</h3>
          <div ref="chart2" class="chart"></div>
        </div>

        <!-- R图 -->
        <div class="card glow echart3" @mouseenter="handleCardHover">
          <h3>R图</h3>
          <div ref="chart3" class="chart"></div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue';
import { initParticles, destroyParticles } from './particles';
import { initCharts, destroyCharts, handleResize } from './charts';
import { useRouter } from 'vue-router';

const router = useRouter();

const go = (path) => {
  router.push(path)
}

const canvasRef = ref(null);
const chart1 = ref(null);
const chart2 = ref(null);
const chart3 = ref(null);

onMounted(() => {
  initParticles(canvasRef.value);
  initCharts(chart1, chart2, chart3);
  window.addEventListener('resize', handleResize);
});

onUnmounted(() => {
  destroyParticles();
  destroyCharts(chart1, chart2, chart3);
  window.removeEventListener('resize', handleResize);
});
</script>

<style scoped>
/* 全局样式 */
.dashboard {
  position: relative;
  min-height: 100vh;
  background: radial-gradient(circle at center, #0a1630 0%, #020714 100%);
  color: #fff;
  overflow: hidden;
}

.particle-canvas {
  position: absolute;
  top: 0;
  left: 0;
  z-index: 0;
}

.container {
  position: relative;
  z-index: 1;
  padding: 2rem;
  max-width: 1400px;
  margin: 0 auto;
}

/* 说明部分 */
.intro {
  margin: 5rem 0 1rem 0;
  text-align: center;
}

.intro h2 {
  font-size: 2rem;
  margin-bottom: 1rem;
  color: #00f7ff;
}

.intro small {
  font-size: 1rem;
  color: #aaa;
}

.intro-content p {
  font-size: 1rem;
  line-height: 1.6;
  color: #ddd;
}

/* 立即使用按钮 */
.use-btn {
  position: relative;
  background: linear-gradient(135deg, #6ab7ff, #4a90e2);
  border: none;
  left: calc(50% - 56px);
  font-size: 1rem;
  padding: 0.75rem 1.5rem;
  border-radius: 30px;
  color: #fff;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 1px;
  box-shadow: 0 4px 15px rgba(106, 183, 255, 0.3);
  cursor: pointer;
  transition: all 0.3s ease;
}

.use-btn:hover {
  background: linear-gradient(135deg, #4a90e2, #6ab7ff);
  box-shadow: 0 6px 20px rgba(106, 183, 255, 0.5);
  transform: translateY(-2px);
}

/* 数据卡片容器 */
.grid-container {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 2rem;
  padding: 1rem;
}

.card {
  background: rgba(16, 24, 48, 0.8);
  backdrop-filter: blur(10px);
  border-radius: 15px;
  padding: 1.5rem;
  transition: all 0.3s ease;
  min-height: 400px;
}

.card:hover {
  transform: translateY(-5px);
  box-shadow: 0 10px 20px rgba(0, 102, 255, 0.3);
}

.chart {
  height: 350px;
  margin-top: 1rem;
}

/* 移动端样式 */
@media (max-width: 768px) {
  .container {
    padding: 1rem;
  }

  .intro h2 {
    font-size: 1.5rem;
  }

  .intro-content p {
    font-size: 0.9rem;
  }

  .grid-container {
    grid-template-columns: 1fr;
  }

  .card {
    min-height: 300px;
  }

  .chart {
    height: 250px;
  }

  .echart1, .echart3 {
    display: none;
  }
}
</style>