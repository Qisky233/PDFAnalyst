// 基础 URL
const BASE_URL = '/api';

// 创建作者
export const createAuthor = async (authorData) => {
  try {
    const response = await fetch(`${BASE_URL}/authors`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(authorData),
    });
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    return await response.json();
  } catch (error) {
    console.error('Failed to create author:', error);
    throw error;
  }
};

// 获取作者列表（分页）
export const getAuthors = async (page = 1) => {
  try {
    const response = await fetch(`${BASE_URL}/authors?page=${page}`);
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    return await response.json();
  } catch (error) {
    console.error('Failed to fetch authors:', error);
    throw error;
  }
};

// 获取作者列表（不分页）
export const getAuthorByNum = async (num = 800) => {
  try {
    const response = await fetch(`${BASE_URL}/authors/num/${num}`);
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    return await response.json();
  } catch (error) {
    console.error('Failed to fetch authors:', error);
    throw error;
  }
};

// 获取单个作者
export const getAuthorById = async (authorId) => {
  try {
    const response = await fetch(`${BASE_URL}/authors/${authorId}`);
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    return await response.json();
  } catch (error) {
    console.error('Failed to fetch author:', error);
    throw error;
  }
};

// 更新作者
export const updateAuthor = async (authorId, authorData) => {
  try {
    const response = await fetch(`${BASE_URL}/authors/${authorId}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(authorData),
    });
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    return await response.json();
  } catch (error) {
    console.error('Failed to update author:', error);
    throw error;
  }
};

// 删除作者
export const deleteAuthor = async (authorId) => {
  try {
    const response = await fetch(`${BASE_URL}/authors/${authorId}`, {
      method: 'DELETE',
    });
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    return await response.json();
  } catch (error) {
    console.error('Failed to delete author:', error);
    throw error;
  }
};