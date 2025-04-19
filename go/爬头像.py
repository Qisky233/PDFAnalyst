import json
import sqlite3
import os
import requests
from bs4 import BeautifulSoup
import logging
from datetime import datetime
import time

# 配置日志
logging.basicConfig(level=logging.DEBUG, format='%(asctime)s - %(levelname)s - %(message)s')

# 数据库连接
db_path = 'tang_poetry.db'
conn = sqlite3.connect(db_path)
cursor = conn.cursor()

# 检查 Authors 表是否存在 imgUrl 字段
cursor.execute('''
    PRAGMA table_info(Authors)
''')
table_info = cursor.fetchall()
if 'imgUrl' not in [col[1] for col in table_info]:
    # 添加 imgUrl 字段
    cursor.execute('''
    ALTER TABLE Authors ADD COLUMN imgUrl TEXT
    ''')
    conn.commit()
    logging.info("Added imgUrl column to Authors table")

# 查询 Authors 表中的所有 name
cursor.execute('SELECT author_id, name FROM Authors')
authors = cursor.fetchall()

# 爬取每个作者的图片URL并写入数据库
for author_id, name in authors:
    url = f"https://www.bing.com/images/search?q={name}"
    try:
        response = requests.get(url, timeout=10)
        response.raise_for_status()

        soup = BeautifulSoup(response.text, 'html.parser')
        img_tag = soup.find('img', class_="mimg")
        if img_tag and 'src' in img_tag.attrs:
            imgUrl = img_tag['src']
            cursor.execute('''
                UPDATE Authors SET imgUrl = ? WHERE author_id = ?
            ''', (imgUrl, author_id))
            conn.commit()
            logging.info(f"Updated {name} with imgUrl: {imgUrl}")
        else:
            logging.warning(f"No image found for {name}")
    except requests.RequestException as e:
        logging.error(f"Request failed for {name}: {e}")
    time.sleep(1)  # 添加延迟，避免请求过快

# 关闭数据库连接
conn.close()
logging.info("Database connection closed")

print("数据导入完成！")