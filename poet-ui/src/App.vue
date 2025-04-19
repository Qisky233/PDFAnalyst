<script setup>
// App level setup
</script>

<template>
  <router-view></router-view>
</template>

<style>
@import url('https://fonts.googleapis.com/css2?family=Noto+Serif+SC:wght@400;500;600;700&display=swap');

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html, body {
  font-family: 'Noto Serif SC', serif;
  background: #000;
  color: #fff;
  overflow-x: hidden;
}

#app {
  width: 100vw;
  min-height: 100vh;
}
</style>
