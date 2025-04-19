<template>
  <div class="home">
    <div class="hero-section">
      <h1 class="title">品墨千秋</h1>
      <h2 class="subtitle">数据里的大唐诗卷</h2>
      <div @click="explorePoets" class="explore-link">
        点击探索诗人星图
        <div class="arrow"></div>
      </div>
    </div>
    <div class="star-map" ref="starMapContainer"></div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import * as THREE from 'three';
import { OrbitControls } from 'three/examples/jsm/controls/OrbitControls';
import { useRouter } from 'vue-router';

const starMapContainer = ref(null);

// Three.js variables
let scene, camera, renderer, controls;
let starField;

const router = useRouter();

// const explorePoets = (event) => {
//   const x = event.clientX;
//   const y = event.clientY;

//   // 创建一个光晕元素
//   const glow = document.createElement('div');
//   glow.style.position = 'fixed';
//   glow.style.width = '100vw';
//   glow.style.height = '100vh';
//   glow.style.top = '0';
//   glow.style.left = '0';
//   glow.style.background = `radial-gradient(circle at ${x}px ${y}px, rgba(255, 255, 255, 0.8), transparent 70%)`;
//   glow.style.zIndex = '1000';
//   glow.style.pointerEvents = 'none';
//   glow.style.transition = 'opacity 0.5s ease';

//   document.body.appendChild(glow);

//   // 触发路由跳转
//   router.push('/poet/1');

//   // 动画结束后移除光晕元素
//   setTimeout(() => {
//     glow.style.opacity = '0';
//     setTimeout(() => {
//       document.body.removeChild(glow);
//     }, 500);
//   }, 500);
// };

// 星空粒子消散效果
const explorePoets = (event) => {
  const transition = document.startViewTransition(() => {
    router.push('/poet');
  });

  transition.ready.then(() => {
    document.documentElement.animate(
      {
        clipPath: [
          'inset(100% 100% 100% 100%)',  // 开始状态：完全收缩
          'inset(0 0 0 0)'               // 结束状态：完全展开
        ],
        opacity: [0, 1]                  // 透明度也从0到1
      },
      {
        duration: 800,
        easing: 'cubic-bezier(0.4, 0, 0.2, 1)',
        pseudoElement: '::view-transition-new(root)',
      }
    );
  });
}

// 淡入淡出 + 缩放效果
// const explorePoets = (event) => {
//   const transition = document.startViewTransition(() => {
//     router.push('/poet/1');
//   });

//   transition.ready.then(() => {
//     document.documentElement.animate(
//       {
//         opacity: [1, 0],
//         transform: ['scale(1)', 'scale(1.2)']
//       },
//       {
//         duration: 500,
//         easing: 'ease-in',
//         pseudoElement: '::view-transition-new(root)',
//       }
//     );
//   });
// };

// 水墨晕染效果
// const explorePoets = (event) => {
//   const x = event.clientX;
//   const y = event.clientY;
  
//   const transition = document.startViewTransition(() => {
//     router.push('/poet/1');
//   });

//   transition.ready.then(() => {
//     document.documentElement.animate(
//       {
//         maskImage: [
//           'radial-gradient(circle at center, black 0%, transparent 0%)',
//           'radial-gradient(circle at center, black 30%, transparent 70%)'
//         ],
//         maskSize: ['100% 100%', '200% 200%'],
//         maskPosition: [`${x}px ${y}px`, `${x}px ${y}px`]
//       },
//       {
//         duration: 600,
//         easing: 'ease-out',
//         pseudoElement: '::view-transition-new(root)',
//       }
//     );
//   });
// };

// 垂直帘幕效果
// const explorePoets = () => {
//   const transition = document.startViewTransition(() => {
//     router.push('/detail');
//   });

//   transition.ready.then(() => {
//     document.documentElement.animate(
//       {
//         clipPath: [
//           'inset(0 0 0 0)',
//           'inset(100% 0 100% 0)'
//         ]
//       },
//       {
//         duration: 700,
//         easing: 'ease-in-out',
//         pseudoElement: '::view-transition-new(root)',
//       }
//     );
//   });
// };

const init = () => {
  // Scene setup
  scene = new THREE.Scene();
  camera = new THREE.PerspectiveCamera(
    75,
    window.innerWidth / window.innerHeight,
    0.1,
    1000
  );

  renderer = new THREE.WebGLRenderer({ antialias: true, alpha: true });
  renderer.setSize(window.innerWidth, window.innerHeight);
  starMapContainer.value.appendChild(renderer.domElement);

  // Camera position
  camera.position.z = 50;

  // Controls
  controls = new OrbitControls(camera, renderer.domElement);
  controls.enableDamping = true;
  controls.dampingFactor = 0.05;

  // Create star field
  createStarField();

  // Animation loop
  animate();
};

const createStarField = () => {
  const geometry = new THREE.BufferGeometry();
  const vertices = [];

  for (let i = 0; i < 5000; i++) {
    vertices.push(
      THREE.MathUtils.randFloatSpread(2000),
      THREE.MathUtils.randFloatSpread(2000),
      THREE.MathUtils.randFloatSpread(2000)
    );
  }

  geometry.setAttribute(
    'position',
    new THREE.Float32BufferAttribute(vertices, 3)
  );

  const material = new THREE.PointsMaterial({
    color: 0xffffff,
    size: 0.5
  });

  starField = new THREE.Points(geometry, material);
  scene.add(starField);
};

const animate = () => {
  requestAnimationFrame(animate);
  controls.update();

  if (starField) {
    starField.rotation.x += 0.0001;
    starField.rotation.y += 0.0001;
  }

  renderer.render(scene, camera);
};

const handleResize = () => {
  camera.aspect = window.innerWidth / window.innerHeight;
  camera.updateProjectionMatrix();
  renderer.setSize(window.innerWidth, window.innerHeight);
};

onMounted(() => {
  init();
  window.addEventListener('resize', handleResize);
});
</script>

<style scoped>
@font-face {
  font-family: 'xingkai';
  src: url('../assets/fonts/方正行楷_GBK.ttf') format('truetype');
  font-weight: normal;
  font-style: normal;
}

.home {
  width: 100vw;
  height: 100vh;
  overflow: hidden;
  position: relative;
  background: linear-gradient(to bottom, #000000, #1a1a2e);
}

.hero-section {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  text-align: center;
  color: #ffffff;
  z-index: 1;
  font-family: '楷体';
  width: 100%;
}

.title {
  font-family: 'xingkai', '楷体';
  font-size: 5rem;
  margin-bottom: 0.5rem;
  font-weight: bold;
  text-shadow: 0 0 10px rgba(255, 255, 255, 0.5);
}

.subtitle {
  font-size: 1.2rem;
  margin-bottom: 2rem;
  opacity: 0.8;
}

.explore-link {
  display: block;
  color: #fff;
  text-decoration: none;
  font-size: 1rem;
  opacity: 0.8;
  transition: all 0.3s ease;
  cursor: pointer;
  margin-top: 2rem;
}

.explore-link:hover {
  opacity: 1;
  transform: translateY(-2px);
}

.arrow {
  width: 20px;
  height: 20px;
  border-right: 2px solid #fff;
  border-bottom: 2px solid #fff;
  transform: rotate(45deg);
  margin: 10px auto 0;
  animation: bounce 2s infinite;
}

@media (max-width: 768px) {
  .title {
    font-size: 2.5rem;
  }
  .subtitle {
    font-size: 1rem;
  }
}

@keyframes bounce {
  0%, 20%, 50%, 80%, 100% {
    transform: translateY(0) rotate(45deg);
  }
  40% {
    transform: translateY(-10px) rotate(45deg);
  }
  60% {
    transform: translateY(-5px) rotate(45deg);
  }
}
</style>