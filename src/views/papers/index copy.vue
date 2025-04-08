<template>
  <div style="margin-top: 120px">
    <label for="file-upload" class="upload-button">导入PDF文件</label>
    <input
      id="file-upload"
      type="file"
      accept=".pdf"
      @change="handleFileUpload"
      style="display: none"
    />
    <div v-for="(item, index) in fileList" :key="index" class="pdf-card" @click="showPreview(item)">
    {{ item.name }}
  </div>
    <div v-if="previewVisible" class="preview-modal">
      <div class="preview-content">
        <object 
          v-if="tempDataUrl"
          :data="tempDataUrl"
          type="application/pdf"
          width="100%"
          height="100%"
        >
        </object>
        <button @click="closePreview">关闭预览</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';

const fileList = ref([]);
const previewVisible = ref(false);
const tempDataUrl = ref(null);
onMounted(() => {
  const savedFiles = localStorage.getItem('pdfFiles');
  if (savedFiles) {
    fileList.value = JSON.parse(savedFiles);
  }
});
function handleFileUpload(event) {
  const uploadedFile = event.target.files[0];
  if (uploadedFile?.type === 'application/pdf') {
    const reader = new FileReader();
    reader.onload = (e) => {
      fileList.value.push({
        name: uploadedFile.name,
        dataUrl: e.target.result,
        lastModified: uploadedFile.lastModified
      });
      localStorage.setItem('pdfFiles', JSON.stringify(fileList.value));
    };
    reader.readAsDataURL(uploadedFile);
  }
}
function showPreview(file) {
  tempDataUrl.value = file.dataUrl;
  previewVisible.value = true;
}
function closePreview() {
  previewVisible.value = false;
  tempDataUrl.value = null;
}
</script>

<style>
.upload-button {
  padding: 12px 24px;
  background: #4CAF50;
  color: white;
  border-radius: 6px;
  cursor: pointer;
  transition: background 0.3s;
}

.upload-button:hover {
  background: #45a049;
}

.pdf-card {
  margin-top: 20px;
  padding: 15px;
  border: 2px solid #ddd;
  border-radius: 8px;
  cursor: pointer;
  max-width: 300px;
  word-break: break-all;
  transition: box-shadow 0.3s;
}

.pdf-card:hover {
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.preview-modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0,0,0,0.5);
  display: flex;
  justify-content: center;
  align-items: center;
}

.preview-content {
  background: white;
  width: 80%;
  height: 90vh;
  padding: 20px;
  border-radius: 8px;
  position: relative;
}

.preview-content object {
  border: 1px solid #ddd;
  margin-bottom: 10px;
}

button {
  padding: 8px 16px;
  background: #007bff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}
</style>