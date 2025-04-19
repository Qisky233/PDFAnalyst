// 基础 URL
const BASE_URL = '/api';

// 1. 获取数据统计
export const getDataStats = async () => {
  try {
    const response = await fetch(`${BASE_URL}/data/stats`);
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    return await response.json();
  } catch (error) {
    console.error('Failed to fetch data stats:', error);
    throw error;
  }
};

// 2. 获取图表数据
export const getEchartData = async (params) => {
  try {
    if (!params) {
      throw new Error('Params is required');
    }
    
    const response = await fetch(`${BASE_URL}/data/echart/${params}`);
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    return await response.json();
  } catch (error) {
    console.error(`Failed to fetch echart data for params ${params}:`, error);
    throw error;
  }
};

// 3. 获取表格数据
export const getTableData = async () => {
  try {
    const response = await fetch(`${BASE_URL}/data/table`);
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    return await response.json();
  } catch (error) {
    console.error('Failed to fetch table data:', error);
    throw error;
  }
};