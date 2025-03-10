import * as echarts from 'echarts'
import { debounce } from 'lodash'
import 'echarts-wordcloud'

const annualData = {
    xAxis: ['2018', '2019', '2020', '2021', '2022', '2023'],
    series: [1232, 1433, 1765, 1982, 2176, 2354]
}


const keywordsData = [
    { name: '人工智能', value: 100 },
    { name: '机器学习', value: 95 },
    { name: '深度学习', value: 90 },
    { name: '神经网络', value: 85 },
    // ...更多数据
]

const institutionsData = [
    { value: 335, name: '大黄' },
    { value: 310, name: '枸杞' },
    { value: 234, name: '麦冬' },
    // ...更多数据
]

let chartInstance1 = null
let chartInstance2 = null
let chartInstance3 = null

export const initCharts = (chart1, chart2, chart3) => {
    // 年度趋势（折线图）
    chartInstance1 = echarts.init(chart1.value)
    chartInstance1.setOption({
        xAxis: {
            type: 'category',
            data: annualData.xAxis,
            axisLabel: { color: '#fff' }
        },
        yAxis: { type: 'value', axisLabel: { color: '#fff' } },
        series: [{
            data: annualData.series,
            type: 'line',
            smooth: true,
            itemStyle: { color: '#00f7ff' },
            lineStyle: { color: '#0066ff', width: 2 }
        }],
        grid: { top: '20%', bottom: '15%' }
    })

    // 词云
    chartInstance2 = echarts.init(chart2.value)
    chartInstance2.setOption({
        series: [{
            type: 'wordCloud',
            data: keywordsData,
            textStyle: {
                color: () => `hsl(${Math.random() * 360}, 100%, 50%)`
            },
            emphasis: {
                focus: true
            }
        }]
    })

    // 机构分布（饼图）
    chartInstance3 = echarts.init(chart3.value)
    chartInstance3.setOption({
        tooltip: { trigger: 'item' },
        series: [{
            type: 'pie',
            radius: '60%',
            data: institutionsData,
            itemStyle: {
                color: (params) => echarts.color.modifyHSL('#0066ff', params.dataIndex * 40)
            },
            label: { color: '#fff' }
        }]
    })
}

export const destroyCharts = (chart1, chart2, chart3) => {
    chartInstance1?.dispose()
    chartInstance2?.dispose()
    chartInstance3?.dispose()
}

export const handleResize = debounce(() => {
    chartInstance1?.resize()
    chartInstance2?.resize()
    chartInstance3?.resize()
}, 300)