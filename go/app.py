import json
import sqlite3
import os

# 数据库连接
db_path = 'tang_poetry.db'
conn = sqlite3.connect(db_path)
cursor = conn.cursor()

# 检查表是否存在
cursor.execute("SELECT name FROM sqlite_master WHERE type='table' AND name='Authors'")
if cursor.fetchone() is None:
    # 创建作者表
    cursor.execute('''
    CREATE TABLE Authors (
        author_id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT UNIQUE,
        description TEXT
    )
    ''')

cursor.execute("SELECT name FROM sqlite_master WHERE type='table' AND name='Poems'")
if cursor.fetchone() is None:
    # 创建诗表
    cursor.execute('''
    CREATE TABLE Poems (
        poem_id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT,
        author_id INTEGER,
        content TEXT,
        FOREIGN KEY (author_id) REFERENCES Authors (author_id)
    )
    ''')

cursor.execute("SELECT name FROM sqlite_master WHERE type='table' AND name='Tang300'")
if cursor.fetchone() is None:
    # 创建唐诗三百首表
    cursor.execute('''
    CREATE TABLE Tang300 (
        id TEXT PRIMARY KEY,
        title TEXT,
        author TEXT,
        content TEXT,
        tags TEXT
    )
    ''')

# 提交表创建
conn.commit()

# 导入作者数据
authors_path = '全唐诗/authors.tang.json'
with open(authors_path, 'r', encoding='utf-8') as f:
    authors = json.load(f)
    for author in authors:
        cursor.execute('''
        INSERT OR IGNORE INTO Authors (name, description) VALUES (?, ?)
        ''', (author['name'], author['desc']))

# 提交作者数据
conn.commit()

# 获取作者 ID 映射
author_id_map = {}
cursor.execute('SELECT name, author_id FROM Authors')
for row in cursor.fetchall():
    author_id_map[row[0]] = row[1]

# 导入诗数据
poetry_files = [f for f in os.listdir('全唐诗') if f.startswith('poet.tang.') and f.endswith('.json')]
for file_name in poetry_files:
    file_path = os.path.join('全唐诗', file_name)
    with open(file_path, 'r', encoding='utf-8') as f:
        poems = json.load(f)
        for poem in poems:
            author_name = poem['author']
            author_id = author_id_map.get(author_name)
            if author_id:
                content = '\n'.join(poem['paragraphs'])
                cursor.execute('''
                INSERT INTO Poems (title, author_id, content) VALUES (?, ?, ?)
                ''', (poem['title'], author_id, content))

# 提交诗数据
conn.commit()

# 导入唐诗三百首数据
tang300_path = '全唐诗/唐诗三百首.json'
with open(tang300_path, 'r', encoding='utf-8') as f:
    tang300_poems = json.load(f)
    for poem in tang300_poems:
        content = '\n'.join(poem['paragraphs'])
        tags = ', '.join(poem['tags'])
        cursor.execute('''
        INSERT INTO Tang300 (id, title, author, content, tags) VALUES (?, ?, ?, ?, ?)
        ''', (poem['id'], poem['title'], poem['author'], content, tags))

# 提交唐诗三百首数据
conn.commit()

# 关闭数据库连接
conn.close()

print("数据导入完成！")