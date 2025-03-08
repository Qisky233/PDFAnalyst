<!-- App.vue -->
<template>
  <div id="app">
    <!-- 导入导航栏组件 -->
    <NavBar style="height: 60px;width: 100%;" @toggle-drawer="toggleDrawer" >1111</NavBar>

    <!-- 移动端抽屉式导航 -->
    <el-drawer v-model="drawerVisible" direction="rtl" size="200px">
      <ul class="drawer-menu">
        <li><a href="/">首页</a></li>
        <li><a href="/papers">论文全览</a></li>
        <li><a href="/analysis">数据分析</a></li>
        <li><a href="/wordCloud">云词统计</a></li>
        <li><a href="/rcharts">生成R图</a></li>
      </ul>
    </el-drawer>

    <router-view></router-view>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import NavBar from '@/components/navBar/index.vue'; // 导入导航栏组件

const drawerVisible = ref(false)
const toggleDrawer = () => {
  drawerVisible.value = !drawerVisible.value;
}
</script>