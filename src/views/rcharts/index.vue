<template>
  <div class="calculator">
    <label class="custom-upload-button">
    <input 
      type="file" 
      @change="handleFileUpload" 
      accept=".pdf" 
      class="hidden-input"
    >
    <span class="button-text">📁 导入PDF文件</span>
  </label>
    <textarea
      v-model="inputText"
      placeholder="将文本粘贴到此处..."
      class="textarea"
    ></textarea>
    <div class="controls">
      <button class="calculate-button">生成R图</button>
    </div>
    
  </div>
</template>

<script>
import * as PDFJS from 'pdfjs-dist';
import 'pdfjs-dist/build/pdf.worker';
export default {
  name: 'Calculator',
  data() {
    return {
      inputText: '',
      wordCount: [],
      chartInstance: null,
      highlightedWord: null,
    };
  },
  methods: {
    async handleFileUpload(event) {
      const file = event.target.files[0];
      if (!file) return;

      try { 
        const fileReader = new FileReader();
        fileReader.onload = async (e) => {
        const arrayBuffer = e.target.result;
        const pdf = await PDFJS.getDocument({ data: arrayBuffer }).promise;
        let textContent = '';

        for (let i = 1; i <= pdf.numPages; i++) {
          const page = await pdf.getPage(i);
          const text = await page.getTextContent();
          textContent += text.items.map((item) => item.str).join(' ');
        }

        // 清理文本内容
        textContent = textContent.trim().replace(/\s+/g, ' ');

        console.log("提取的文本内容：", textContent); // 打印提取的文本

        if (textContent) {
          this.inputText = textContent; // 将解析的文本填充到textarea
          this.$message.success('PDF文件解析成功');
        } else {
          this.$message.warning('PDF文件中没有可提取的文本');
        }
      };
        fileReader.readAsArrayBuffer(file);
      } catch (error) {
        console.error(error);
        this.$message.error('PDF文件解析失败');
      }
    },
    
  }
};
</script>

<style scoped>
.calculator {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-top: 80px;
  padding-bottom: 40px;
}

.textarea {
  width: 90%;
  height: 300px;
  margin-bottom: 20px;
  border: 1px solid #e4e7ed;
  border-radius: 4px;
  padding: 12px;
  font-size: 14px;
}


.word-list .highlight {
  background-color: #f3f6f3;
  border-radius: 4px;
}
.hidden-input {
  position: absolute;
  width: 0.1px;
  height: 0.1px;
  opacity: 0;
  overflow: hidden;
  z-index: -1;
}

/* 自定义按钮样式 */
.custom-upload-button {
  display: inline-block;
  padding: 12px 24px;
  background: #409EFF;
  color: white;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  margin-bottom: 20px;
}
.calculate-button {
    padding: 0 35px;
    height: 40px;
    background: #409eff;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    transition: all 0.3s;
  }
  
  .calculate-button:hover {
    background: #66b1ff;
  }
.custom-upload-button:hover {
  background: #66b1ff;
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0,0,0,0.15);
}

.button-text::after {
  /* content: "（支持PDF格式）"; */
  font-size: 0.8em;
  opacity: 0.8;
  margin-left: 8px;
}
</style>