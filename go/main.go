package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type Author struct {
	AuthorID    int    `json:"author_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Poems       []Poem `json:"poems"`       // 作者的部分诗作
	TotalPoems  int    `json:"total_poems"` // 作者的诗作总数
}

type Poem struct {
	PoemID   int    `json:"poem_id"`
	Title    string `json:"title"`
	AuthorID int    `json:"author_id"`
	Content  string `json:"content"`
}

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./tang_poetry.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Printf("Connected to database: %s", "./tang_poetry.db")
	defer db.Close()

	router := gin.Default()

	router.POST("/authors", createAuthor)
	router.GET("/authors", getAuthors)
	router.GET("/authors/:id", getAuthor)
	router.PUT("/authors/:id", updateAuthor)
	router.DELETE("/authors/:id", deleteAuthor)

	router.POST("/poems", createPoem)
	router.GET("/poems", getPoems)
	router.GET("/poems/:id", getPoem)
	router.PUT("/poems/:id", updatePoem)
	router.DELETE("/poems/:id", deletePoem)

	router.GET("/search/authors", searchAuthors)
	router.GET("/search/poems", searchPoems)
	router.GET("/authors/:id/poems", getAuthorAllPoems)

	router.Run(":8080")
}

func createAuthor(c *gin.Context) {
	var author Author
	if err := c.ShouldBindJSON(&author); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stmt, err := db.Prepare("INSERT INTO Authors (name, description) VALUES (?, ?)")
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(author.Name, author.Description)
	if err != nil {
		log.Printf("Error executing statement: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Author created: %v", author)
	c.JSON(http.StatusCreated, gin.H{"message": "Author created"})
}

func getAuthors(c *gin.Context) {
	// 获取请求参数中的页码，默认为第一页
	pageStr := c.Query("page")
	page := 1
	if pageStr != "" {
		page, _ = strconv.Atoi(pageStr)
	}
	if page < 1 {
		page = 1
	}

	// 设置每页显示的条数
	const pageSize = 6
	offset := (page - 1) * pageSize

	// 计算总数
	var totalAuthors int
	err := db.QueryRow("SELECT COUNT(*) FROM Authors").Scan(&totalAuthors)
	if err != nil {
		log.Printf("Error querying total authors: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rows, err := db.Query("SELECT author_id, name, description FROM Authors LIMIT ? OFFSET ?", pageSize, offset)
	if err != nil {
		log.Printf("Error querying database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	authors := []Author{}
	for rows.Next() {
		var author Author
		if err := rows.Scan(&author.AuthorID, &author.Name, &author.Description); err != nil {
			log.Printf("Error scanning row: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		authors = append(authors, author)
	}

	// 返回分页结果和总数
	c.JSON(http.StatusOK, gin.H{
		"page":      page,
		"page_size": pageSize,
		"total":     totalAuthors,
		"data":      authors,
	})
}

func getAuthor(c *gin.Context) {
	id := c.Param("id")
	var author Author
	err := db.QueryRow("SELECT author_id, name, description FROM Authors WHERE author_id = ?", id).Scan(&author.AuthorID, &author.Name, &author.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Author not found: %s", id)
			c.JSON(http.StatusNotFound, gin.H{"error": "Author not found"})
		} else {
			log.Printf("Error querying database: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	log.Printf("Author retrieved: %v", author)
	c.JSON(http.StatusOK, author)
}

func updateAuthor(c *gin.Context) {
	id := c.Param("id")
	var author Author
	if err := c.ShouldBindJSON(&author); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stmt, err := db.Prepare("UPDATE Authors SET name = ?, description = ? WHERE author_id = ?")
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(author.Name, author.Description, id)
	if err != nil {
		log.Printf("Error executing statement: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Author updated: %v", author)
	c.JSON(http.StatusOK, gin.H{"message": "Author updated"})
}

func deleteAuthor(c *gin.Context) {
	id := c.Param("id")
	stmt, err := db.Prepare("DELETE FROM Authors WHERE author_id = ?")
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		log.Printf("Error executing statement: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Author deleted: %s", id)
	c.JSON(http.StatusOK, gin.H{"message": "Author deleted"})
}

func createPoem(c *gin.Context) {
	var poem Poem
	if err := c.ShouldBindJSON(&poem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stmt, err := db.Prepare("INSERT INTO Poems (title, author_id, content) VALUES (?, ?, ?)")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(poem.Title, poem.AuthorID, poem.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Poem created"})
}

func getPoems(c *gin.Context) {
	// 获取请求参数中的页码，默认为第一页
	pageStr := c.Query("page")
	page := 1
	if pageStr != "" {
		page, _ = strconv.Atoi(pageStr)
	}
	if page < 1 {
		page = 1
	}

	// 设置每页显示的条数
	const pageSize = 6
	offset := (page - 1) * pageSize

	// 计算总数
	var totalPoems int
	err := db.QueryRow("SELECT COUNT(*) FROM Poems").Scan(&totalPoems)
	if err != nil {
		log.Printf("Error querying total poems: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rows, err := db.Query("SELECT poem_id, title, author_id, content FROM Poems LIMIT ? OFFSET ?", pageSize, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	poems := []Poem{}
	for rows.Next() {
		var poem Poem
		if err := rows.Scan(&poem.PoemID, &poem.Title, &poem.AuthorID, &poem.Content); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		poems = append(poems, poem)
	}

	// 返回分页结果和总数
	c.JSON(http.StatusOK, gin.H{
		"page":      page,
		"page_size": pageSize,
		"total":     totalPoems,
		"data":      poems,
	})
}

func getPoem(c *gin.Context) {
	id := c.Param("id")
	var poem Poem
	err := db.QueryRow("SELECT poem_id, title, author_id, content FROM Poems WHERE poem_id = ?", id).Scan(&poem.PoemID, &poem.Title, &poem.AuthorID, &poem.Content)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Poem not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, poem)
}

func updatePoem(c *gin.Context) {
	id := c.Param("id")
	var poem Poem
	if err := c.ShouldBindJSON(&poem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stmt, err := db.Prepare("UPDATE Poems SET title = ?, author_id = ?, content = ? WHERE poem_id = ?")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(poem.Title, poem.AuthorID, poem.Content, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Poem updated"})
}

func deletePoem(c *gin.Context) {
	// 获取 URL 参数中的 poem_id
	id := c.Param("id")

	// 准备 SQL 语句
	stmt, err := db.Prepare("DELETE FROM Poems WHERE poem_id = ?")
	if err != nil {
		// 如果准备语句失败，返回 500 Internal Server Error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法准备 SQL 语句"})
		return
	}
	defer stmt.Close()

	// 执行删除操作
	result, err := stmt.Exec(id)
	if err != nil {
		// 如果执行失败，返回 500 Internal Server Error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除诗时发生错误"})
		return
	}

	// 获取删除的行数
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		// 如果获取受影响行数失败，返回 500 Internal Server Error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法获取受影响的行数"})
		return
	}

	if rowsAffected == 0 {
		// 如果没有行被删除，返回 404 Not Found
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到诗"})
		return
	}

	// 如果删除成功，返回 200 OK
	c.JSON(http.StatusOK, gin.H{"message": "Poem deleted"})
}

// 搜索作者（模糊匹配）
func searchAuthors(c *gin.Context) {
	// 获取查询参数中的 name
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name 参数不能为空"})
		return
	}

	// 获取请求参数中的页码，默认为第一页
	pageStr := c.Query("page")
	page := 1
	if pageStr != "" {
		page, _ = strconv.Atoi(pageStr)
	}
	if page < 1 {
		page = 1
	}

	// 设置每页显示的条数
	const pageSize = 6
	offset := (page - 1) * pageSize

	// 查询作者总数
	var totalAuthors int
	err := db.QueryRow("SELECT COUNT(*) FROM Authors WHERE name LIKE ?", "%"+name+"%").Scan(&totalAuthors)
	if err != nil {
		log.Printf("Error querying total authors: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 查询作者
	rows, err := db.Query("SELECT author_id, name, description FROM Authors WHERE name LIKE ? LIMIT ? OFFSET ?", "%"+name+"%", pageSize, offset)
	if err != nil {
		log.Printf("Error querying database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	authors := []Author{}
	for rows.Next() {
		var author Author
		if err := rows.Scan(&author.AuthorID, &author.Name, &author.Description); err != nil {
			log.Printf("Error scanning row: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// 查询该作者的部分诗作（分页显示）
		poems, totalPoems, err := getAuthorPoems(author.AuthorID, 1, 6) // 默认为第一页，每页6条
		if err != nil {
			log.Printf("Error querying poems for author %d: %v", author.AuthorID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		author.Poems = poems
		author.TotalPoems = totalPoems
		authors = append(authors, author)
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"page":      page,
		"page_size": pageSize,
		"total":     totalAuthors,
		"data":      authors,
	})
}

// 辅助函数：获取指定作者的部分诗作（分页显示）
func getAuthorPoems(authorID int, page int, pageSize int) ([]Poem, int, error) {
	offset := (page - 1) * pageSize

	// 查询该作者的诗作总数
	var totalPoems int
	err := db.QueryRow("SELECT COUNT(*) FROM Poems WHERE author_id = ?", authorID).Scan(&totalPoems)
	if err != nil {
		return nil, 0, err
	}

	// 查询该作者的部分诗作
	rows, err := db.Query("SELECT poem_id, title, content FROM Poems WHERE author_id = ? LIMIT ? OFFSET ?", authorID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	poems := []Poem{}
	for rows.Next() {
		var poem Poem
		if err := rows.Scan(&poem.PoemID, &poem.Title, &poem.Content); err != nil {
			return nil, 0, err
		}
		poems = append(poems, poem)
	}

	return poems, totalPoems, nil
}

// 搜索诗作（模糊匹配）
func searchPoems(c *gin.Context) {
	// 获取查询参数中的 name
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name 参数不能为空"})
		return
	}

	// 获取请求参数中的页码，默认为第一页
	pageStr := c.Query("page")
	page := 1
	if pageStr != "" {
		page, _ = strconv.Atoi(pageStr)
	}
	if page < 1 {
		page = 1
	}

	// 设置每页显示的条数
	const pageSize = 6
	offset := (page - 1) * pageSize

	// 查询诗作总数
	var totalPoems int
	err := db.QueryRow("SELECT COUNT(*) FROM Poems WHERE title LIKE ? OR content LIKE ?", "%"+name+"%", "%"+name+"%").Scan(&totalPoems)
	if err != nil {
		log.Printf("Error querying total poems: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 查询诗作
	rows, err := db.Query("SELECT poem_id, title, author_id, content FROM Poems WHERE title LIKE ? OR content LIKE ? LIMIT ? OFFSET ?", "%"+name+"%", "%"+name+"%", pageSize, offset)
	if err != nil {
		log.Printf("Error querying database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	poems := []Poem{}
	for rows.Next() {
		var poem Poem
		if err := rows.Scan(&poem.PoemID, &poem.Title, &poem.AuthorID, &poem.Content); err != nil {
			log.Printf("Error scanning row: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		poems = append(poems, poem)
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"page":      page,
		"page_size": pageSize,
		"total":     totalPoems,
		"data":      poems,
	})
}

func getAuthorAllPoems(c *gin.Context) {
	// 获取 URL 参数中的 author_id
	authorIDStr := c.Param("id")
	authorID, err := strconv.Atoi(authorIDStr)
	if err != nil {
		log.Printf("Error parsing author_id: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid author_id"})
		return
	}

	// 获取请求参数中的页码，默认为第一页
	pageStr := c.Query("page")
	page := 1
	if pageStr != "" {
		page, _ = strconv.Atoi(pageStr)
	}
	if page < 1 {
		page = 1
	}

	// 设置每页显示的条数
	const pageSize = 6
	offset := (page - 1) * pageSize

	// 查询该作者的诗作总数
	var totalPoems int
	err = db.QueryRow("SELECT COUNT(*) FROM Poems WHERE author_id = ?", authorID).Scan(&totalPoems)
	if err != nil {
		log.Printf("Error querying total poems for author %d: %v", authorID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 查询该作者的诗作
	rows, err := db.Query("SELECT poem_id, title, content FROM Poems WHERE author_id = ? LIMIT ? OFFSET ?", authorID, pageSize, offset)
	if err != nil {
		log.Printf("Error querying poems for author %d: %v", authorID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	poems := []Poem{}
	for rows.Next() {
		var poem Poem
		if err := rows.Scan(&poem.PoemID, &poem.Title, &poem.Content); err != nil {
			log.Printf("Error scanning row: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		poems = append(poems, poem)
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"page":      page,
		"page_size": pageSize,
		"total":     totalPoems,
		"data":      poems,
	})
}
