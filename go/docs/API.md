# 唐诗数据库API文档

## 基础信息
- **Base URL**: `http://localhost:8080`
- **Content-Type**: `application/json`

## 作者管理接口

### 1. 创建作者
`POST /authors`
```json
{
  "name": "李白",
  "description": "唐代著名诗人"
}
```

### 2. 获取作者列表（分页）
`GET /authors?page=1`
```json
{
  "page": 1,
  "page_size": 6,
  "total": 100,
  "data": [
    {
      "author_id": 1,
      "name": "李白",
      "description": "唐代著名诗人..."
    }
  ]
}
```

### 3. 获取单个作者
`GET /authors/1`
```json
{
  "author_id": 1,
  "name": "李白",
  "description": "唐代著名诗人..."
}
```

### 4. 更新作者
`PUT /authors/1`
```json
{
  "name": "李白(更新)",
  "description": "更新后的描述"
}
```

### 5. 删除作者
`DELETE /authors/1`

---

## 诗作管理接口

### 1. 创建诗作
`POST /poems`
```json
{
  "title": "静夜思",
  "author_id": 1,
  "content": "床前明月光..."
}
```

### 2. 获取诗作列表（分页）
`GET /poems?page=1`
```json
{
  "page": 1,
  "page_size": 6,
  "total": 50,
  "data": [
    {
      "poem_id": 1,
      "title": "静夜思",
      "author_id": 1,
      "content": "床前明月光..."
    }
  ]
}
```

### 3. 获取单个诗作
`GET /poems/1`
```json
{
  "poem_id": 1,
  "title": "静夜思",
  "author_id": 1,
  "content": "床前明月光..."
}
```

### 4. 更新诗作
`PUT /poems/1`
```json
{
  "title": "静夜思(更新)",
  "author_id": 1,
  "content": "更新后的内容..."
}
```

### 5. 删除诗作
`DELETE /poems/1`

---

## 搜索接口

### 6. 模糊搜索作者
`GET /search/authors`
- **查询参数**：
  - `name`：作者姓名（模糊匹配）
  - `page`：页码，默认为 1
- **响应**：
  - **状态码**：200 OK
  - **响应体**：
  ```json
  {
    "page": 1,
    "page_size": 6,
    "total": 2,
    "data": [
      {
        "author_id": 1,
        "name": "李白",
        "description": "唐代著名诗人..."
      },
      {
        "author_id": 2,
        "name": "李清照",
        "description": "宋代女词人..."
      }
    ]
  }
  ```
  - **状态码**：400 Bad Request
  - **响应体**：
  ```json
  {
    "error": "name 参数不能为空"
  }
  ```
- **示例**：
  - **请求**：
  ```bash
  curl -X GET http://localhost:8080/search/authors?name=李&page=1
  ```
  - **响应**：
  ```json
  {
    "page": 1,
    "page_size": 6,
    "total": 2,
    "data": [
      {
        "author_id": 1,
        "name": "李白",
        "description": "唐代著名诗人..."
      },
      {
        "author_id": 2,
        "name": "李清照",
        "description": "宋代女词人..."
      }
    ]
  }
  ```

### 7. 模糊搜索诗作
`GET /search/poems`
- **查询参数**：
  - `name`：诗的标题或内容（模糊匹配）
  - `page`：页码，默认为 1
- **响应**：
  - **状态码**：200 OK
  - **响应体**：
  ```json
  {
    "page": 1,
    "page_size": 6,
    "total": 5,
    "data": [
      {
        "poem_id": 1,
        "title": "静夜思",
        "author_id": 1,
        "content": "床前明月光，疑是地上霜。举头望明月，低头思故乡。"
      },
      {
        "poem_id": 2,
        "title": "将进酒",
        "author_id": 1,
        "content": "君不见黄河之水天上来，奔流到海不复回。"
      }
    ]
  }
  ```
  - **状态码**：400 Bad Request
  - **响应体**：
  ```json
  {
    "error": "name 参数不能为空"
  }
  ```
- **示例**：
  - **请求**：
  ```bash
  curl -X GET http://localhost:8080/search/poems?name=明月&page=1
  ```
  - **响应**：
  ```json
  {
    "page": 1,
    "page_size": 6,
    "total": 5,
    "data": [
      {
        "poem_id": 1,
        "title": "静夜思",
        "author_id": 1,
        "content": "床前明月光，疑是地上霜。举头望明月，低头思故乡。"
      },
      {
        "poem_id": 2,
        "title": "将进酒",
        "author_id": 1,
        "content": "君不见黄河之水天上来，奔流到海不复回。"
      }
    ]
  }
  ```

### 8. 精准搜索作者的所有诗作
`GET /authors/:id/poems`
- **路径参数**：
  - `id`：作者 ID
- **查询参数**：
  - `page`：页码，默认为 1
- **响应**：
  - **状态码**：200 OK
  - **响应体**：
  ```json
  {
    "page": 1,
    "page_size": 6,
    "total": 10,
    "data": [
      {
        "poem_id": 1,
        "title": "静夜思",
        "author_id": 1,
        "content": "床前明月光，疑是地上霜。举头望明月，低头思故乡。"
      },
      {
        "poem_id": 2,
        "title": "将进酒",
        "author_id": 1,
        "content": "君不见黄河之水天上来，奔流到海不复回。"
      }
    ]
  }
  ```
  - **状态码**：404 Not Found
  - **响应体**：
  ```json
  {
    "error": "Author not found"
  }
  ```
- **示例**：
  - **请求**：
  ```bash
  curl -X GET http://localhost:8080/authors/1/poems?page=1
  ```
  - **响应**：
  ```json
  {
    "page": 1,
    "page_size": 6,
    "total": 10,
    "data": [
      {
        "poem_id": 1,
        "title": "静夜思",
        "author_id": 1,
        "content": "床前明月光，疑是地上霜。举头望明月，低头思故乡。"
      },
      {
        "poem_id": 2,
        "title": "将进酒",
        "author_id": 1,
        "content": "君不见黄河之水天上来，奔流到海不复回。"
      }
    ]
  }
  ```

---

## 响应状态码
| 状态码 | 说明 |
|--------|------|
| 200 | 请求成功 |
| 201 | 创建成功 |
| 400 | 请求参数错误 |
| 404 | 资源未找到 |
| 500 | 服务器内部错误 |

---

以上是完整的 API 文档，包含了作者管理、诗作管理以及搜索接口的详细说明。你可以将这些内容复制到 Markdown 编辑器中使用。