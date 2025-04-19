<template>
    <div class="poet-container">
        <navBar />
      <div class="star-map" ref="starMapContainer"></div>
      <div v-if="selectedAuthor" class="author-details" @click.stop>
        <div class="author-header">
          <h2 class="author-name">{{ selectedAuthor.name }}</h2>
          <button class="detail-btn" @click="goToDetail">查看详情</button>
        </div>
        <div class="description">{{ selectedAuthor.description }}</div>
        <div class="stats">
          <div class="stat-item">
            <span class="stat-label">诗作数量(仅收录)</span>
            <span class="stat-value">{{ selectedAuthor.total_poems || 0 }}</span>
          </div>
          <div class="stat-item">
            <button class="detail-btn" @click.stop="toggleSpeech">
              <svg v-if="!isSpeaking" viewBox="0 0 24 24" width="20" height="20">
                <path fill="currentColor" d="M12 3v10.55c-.59-.34-1.27-.55-2-.55-2.21 0-4 1.79-4 4s1.79 4 4 4 4-1.79 4-4V7h4V3h-6zm-1 12.5c-.83 0-1.5.67-1.5 1.5s.67 1.5 1.5 1.5 1.5-.67 1.5-1.5-.67-1.5-1.5-1.5z"/>
              </svg>
              <svg v-else viewBox="0 0 24 24" width="20" height="20">
                <path fill="currentColor" d="M3 9v6h4l5 5V4L7 9H3zm13.5 3c0-1.77-1.02-3.29-2.5-4.03v8.05c1.48-.73 2.5-2.25 2.5-4.02zM14 3.23v2.06c2.89.86 5 3.54 5 6.71s-2.11 5.85-5 6.71v2.06c4.01-.91 7-4.49 7-8.77s-2.99-7.86-7-8.77z"/>
              </svg>
              <span class="speech-label">{{ isSpeaking ? '停止朗读' : '朗读介绍' }}</span>
            </button>
          </div>
        </div>
      </div>
      <div class="debug-console" v-if="showDebug" @click.stop>
        <h3>调试信息</h3>
        <pre>{{ debugInfo }}</pre>
      </div>
      <button class="debug-toggle" @click="showDebug = !showDebug">
        {{ showDebug ? '隐藏' : '显示' }}调试信息
      </button>
    </div>
  </template>
  
<script setup>
import { ref, onMounted, onUnmounted, onBeforeUnmount  } from 'vue';
import * as THREE from 'three';
import { OrbitControls } from 'three/examples/jsm/controls/OrbitControls';
import { FontLoader } from 'three/examples/jsm/loaders/FontLoader';
import { TextGeometry } from 'three/examples/jsm/geometries/TextGeometry';
import { getAuthors, getAuthorByNum } from '@/api/poet.js';
import { searchAuthor } from '@/api/search.js';
import { useRouter } from 'vue-router';
import gsap from 'gsap';

import navBar from '@/components/navBar.vue';

const router = useRouter();
  
// 调试开关
const showDebug = ref(true);
const debugInfo = ref('初始化中...');

const props = defineProps({
  selectedAuthor: Object
});

const isSpeaking = ref(false);
let speechSynthesis = null;
let speechUtterance = null;

// 在组件中添加检查方法
const checkSpeechSupport = () => {
  if (!('speechSynthesis' in window)) {
    console.error('您的浏览器不支持语音合成API');
    return false;
  }
  return true;
};
const toggleSpeech = async () => {
  if (!checkSpeechSupport()) {
    alert('您的浏览器不支持文本朗读功能');
    return;
  }

  if (isSpeaking.value) {
    stopSpeech();
    return;
  }

  if (!selectedAuthor.value?.description) {
    alert('没有可朗读的内容');
    return;
  }

  try {
    // 确保语音API已加载
    if (!speechSynthesis) {
      speechSynthesis = window.speechSynthesis;
      
      // 有些浏览器需要先获取语音列表才能正常工作
      await new Promise(resolve => {
        const voices = speechSynthesis.getVoices();
        if (voices.length > 0) {
          resolve();
        } else {
          speechSynthesis.onvoiceschanged = resolve;
        }
      });
    }

    stopSpeech(); // 停止当前朗读
    
    speechUtterance = new SpeechSynthesisUtterance(selectedAuthor.value.description);
    speechUtterance.lang = 'zh-CN';
    speechUtterance.rate = 0.9;
    speechUtterance.pitch = 1;
    
    // 尝试选择中文语音
    const voices = speechSynthesis.getVoices();
    const naturalVoice = voices.find(voice => 
      voice.name.includes('Google') && voice.lang.includes('zh') // 选择Google提供的中文语音
    );
    if (naturalVoice) {
      speechUtterance.voice = naturalVoice;
    } else {
      console.warn('未找到更自然的语音，使用默认语音');
    }
    
    speechUtterance.onend = () => {
      isSpeaking.value = false;
      console.log('朗读结束');
    };
    
    speechUtterance.onerror = (e) => {
      console.error('朗读出错:', e);
      isSpeaking.value = false;
      alert('朗读出错: ' + e.error);
    };
    
    speechSynthesis.speak(speechUtterance);
    isSpeaking.value = true;
    console.log('开始朗读');
  } catch (error) {
    console.error('朗读异常:', error);
    alert('朗读功能异常: ' + error.message);
  }
};
const startSpeech = () => {
  if (!speechSynthesis) {
    speechSynthesis = window.speechSynthesis;
  }
  
  stopSpeech(); // 停止当前正在进行的朗读
  
  speechUtterance = new SpeechSynthesisUtterance(props.selectedAuthor.description);
  speechUtterance.lang = 'zh-CN';
  speechUtterance.rate = 0.9;
  speechUtterance.pitch = 1;
  
  speechUtterance.onend = () => {
    isSpeaking.value = false;
  };
  
  speechUtterance.onerror = (e) => {
    console.error('朗读出错:', e);
    isSpeaking.value = false;
  };
  
  speechSynthesis.speak(speechUtterance);
  isSpeaking.value = true;
};

const stopSpeech = () => {
  if (speechSynthesis) {
    if (speechUtterance) {
      // 移除事件监听器
      speechUtterance.onend = null;
      speechUtterance.onerror = null;
    }
    speechSynthesis.cancel();
  }
  isSpeaking.value = false;
};

// 组件卸载时停止朗读
onBeforeUnmount(() => {
  stopSpeech();
});
  
// 更新调试信息
const updateDebugInfo = (message) => {
  const timestamp = new Date().toLocaleTimeString();
  debugInfo.value = `${timestamp}: ${message}\n${debugInfo.value}`;
  console.log(message); // 同时输出到控制台
};
  
// DOM引用和状态
const starMapContainer = ref(null);
const selectedAuthor = ref(null);
const authors = ref([]);
  
// 字体URL - 使用更可靠的CDN
const FONT_URL = 'https://cdn.jsdelivr.net/npm/three@0.132.2/examples/fonts/helvetiker_regular.typeface.json';

  
  
  // 获取诗人数据
  const getAuthorsData = async () => {
    try {
      updateDebugInfo('开始获取诗人数据...');
      const response = await getAuthorByNum();
      console.log('获取到的诗人数据:', response?.data);
      updateDebugInfo(`获取到 ${response?.data?.length || 0} 位诗人数据`);
      
      authors.value = response?.data || [];
      
      if (authors.value.length === 0) {
        updateDebugInfo('警告：获取到的诗人数据为空！');
      }
      
      initScene();
    } catch (error) {
      updateDebugInfo(`获取诗人数据失败: ${error.message}`);
      console.error('完整错误:', error);
    }
  };

  // 添加获取诗人详情的方法
    const fetchAuthorDetail = async (authorName) => {
    try {
        updateDebugInfo(`开始获取诗人 ${authorName} 的详细信息...`);
        const response = await searchAuthor(authorName)
        
        if (response?.data?.length > 0) {
        // 找到名字完全匹配的诗人
        const exactMatch = response.data.find(author => author.name === authorName);
        if (exactMatch) {
            updateDebugInfo(`找到诗人 ${authorName} 的详细信息`);
            return exactMatch;
        }
        }
        updateDebugInfo(`未找到诗人 ${authorName} 的详细信息`);
        return null;
    } catch (error) {
        updateDebugInfo(`获取诗人详情失败: ${error.message}`);
        console.error('完整错误:', error);
        return null;
    }
    };
  
  // Three.js变量
  let scene = null;
  let camera = null;
  let renderer = null;
  let controls = null;
  let starField = null;
  let authorNodes = [];
  let textLabels = [];
  let raycaster = new THREE.Raycaster();
  let mouse = new THREE.Vector2();
  let font = null;
  
// 创建星空背景
const createStarField = () => {
  updateDebugInfo('创建星空背景...');
  try {
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
    updateDebugInfo('星空背景创建成功');
  } catch (error) {
    updateDebugInfo(`创建星空背景失败: ${error.message}`);
  }
};
  
// 初始化场景
const initScene = () => {
  updateDebugInfo('开始初始化场景...');
  
  // 检查容器是否存在
  if (!starMapContainer.value) {
    updateDebugInfo('错误：starMapContainer DOM元素未找到！');
    return;
  }

  // 清除旧场景
  if (scene) {
    updateDebugInfo('清理旧场景...');
    while(scene.children.length > 0) { 
      scene.remove(scene.children[0]); 
    }
  }

  // 创建新场景
  try {
    scene = new THREE.Scene();
    updateDebugInfo('场景对象创建成功');
    
    camera = new THREE.PerspectiveCamera(
      45,
      window.innerWidth / window.innerHeight,
      0.1,
      1000
    );
    camera.position.z = 200;
    updateDebugInfo('相机设置完成');

    renderer = new THREE.WebGLRenderer({ 
      antialias: true, 
      alpha: true,
      preserveDrawingBuffer: true
    });
    renderer.setPixelRatio(window.devicePixelRatio);
    renderer.setSize(window.innerWidth, window.innerHeight);
    
    // 清除旧容器内容
    while (starMapContainer.value.firstChild) {
      starMapContainer.value.removeChild(starMapContainer.value.firstChild);
    }
    starMapContainer.value.appendChild(renderer.domElement);
    updateDebugInfo('渲染器设置完成');

    controls = new OrbitControls(camera, renderer.domElement);
    controls.enableDamping = true;
    controls.dampingFactor = 0.05;
    controls.enableZoom = true; // 确保启用缩放功能
    controls.minDistance = 50; // 设置最小缩放距离
    controls.maxDistance = 500; // 设置最大缩放距离
    updateDebugInfo('控制器设置完成');

    // 创建星空背景
    createStarField();
    
    // 加载字体并创建节点
    loadFontAndCreateNodes();
    
    // 开始动画循环
    animate();
    updateDebugInfo('场景初始化完成');
  } catch (error) {
    updateDebugInfo(`场景初始化失败: ${error.message}`);
    console.error('完整错误:', error);
  }
};
  
// 加载字体并创建节点
const loadFontAndCreateNodes = () => {
  updateDebugInfo('开始加载字体...');
  
  if (font) {
    updateDebugInfo('使用已加载的字体');
    createAuthorNodes();
    return;
  }

  const fontLoader = new FontLoader();
  fontLoader.load(
    FONT_URL,
    (loadedFont) => {
      updateDebugInfo('字体加载成功');
      font = loadedFont;
      createAuthorNodes();
    },
    (progress) => {
      updateDebugInfo(`字体加载中: ${(progress.loaded / progress.total * 100).toFixed(1)}%`);
    },
    (error) => {
      updateDebugInfo(`字体加载失败: ${error.message}`);
      console.error('字体加载错误:', error);
      createAuthorNodesWithoutText();
    }
  );
};
  
  // 不带文本的节点创建
  const createAuthorNodesWithoutText = () => {
    updateDebugInfo('创建不带文本的节点...');
    authors.value.forEach((author, index) => {
      createAuthorNode(author, index);
    });
  };
  
// 创建作者节点
const createAuthorNodes = () => {
  updateDebugInfo(`开始创建 ${authors.value.length} 个作者节点...`);
  authorNodes = [];
  textLabels = [];

  if (authors.value.length === 0) {
    updateDebugInfo('警告：没有作者数据可创建节点');
    return;
  }

  try {
    authors.value.forEach((author, index) => {
      createAuthorNode(author, index);
      
      if (font) {
        createTextLabel(author, index);
      }
    });
    updateDebugInfo('作者节点创建完成');
  } catch (error) {
    updateDebugInfo(`创建作者节点失败: ${error.message}`);
  }
};
  
// 创建单个作者节点
// const createAuthorNode = (author, index) => {
//   try {
//     const size = Math.max(0.5, Math.log(author.total_poems + 1) * 1.5);
//     const geometry = new THREE.SphereGeometry(size, 32, 32);
//     const material = new THREE.MeshBasicMaterial({
//       color: getColorByIndex(index),
//       transparent: true,
//       opacity: 0.8
//     });
//     const node = new THREE.Mesh(geometry, material);

//     node.position.set(
//       THREE.MathUtils.randFloatSpread(50),
//       THREE.MathUtils.randFloatSpread(50),
//       THREE.MathUtils.randFloatSpread(50)
//     );

//     node.userData = { author, index };
//     scene.add(node);
//     authorNodes.push(node);

//     // 检查是否有图片URL
//     if (author.imgUrl) {
//       // 创建图片材质
//       const texture = new THREE.TextureLoader().load(author.imgUrl);
//       const spriteMaterial = new THREE.SpriteMaterial({ map: texture });
//       const sprite = new THREE.Sprite(spriteMaterial);
//       sprite.scale.set(10, 10, 1); // 设置图片大小为30x30像素
//       sprite.position.copy(node.position);
//       scene.add(sprite);
//       textLabels.push(sprite);
//     }
//   } catch (error) {
//     updateDebugInfo(`创建作者节点 ${author.name} 失败: ${error.message}`);
//   }
// };

const createAuthorNode = (author, index) => {
  try {
    const size = Math.max(2.5, Math.log(author.total_poems + 2) * 2);

    // 调试信息
    console.log(`创建节点: ${author.name}, 图片URL: ${author.imgUrl || '无'}`);

    if (author.imgUrl) {
      const img = new Image();
      img.crossOrigin = "Anonymous";
      img.src = author.imgUrl;

      img.onload = function () {
        console.log('图片加载成功，创建纹理:', author.imgUrl);

        const canvas = document.createElement('canvas');
        const ctx = canvas.getContext('2d');
        const diameter = 128;
        canvas.width = diameter;
        canvas.height = diameter;

        ctx.beginPath();
        ctx.arc(diameter / 2, diameter / 2, diameter / 2, 0, Math.PI * 2);
        ctx.closePath();
        ctx.clip();
        ctx.drawImage(img, 0, 0, diameter, diameter);

        const texture = new THREE.CanvasTexture(canvas);
        const spriteMaterial = new THREE.SpriteMaterial({
          map: texture,
          transparent: true,
          opacity: 0.6,
        });

        const sprite = new THREE.Sprite(spriteMaterial);
        sprite.scale.set(size * 4, size * 4, 1.5); // 调整节点大小

        // 使用球坐标系随机分布节点位置
        const radius = 100 + Math.random() * 400; // 随机半径
        const phi = Math.acos(1 - 2 * Math.random()); // 随机角度
        const theta = Math.random() * Math.PI * 2; // 随机角度


        // sprite.position.set(
        //   THREE.MathUtils.randFloatSpread(50),
        //   THREE.MathUtils.randFloatSpread(50),
        //   THREE.MathUtils.randFloatSpread(50)
        // );
        sprite.position.set(
          radius * Math.sin(phi) * Math.cos(theta),
          radius * Math.sin(phi) * Math.sin(theta),
          radius * Math.cos(phi)
        );


        // 添加自定义属性以便识别
        sprite.userData = {
          type: 'author',
          author: author,
          index: index,
        };

        // 添加点击事件监听
        sprite.addEventListener('click', () => {
          console.log(`点击了作者节点: ${author.name}`);
          // 这里可以添加你的点击事件逻辑
        });

        scene.add(sprite);
        textLabels.push(sprite);

        // 添加文本标签
        if (font) {
          createTextLabel(author, sprite);
        }

        console.log('精灵创建完成');
      };

      img.onerror = function () {
        console.error('图片加载失败:', author.imgUrl);
      };
    }
  } catch (error) {
    console.error(`创建作者节点 ${author.name} 失败:`, error);
  }
};
// 临时解决方案：使用精灵(Sprite)显示中文
const createTextLabel = (author, index) => {
  const node = authorNodes[index];
  if (!node) return;

  const canvas = document.createElement('canvas');
  canvas.width = 256;
  canvas.height = 128;
  const context = canvas.getContext('2d');
  context.fillStyle = 'rgba(255, 255, 255, 0.9)';
  context.font = '24px Microsoft YaHei';  // 使用系统自带的中文字体
  context.fillText(author.name, 10, 50);

  const texture = new THREE.CanvasTexture(canvas);
  const material = new THREE.SpriteMaterial({ map: texture });
  const sprite = new THREE.Sprite(material);
  
  sprite.position.set(
      node.position.x,
      node.position.y + node.geometry.parameters.radius + 1,
      node.position.z
  );
  sprite.scale.set(10, 5, 1);
  scene.add(sprite);
  textLabels.push(sprite);
};
  
// 获取颜色
const getColorByIndex = (index) => {
  const colors = [
    0xff0000, 0x00ff00, 0x0000ff, 0xffff00, 0xff00ff, 0x00ffff,
    0xff9900, 0x9900ff, 0x0099ff, 0xff0099, 0x99ff00, 0x00ff99
  ];
  return colors[index % colors.length];
};
  
// 动画循环
const animate = () => {
  requestAnimationFrame(animate);
  
  if (controls) controls.update();
  
  // 更新文本标签
  if (textLabels.length > 0) {
    textLabels.forEach(label => {
      if (label.userData.parentNode) {
        const node = label.userData.parentNode;
        const textWidth = label.geometry.boundingBox?.max.x - label.geometry.boundingBox?.min.x || 0;
        
        label.position.set(
          node.position.x - (textWidth / 2),
          node.position.y + node.geometry.parameters.radius + 1,
          node.position.z
        );
        label.lookAt(camera.position);
      }
    });
  }

  if (renderer && scene && camera) {
    renderer.render(scene, camera);
  }
};
  
// 窗口大小调整
const handleResize = () => {
  if (camera && renderer) {
    camera.aspect = window.innerWidth / window.innerHeight;
    camera.updateProjectionMatrix();
    renderer.setSize(window.innerWidth, window.innerHeight);
    updateDebugInfo('窗口大小调整处理完成');
  }
};
  
// 鼠标移动
const handleMouseMove = (event) => {
  mouse.x = (event.clientX / window.innerWidth) * 2 - 1;
  mouse.y = -(event.clientY / window.innerHeight) * 2 + 1;
};

// 点击处理
const handleClick = async () => {
  if (!raycaster || !camera || textLabels.length === 0) return;

  raycaster.setFromCamera(mouse, camera);
  const intersects = raycaster.intersectObjects(textLabels);

  if (intersects.length > 0) {
    const clickedNode = intersects[0].object;
    const author = clickedNode.userData.author;

    // 获取诗人详情
    const detailedAuthor = await fetchAuthorDetail(author.name);

    // 更新selectedAuthor
    selectedAuthor.value = {
      ...author,
      ...(detailedAuthor || {}),
      // 保留原始数据，只更新total_poems
      total_poems: detailedAuthor?.total_poems || author.total_poems || 0,
    };

    updateDebugInfo(`选中作者: ${selectedAuthor.value.name}, 诗作数量: ${selectedAuthor.value.total_poems}`);

    // 高亮节点
    textLabels.forEach((node) => {
      node.material.opacity = node === clickedNode ? 0.9 : 0.6;
    });

    // 将相机位置调整到点击的节点
    // controls.target.copy(clickedNode.position);
    // camera.lookAt(clickedNode.position);
    // 使用 gsap 实现平滑过渡
    gsap.to(camera.position, {
      x: clickedNode.position.x,
      y: clickedNode.position.y,
      z: clickedNode.position.z + 50, // 保持一定的距离
      duration: 1.5, // 动画持续时间
      ease: "power2.inOut" // 缓动效果
    });

    gsap.to(controls.target, {
      x: clickedNode.position.x,
      y: clickedNode.position.y,
      z: clickedNode.position.z,
      duration: 1.5, // 动画持续时间
      ease: "power2.inOut" // 缓动效果
    });
  } else {
    selectedAuthor.value = null;
    textLabels.forEach((node) => {
      node.material.opacity = 0.6;
    });
    updateDebugInfo('取消选中作者');
  }
};

// 添加跳转详情页方法
const goToDetail = () => {
  if (selectedAuthor.value) {
    updateDebugInfo(`跳转到诗人 ${selectedAuthor.value.name} 的详情页`);
    router.push({
      path: '/detail',
      query: {
        authorId: selectedAuthor.value.author_id,
        authorName: selectedAuthor.value.name
      }
    });
  }
};
  
// 生命周期
onMounted(() => {
  updateDebugInfo('组件挂载完成');
  getAuthorsData();
  window.addEventListener('resize', handleResize);
  window.addEventListener('mousemove', handleMouseMove);
  window.addEventListener('click', handleClick);

  // 检查语音支持
  if ('speechSynthesis' in window) {
    updateDebugInfo('浏览器支持语音合成API');
    console.log('可用语音列表:', window.speechSynthesis.getVoices());
  } else {
    updateDebugInfo('浏览器不支持语音合成API');
  }
});
  
onUnmounted(() => {
  updateDebugInfo('组件卸载中...');
  window.removeEventListener('resize', handleResize);
  window.removeEventListener('mousemove', handleMouseMove);
  window.removeEventListener('click', handleClick);
  
  if (renderer) {
    renderer.dispose();
    updateDebugInfo('渲染器已清理');
  }
});
</script>
  
<style scoped>
.poet-container {
  position: relative;
  width: 100%;
  height: 100vh;
  overflow: hidden;
  background-color: #000;
}

.star-map {
  width: 100%;
  height: 100%;
}

.author-details {
  position: absolute;
  top: 60px;
  right: 20px;
  width: 300px;
  background: rgba(0, 0, 0, 0.7);
  color: white;
  padding: 20px;
  border-radius: 10px;
  backdrop-filter: blur(5px);
  z-index: 10;
  max-height: 100vh;
  overflow-y: auto;
}

.author-details h2 {
  margin: 0 0 10px 0;
  color: #ffcc00;
}

.description {
  margin-bottom: 20px;
  line-height: 1.5;
  font-size: 0.9em;
  white-space: pre-line;
}

.stats {
  display: flex;
  margin-bottom: 20px;
}

.stat-item {
  margin-right: 20px;
  text-align: center;
}

.stat-label {
  display: block;
  font-size: 0.8em;
  color: #aaa;
}

.stat-value {
  font-size: 1.5em;
  color: #ffcc00;
}

.debug-console {
  position: absolute;
  bottom: 20px;
  left: 20px;
  width: 400px;
  max-height: 200px;
  background: rgba(0, 0, 0, 0.7);
  color: #0f0;
  padding: 10px;
  border-radius: 5px;
  font-family: monospace;
  font-size: 12px;
  overflow-y: auto;
  z-index: 100;
}

.debug-console pre {
  margin: 0;
  white-space: pre-wrap;
}

.debug-toggle {
  position: absolute;
  bottom: 20px;
  left: 430px;
  background: rgba(0, 0, 0, 0.7);
  color: white;
  border: none;
  padding: 5px 10px;
  border-radius: 5px;
  cursor: pointer;
  z-index: 100;
}
/* 作者头部样式 */
.author-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}
.author-name {
  margin: 0;
  flex: 1;
}
/* 文本按钮样式 */
.detail-btn {
  color: #ffcc00;
  background: transparent;
  border: none;
  padding: 0;
  cursor: pointer;
  font-size: 14px;
  text-decoration: underline;
  transition: all 0.3s;
  white-space: nowrap;
  margin-left: 15px;
}

.detail-btn:hover {
  color: #ffdd33;
  text-decoration: none;
}
</style>