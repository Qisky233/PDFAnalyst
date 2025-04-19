<template>
  <div class="image-review-container">
    <div class="header">
      <h1>作者图片批量审核系统</h1>
      <div class="controls">
        <button class="load-btn" @click="loadAuthors">
          <i class="icon-refresh"></i> 加载100位作者
        </button>
        <button class="save-btn" @click="saveChanges" :disabled="!hasChanges">
          <i class="icon-save"></i> 保存所有更改
        </button>
        <button class="remove-btn" @click="batchRemoveImages">
          <i class="icon-delete"></i> 批量删除选中
        </button>
      </div>
      <div class="stats">
        <div class="stat-item">
          <span class="stat-label">总作者数:</span>
          <span class="stat-value">{{ authors.length }}</span>
        </div>
        <div class="stat-item">
          <span class="stat-label">待审核:</span>
          <span class="stat-value">{{ pendingReviewCount }}</span>
        </div>
        <div class="stat-item">
          <span class="stat-label">已修改:</span>
          <span class="stat-value">{{ modifiedCount }}</span>
        </div>
        <div class="stat-item">
          <span class="stat-label">已选中:</span>
          <span class="stat-value">{{ selectedAuthors.length }}</span>
        </div>
      </div>
    </div>

    <!-- 图片网格 -->
    <div class="image-grid-container">
      <div class="image-grid">
        <div 
          v-for="(author, index) in authors" 
          :key="author.author_id"
          class="image-cell"
          :class="{ 
            'no-image': !author.imgUrl,
            'selected': selectedAuthors.includes(author.author_id),
            'modified': author.imgUrl !== author.originalImgUrl,
            'new-row': index % 10 === 0
          }"
          @click="toggleSelect(author.author_id)"
          @dblclick="showDetail(author)"
        >
          <div class="image-wrapper">
            <img 
              v-if="author.imgUrl" 
              :src="author.imgUrl" 
              :alt="author.name"
              @error="handleImageError(author)"
            />
            <div v-else class="no-image-placeholder">
              <i class="icon-no-image"></i>
            </div>
            <div v-if="selectedAuthors.includes(author.author_id)" class="selection-check">
              <i class="icon-check"></i>
            </div>
          </div>
          <div class="author-name">{{ author.name }}</div>
          <div v-if="author.imgUrl !== author.originalImgUrl" class="modified-badge">
            <i class="icon-modified"></i>
          </div>
        </div>
      </div>
    </div>

    <!-- 批量操作工具栏 -->
    <div v-if="selectedAuthors.length > 0" class="batch-toolbar">
      <div class="batch-toolbar-content">
        <span class="selected-count">已选中 {{ selectedAuthors.length }} 位作者</span>
        <button class="toolbar-btn remove-btn" @click="batchRemoveImages">
          <i class="icon-delete"></i> 批量删除图片
        </button>
        <button class="toolbar-btn cancel-btn" @click="clearSelection">
          <i class="icon-cancel"></i> 取消选择
        </button>
      </div>
    </div>

    <!-- 详情弹窗 -->
    <div v-if="selectedAuthor" class="modal-overlay" @click.self="closeDetail">
      <div class="modal-content">
        <div class="modal-header">
          <h2>{{ selectedAuthor.name }}</h2>
          <button class="close-btn" @click="closeDetail">
            <i class="icon-close"></i>
          </button>
        </div>
        
        <div class="modal-body">
          <div class="image-preview-container">
            <div class="image-preview">
              <img 
                v-if="selectedAuthor.imgUrl" 
                :src="selectedAuthor.imgUrl" 
                :alt="selectedAuthor.name"
                @error="handleImageError(selectedAuthor)"
              />
              <div v-else class="no-image-large">
                <i class="icon-no-image-large"></i>
                <span>无图片</span>
              </div>
            </div>
            <button class="test-btn" @click="testImage">
              <i class="icon-test"></i> 测试图片
            </button>
          </div>
          
          <div class="author-info">
            <div class="form-group">
              <label>图片URL:</label>
              <div class="input-with-button">
                <input v-model="editForm.imgUrl" placeholder="输入新的图片URL" />
              </div>
            </div>
            
            <div class="form-group">
              <label>姓名:</label>
              <input v-model="editForm.name" />
            </div>
            
            <div class="form-group">
              <label>描述:</label>
              <textarea v-model="editForm.description" rows="5"></textarea>
            </div>
            
            <div class="action-buttons">
              <button class="btn remove-btn" @click="removeImage">
                <i class="icon-delete"></i> 删除图片
              </button>
              <button class="btn save-btn" @click="saveAuthor">
                <i class="icon-save"></i> 保存更改
              </button>
              <button class="btn cancel-btn" @click="closeDetail">
                <i class="icon-cancel"></i> 取消
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue';
import axios from 'axios';

const BASE_URL = 'http://localhost:8080';

// 数据状态
const authors = ref([]);
const selectedAuthor = ref(null);
const selectedAuthors = ref([]);
const editForm = reactive({
  name: '',
  description: '',
  imgUrl: ''
});

// 计算属性
const hasChanges = computed(() => authors.value.some(a => a.imgUrl !== a.originalImgUrl));
const modifiedCount = computed(() => authors.value.filter(a => a.imgUrl !== a.originalImgUrl).length);
const pendingReviewCount = computed(() => authors.value.filter(a => !a.imgUrl).length);

// 加载100位作者数据
const loadAuthors = async () => {
  try {
    const response = await axios.get(`${BASE_URL}/authors/num/800`);
    authors.value = response.data.data.map(author => ({
      ...author,
      originalImgUrl: author.imgUrl
    }));
    selectedAuthors.value = [];
  } catch (error) {
    console.error('加载作者数据失败:', error);
    alert('加载数据失败，请检查API服务');
  }
};

// 切换选择状态
const toggleSelect = (authorId) => {
  const index = selectedAuthors.value.indexOf(authorId);
  if (index === -1) {
    selectedAuthors.value.push(authorId);
  } else {
    selectedAuthors.value.splice(index, 1);
  }
};

// 清空选择
const clearSelection = () => {
  selectedAuthors.value = [];
};

// 显示作者详情
const showDetail = (author) => {
  selectedAuthor.value = author;
  editForm.name = author.name;
  editForm.description = author.description;
  editForm.imgUrl = author.imgUrl;
};

// 关闭详情弹窗
const closeDetail = () => {
  selectedAuthor.value = null;
};

// 处理图片加载错误
const handleImageError = (author) => {
  author.imgUrl = '';
};

// 测试图片URL
const testImage = () => {
  if (!editForm.imgUrl) return;
  
  const img = new Image();
  img.onload = () => alert('图片加载成功！');
  img.onerror = () => alert('图片加载失败，请检查URL');
  img.src = editForm.imgUrl;
};

// 保存单个作者修改
const saveAuthor = async () => {
  if (!selectedAuthor.value) return;

  try {
    const response = await axios.put(
      `${BASE_URL}/authors/${selectedAuthor.value.author_id}`,
      {
        name: editForm.name,
        description: editForm.description,
        imgUrl: editForm.imgUrl
      }
    );

    const updatedAuthor = response.data;
    const index = authors.value.findIndex(a => a.author_id === updatedAuthor.author_id);
    if (index !== -1) {
      authors.value[index] = {
        ...updatedAuthor,
        originalImgUrl: updatedAuthor.imgUrl
      };
    }

    closeDetail();
    alert('修改已保存！');
  } catch (error) {
    console.error('保存失败:', error);
    alert('保存失败，请检查网络或API服务');
  }
};

// 删除单个作者图片
const removeImage = async () => {
  if (!selectedAuthor.value) return;

  try {
    const response = await axios.put(
      `${BASE_URL}/authors/${selectedAuthor.value.author_id}`,
      {
        ...selectedAuthor.value,
        imgUrl: ''
      }
    );

    const updatedAuthor = response.data;
    const index = authors.value.findIndex(a => a.author_id === updatedAuthor.author_id);
    if (index !== -1) {
      authors.value[index] = {
        ...updatedAuthor,
        originalImgUrl: updatedAuthor.imgUrl
      };
    }

    closeDetail();
    alert('图片已删除！');
  } catch (error) {
    console.error('删除失败:', error);
    alert('删除失败，请检查网络或API服务');
  }
};

// 批量删除图片
const batchRemoveImages = async () => {
  if (selectedAuthors.value.length === 0) return;

  if (!confirm(`确定要批量删除 ${selectedAuthors.value.length} 位作者的图片吗？`)) {
    return;
  }

  try {
    const requests = selectedAuthors.value.map(authorId => {
      const author = authors.value.find(a => a.author_id === authorId);
      return axios.put(`${BASE_URL}/authors/${authorId}`, {
        ...author,
        imgUrl: ''
      });
    });

    const responses = await Promise.all(requests);

    responses.forEach(response => {
      const updatedAuthor = response.data;
      const index = authors.value.findIndex(a => a.author_id === updatedAuthor.author_id);
      if (index !== -1) {
        authors.value[index] = {
          ...updatedAuthor,
          originalImgUrl: updatedAuthor.imgUrl
        };
      }
    });

    selectedAuthors.value = [];
    alert(`成功删除 ${responses.length} 位作者的图片！`);
  } catch (error) {
    console.error('批量删除失败:', error);
    alert('部分删除失败，请检查网络或API服务');
  }
};

// 保存所有更改
const saveChanges = async () => {
  if (!hasChanges.value) return;

  const changedAuthors = authors.value.filter(a => a.imgUrl !== a.originalImgUrl);
  
  if (!confirm(`确定要保存 ${changedAuthors.length} 处修改吗？`)) {
    return;
  }

  try {
    const requests = changedAuthors.map(author => 
      axios.put(`${BASE_URL}/authors/${author.author_id}`, {
        name: author.name,
        description: author.description,
        imgUrl: author.imgUrl
      })
    );

    const responses = await Promise.all(requests);

    responses.forEach(response => {
      const updatedAuthor = response.data;
      const index = authors.value.findIndex(a => a.author_id === updatedAuthor.author_id);
      if (index !== -1) {
        authors.value[index].originalImgUrl = updatedAuthor.imgUrl;
      }
    });

    alert(`成功保存 ${responses.length} 处修改！`);
  } catch (error) {
    console.error('保存失败:', error);
    alert('部分修改保存失败，请检查网络或API服务');
  }
};
</script>

<style scoped>
/* 基础样式 */
:root {
  --primary-color: #3498db;
  --secondary-color: #2ecc71;
  --danger-color: #e74c3c;
  --warning-color: #f39c12;
  --light-color: #ecf0f1;
  --dark-color: #34495e;
  --gray-color: #95a5a6;
  --border-radius: 8px;
  --box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  --transition: all 0.3s ease;
}

* {
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}

body {
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  line-height: 1.6;
  color: #333;
  /* background-color: #f5f7fa; */
}

/* 图标样式 */
.icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

/* 容器样式 */
.image-review-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

/* 头部样式 */
.header {
  /* background-color: white; */
  padding: 20px;
  border-radius: var(--border-radius);
  box-shadow: var(--box-shadow);
  margin-bottom: 20px;
}

.header h1 {
  color: var(--dark-color);
  margin-bottom: 15px;
  font-size: 24px;
  font-weight: 600;
}

.controls {
  display: flex;
  gap: 10px;
  margin-bottom: 15px;
}

button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 5px;
  padding: 8px 16px;
  border: none;
  border-radius: var(--border-radius);
  cursor: pointer;
  font-weight: 500;
  transition: var(--transition);
}

button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.load-btn {
  background-color: var(--primary-color);
  color: white;
}

.save-btn {
  background-color: var(--secondary-color);
  color: white;
}

.remove-btn {
  background-color: var(--danger-color);
  color: white;
}

.cancel-btn {
  background-color: var(--gray-color);
  color: white;
}

.test-btn {
  background-color: var(--warning-color);
  color: white;
  margin-top: 10px;
}

.stats {
  display: flex;
  gap: 20px;
  flex-wrap: wrap;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 5px;
}

.stat-label {
  color: var(--gray-color);
}

.stat-value {
  font-weight: 600;
  color: var(--dark-color);
}

/* 图片网格容器 */
.image-grid-container {
  /* background-color: white; */
  padding: 20px;
  border-radius: var(--border-radius);
  box-shadow: var(--box-shadow);
}

.image-grid {
  display: grid;
  grid-template-columns: repeat(10, 1fr);
  gap: 15px;
}

.image-cell {
  display: flex;
  flex-direction: column;
  align-items: center;
  cursor: pointer;
  transition: var(--transition);
  position: relative;
  padding: 5px;
  border-radius: var(--border-radius);
}

.image-cell:hover {
  background-color: var(--light-color);
}

.image-cell.selected {
  background-color: rgba(52, 152, 219, 0.1);
  box-shadow: 0 0 0 2px var(--primary-color);
}

.image-cell.no-image {
  opacity: 0.7;
}

.image-wrapper {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  overflow: hidden;
  background-color: #f0f0f0;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  border: 2px solid #ddd;
}

.image-wrapper img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.no-image-placeholder {
  color: var(--gray-color);
  font-size: 24px;
}

.selection-check {
  position: absolute;
  top: 0;
  right: 0;
  background-color: var(--primary-color);
  color: white;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
}

.author-name {
  margin-top: 8px;
  font-size: 12px;
  text-align: center;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  width: 100%;
  color: var(--dark-color);
}

.modified-badge {
  position: absolute;
  top: 0;
  right: 0;
  background-color: var(--secondary-color);
  color: white;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
}

/* 批量操作工具栏 */
.batch-toolbar {
  position: fixed;
  bottom: 20px;
  left: 50%;
  transform: translateX(-50%);
  /* background-color: white; */
  padding: 10px 20px;
  border-radius: 30px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.2);
  z-index: 100;
  min-width: 300px;
}

.batch-toolbar-content {
  display: flex;
  align-items: center;
  gap: 15px;
}

.selected-count {
  font-weight: 600;
  color: var(--dark-color);
}

.toolbar-btn {
  padding: 6px 12px;
  font-size: 13px;
}

/* 模态框样式 */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
  backdrop-filter: blur(5px);
}

.modal-content {
  background-color: white;
  border-radius: var(--border-radius);
  width: 90%;
  max-width: 800px;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.2);
  animation: modalFadeIn 0.3s ease;
}

@keyframes modalFadeIn {
  from {
    opacity: 0;
    transform: translateY(-20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #eee;
  position: sticky;
  top: 0;
  background-color: white;
  z-index: 10;
}

.modal-header h2 {
  color: var(--dark-color);
  font-size: 20px;
}

.close-btn {
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: var(--gray-color);
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: var(--transition);
}

.close-btn:hover {
  background-color: var(--light-color);
}

.modal-body {
  display: flex;
  padding: 20px;
  gap: 30px;
}

.image-preview-container {
  flex: 0 0 300px;
}

.image-preview {
  width: 100%;
  height: 300px;
  border-radius: var(--border-radius);
  background-color: #f0f0f0;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  border: 1px solid #ddd;
}

.image-preview img {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
}

.no-image-large {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: var(--gray-color);
  font-size: 18px;
  gap: 10px;
}

.no-image-large i {
  font-size: 48px;
}

.author-info {
  flex: 1;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-weight: 600;
  color: var(--dark-color);
}

.form-group input,
.form-group textarea {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: var(--border-radius);
  font-family: inherit;
  transition: var(--transition);
}

.form-group input:focus,
.form-group textarea:focus {
  outline: none;
  border-color: var(--primary-color);
  box-shadow: 0 0 0 2px rgba(52, 152, 219, 0.2);
}

.form-group textarea {
  min-height: 120px;
  resize: vertical;
}

.action-buttons {
  display: flex;
  gap: 10px;
  margin-top: 30px;
}

.action-buttons .btn {
  flex: 1;
  padding: 10px;
}

/* 响应式调整 */
@media (max-width: 1024px) {
  .image-grid {
    grid-template-columns: repeat(5, 1fr);
  }
}

@media (max-width: 768px) {
  .modal-body {
    flex-direction: column;
  }
  
  .image-preview-container {
    flex: auto;
    margin-bottom: 20px;
  }
  
  .controls {
    flex-wrap: wrap;
  }
}

@media (max-width: 480px) {
  .image-grid {
    grid-template-columns: repeat(3, 1fr);
  }
  
  .stats {
    flex-direction: column;
    gap: 5px;
  }
  
  .batch-toolbar {
    width: 90%;
    border-radius: var(--border-radius);
  }
  
  .batch-toolbar-content {
    flex-direction: column;
    gap: 10px;
  }
  
  .toolbar-btn {
    width: 100%;
  }
}
</style>