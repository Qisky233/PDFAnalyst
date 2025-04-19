// 基础 URL
const BASE_URL = 'http://localhost:8080';

// 模糊搜索作者（by name）
export const searchAuthor = async (name) => {
    try {
        const response = await fetch(`${BASE_URL}/search/authors?name=${name}`);
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            return await response.json();
    } catch (error) {
        console.error('Failed to create author:', error);
        throw error;
    }
};

// 精准搜索作者的所有诗作
export const searchAllPoems = async (authorId, page) => {
    try {
        const response = await fetch(`${BASE_URL}/authors/${authorId}/poems?page=${page}`);
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            return await response.json();
    } catch (error) {
        console.error('Failed to create author:', error);
        throwerror;
    }
}


export const searchPoem = async (q, page) => {
    try {
        const response = await fetch(`${BASE_URL}/search/poems?name=${q}&page=${page}`);
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            return await response.json();
    } catch (error) {
        console.error('Failed to create author:', error);
        throw error;
    }
}
