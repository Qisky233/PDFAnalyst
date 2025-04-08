<template>
  <div class="container">
    <!-- 工具栏 -->
    <div class="toolbar">
      <el-upload
          class="upload-demo"
          action="#"
          :on-change="handleUpload"
          :auto-upload="false"
          :show-file-list="false"
      >
        <el-button type="primary">上传PDF</el-button>
      </el-upload>
      <div class="upload-history">
        <h3>上传记录</h3>
        <ul>
          <li
              v-for="(file, index) in uploadedFiles"
              :key="index"
              @click="handleSelectFile(file)"
              :class="{ active: file.isActive }"
          >
            {{ file.name }}
          </li>
        </ul>
      </div>
    </div>

    <!-- 内容浏览栏 -->
    <div class="content">
      <iframe
          v-if="currentPdf"
          :src="currentPdf"
          width="100%"
          height="100%"
          style="border: none"
      ></iframe>
      <div v-else class="placeholder">
        <span>请上传或选择PDF文件</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue';
import { ElMessage } from 'element-plus';

// IndexedDB 配置
const dbName = 'uploadDatabase';
const storeName = 'uploadedFiles';
const dbVersion = 1;

// 打开/创建数据库
const openDB = () => {
  return new Promise((resolve, reject) => {
    const request = indexedDB.open(dbName, dbVersion);

    request.onupgradeneeded = (event) => {
      const db = event.target.result;
      if (!db.objectStoreNames.contains(storeName)) {
        db.createObjectStore(storeName, { keyPath: 'name' });
      }
    };

    request.onsuccess = (event) => resolve(event.target.result);
    request.onerror = (event) => reject(event.target.error);
  });
};

// 存储文件到IndexedDB
const addFileToDB = async (file) => {
  return new Promise(async (resolve) => {
    const reader = new FileReader();
    reader.readAsArrayBuffer(file.raw);
    
    reader.onload = async () => {
      const db = await openDB();
      const transaction = db.transaction(storeName, 'readwrite');
      const store = transaction.objectStore(storeName);
      
      await store.put({
        name: file.name,
        data: reader.result, // 存储ArrayBuffer
        isActive: false
      });
      resolve();
    };
  });
};

// 从IndexedDB获取文件
const getAllFilesFromDB = async () => {
  const db = await openDB();
  return new Promise((resolve) => {
    const transaction = db.transaction(storeName, 'readonly');
    const store = transaction.objectStore(storeName);
    const request = store.getAll();

    request.onsuccess = () => resolve(request.result);
  });
};

// 更新文件激活状态
const updateFileInDB = async (fileName, isActive) => {
  const db = await openDB();
  const transaction = db.transaction(storeName, 'readwrite');
  const store = transaction.objectStore(storeName);
  const request = store.get(fileName);

  request.onsuccess = () => {
    const data = request.result;
    data.isActive = isActive;
    store.put(data);
  };
};

// Vue组件逻辑
const uploadedFiles = ref([]);
const currentPdf = ref('');

// 文件上传处理
const handleUpload = async (file) => {
  // 检查重复
  if (uploadedFiles.value.some(f => f.name === file.name)) {
    ElMessage.warning('文件已存在!');
    return;
  }

  try {
    // 存储原始数据到数据库
    await addFileToDB(file);
    
    // 生成新的Blob URL
    const blob = new Blob([file.raw], { type: 'application/pdf' });
    const blobUrl = URL.createObjectURL(blob);
    
    uploadedFiles.value.push({
      name: file.name,
      url: blobUrl,
      isActive: false
    });
    
    // 自动选中新文件
    handleSelectFile(uploadedFiles.value[uploadedFiles.value.length - 1]);
  } catch (error) {
    ElMessage.error('文件上传失败: ' + error.message);
  }
};

// 文件选择逻辑
const handleSelectFile = async (selectedFile) => {
  // 更新激活状态
  uploadedFiles.value.forEach(file => {
    file.isActive = file.name === selectedFile.name;
    if (file.isActive) {
      currentPdf.value = file.url;
      updateFileInDB(file.name, true);
    } else {
      updateFileInDB(file.name, false);
    }
  });
};

// 页面初始化
onMounted(async () => {
  try {
    const savedFiles = await getAllFilesFromDB();
    
    // 重建Blob URL
    uploadedFiles.value = savedFiles.map(file => ({
      ...file,
      url: URL.createObjectURL(new Blob([file.data], { type: 'application/pdf' }))
    }));

    // 恢复选中状态
    const activeFile = uploadedFiles.value.find(f => f.isActive);
    if (activeFile) currentPdf.value = activeFile.url;
  } catch (error) {
    ElMessage.error('加载历史记录失败: ' + error.message);
  }
});

// 清理Blob URL
onBeforeUnmount(() => {
  uploadedFiles.value.forEach(file => URL.revokeObjectURL(file.url));
});
</script>

<style scoped>
/* 样式保持不变 */
.active {
  color: #409eff;
  font-weight: bold;
}

.container {
  display: flex;
  height: calc(100vh - 65px);
  margin-top: 65px;
  background-color: #f5f5f5;
}

.toolbar {
  width: 250px;
  padding: 20px;
  background-color: #fff;
  box-shadow: 2px 0 5px rgba(0, 0, 0, 0.1);
}

.upload-history {
  margin-top: 20px;
}

.upload-history ul {
  list-style: none;
  padding: 0;
}

.upload-history li {
  margin: 10px 0;
  cursor: pointer;
}

.upload-history li:hover {
  color: #409eff;
}

.content {
  flex: 1;
  padding: 20px;
  background-color: #fff;
  display: flex;
  justify-content: center;
  align-items: center;
}

.placeholder {
  text-align: center;
}

@media (max-width: 768px) {
  .container {
    flex-direction: column;
  }

  .toolbar {
    width: 100%;
    height: auto;
    box-shadow: none;
  }

  .content {
    padding: 10px;
  }
}
</style>