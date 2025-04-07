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
      <el-select 
        v-model="selectedNumber" 
        placeholder="显示数量"
        class="number-select"
        popper-class="select-popper"
      >
        <el-option
          v-for="num in numberOptions"
          :key="num"
          :label="`前${num}个`"
          :value="num"
        />
      </el-select>
      <el-select 
        v-model="selectedTheme" 
        placeholder="颜色主题"
        class="theme-select"
        popper-class="select-popper"
      >
        <el-option
          v-for="theme in themes"
          :key="theme.value"
          :label="theme.label"
          :value="theme.value"
        />
      </el-select>
      <el-select 
        v-model="selectedShape" 
        placeholder="词云形状"
        class="shape-select"
        popper-class="select-popper"
      >
        <el-option
          v-for="shape in shapes"
          :key="shape.value"
          :label="shape.label"
          :value="shape.value"
        />
      </el-select>
      <button @click="calculate" class="calculate-button">计算</button>
      <button @click="exportWordCloud" class="export-button">导出词云</button>
    </div>
    <div v-show="wordCount.length" class="results-container">
      <!-- 词云图 -->
      <div ref="wordcloudChart" class="chart-container"></div>
      <!-- 高频词列表 -->
      <div class="word-list">
        <h3>高频词语：</h3>
        <div class="list-container">
        <ol>
          <li v-for="([word, count], index) in wordCount" :key="index"
              :class="{ 'highlight': highlightedWord === word }"
              @mouseenter="highlightWord(word)"
              @mouseleave="highlightWord(null)">
            <span class="word">{{ word }}</span>
            <span class="count">{{ count }}</span>
          </li>
        </ol>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { Segment, useDefault } from 'segmentit';
import * as echarts from 'echarts';
import 'echarts-wordcloud';
import { debounce } from 'lodash';
import * as PDFJS from 'pdfjs-dist';
import 'pdfjs-dist/build/pdf.worker';
export default {
  name: 'Calculator',
  data() {
    return {
      inputText: '',
      wordCount: [],
      selectedNumber: 20,
      numberOptions: [5, 10, 15, 20, 30, 40, 50],
      selectedTheme: 'random',
      themes: [
        { label: '随机颜色', value: 'random' },
        { label: '莫兰迪色系', value: 'muted' },
        { label: '蓝色系', value: 'blue' },
        { label: '红色系', value: 'red' },
        { label: '绿色系', value: 'green' },
        { label: '暖色系', value: 'warm' },
        { label: '亮色系', value: 'bright' },
      ],
      selectedShape: 'circle',
      shapes: [
        { label: '圆形', value: 'circle' },
        { label: '星形', value: 'star' },
        { label: '正方形', value: 'square' },
        { label: '三角形', value: 'triangle' },
        { label: '菱形', value: 'diamond' },
        { label: '心形', value: 'cardioid' },
        { label: '五边形', value: 'pentagon' },
      ],
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
    calculate() {
      if (!this.inputText.trim()) {
        this.$message.error('请输入要分析的文本');
        return;
      }

      const segment = useDefault(new Segment());
      const words = segment.doSegment(this.inputText, { simple: true });

      const countMap = words.reduce((acc, word) => {
        const trimmed = word.trim();
        if (trimmed && trimmed.length > 1 && /[\u4e00-\u9fa5]/.test(trimmed)) {
          acc[trimmed] = (acc[trimmed] || 0) + 1;
        }
        return acc;
      }, {});

      this.wordCount = Object.entries(countMap)
        .sort((a, b) => b[1] - a[1])
        .slice(0, this.selectedNumber);

      this.$nextTick(() => {
        this.updateWordCloud();
      });
    },
    updateWordCloud() {
      if (!this.chartInstance) {
        if (this.$refs.wordcloudChart) {
          this.chartInstance = echarts.init(this.$refs.wordcloudChart);
        } else {
          return;
        }
      }
      
      const cloudData = this.wordCount.map(([word, count]) => ({
        name: word,
        value: count
      }));

      const option = {
        title: {
          text: '词云展示',
          left: 'center',
          textStyle: {
            color: '#666'
          }
        },
        tooltip: {
          show: true,
          formatter: function (params) {
            return `${params.name}: ${params.value}`;
          }
        },
        series: [ {
          type: 'wordCloud',
          data: cloudData,
          shape: this.selectedShape, 
          sizeRange: [20, 80],
          rotationRange: [0, 0],
          gridSize: 8,
          textStyle: {
            color: this.getColorFunction(),
            emphasis: {
              shadowBlur: 10,
              shadowColor: '#333'
            }
          }
        }]
      };

      this.chartInstance.clear();
      this.chartInstance.setOption(option);
      this.chartInstance.resize();
      this.chartInstance.on('mouseover', (params) => {
        this.highlightWord(params.name);
      });
      this.chartInstance.on('mouseout', () => {
        this.highlightWord(null);
      });
    }, 
    getColorFunction() {
      const theme = this.selectedTheme;
      switch (theme) {
        case 'random':
          return () => `hsl(${Math.random() * 360}, 70%, 50%)`;
        case 'blue':
          return () => `hsl(${Math.random() * 90 + 180}, 70%, 50%)`; // 蓝色系
        case 'red':
          return () => `hsl(${Math.random() * 30 + 0}, 70%, 50%)`; // 红色系
        case 'green':
          return () => `hsl(${Math.random() * 90 + 90}, 70%, 50%)`;// 绿色系
        case 'warm':
          return () => `hsl(${Math.random() * 60 + 20}, 70%, 50%)`; // 暖色系（黄色到红色）
        case 'bright':
          return () => `hsl(${Math.random() * 360}, 80%, 70%)`; // 亮色系（高亮度、高饱和度）
        case 'muted':
          const mutedColors = ['#a3b7a0', '#d6c7d9', '#e8d9d0', '#c8bfd0', '#b2c3b3', '#b1c2d6', '#c8d5d4', '#dfd6d7', '#c2d2b5', '#9d91a7', '#b8def6', '#dfecda'];
          return () => mutedColors[Math.floor(Math.random() * mutedColors.length)]; // 莫兰迪色系（预定义颜色列表）
        default:
          return () => `hsl(${Math.random() * 360}, 70%, 50%)`;
      }
    },
    handleResize: debounce(function() {
      this.chartInstance?.resize();
    }, 300),
    highlightWord(word) {
      this.highlightedWord = word;
      if (this.chartInstance) {
        this.chartInstance.dispatchAction({
          type: 'downplay', 
        });
        if (word) {
          this.chartInstance.dispatchAction({
            type: 'highlight', 
            name: word
          });
        }
      }
    },
    exportWordCloud() {
      if (!this.chartInstance) {
        this.$message.error('请先生成词云');
        return;
      }
      const url = this.chartInstance.getDataURL({
        type: 'png',
        pixelRatio: 2, // 导出的图片分辨率比例，默认是 1
        backgroundColor: '#fff' 
      });
      const link = document.createElement('a');
      link.href = url;
      link.download = 'wordcloud.png'; 
      link.click();
    }
  }
};
</script>

<style scoped>
@import url(../../css/wordcloud.css);
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

.controls {
  width: 90%;
  display: flex;
  gap: 15px;
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