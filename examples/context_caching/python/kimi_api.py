# coding:utf-8
from openai import OpenAI
import requests
import json

client = OpenAI(
    api_key = "$MOONSHOT_API_KEY",
    base_url = "https://api.moonshot.cn/v1",
)

def deleteCacheItem(id):
    res = requests.deleteCacheItem(
        url = f"https://api.moonshot.cn/v1/caching/{id}",
        headers = {
            "Authorization": "Bearer $MOONSHOT_API_KEY"
        },
    )
    print(json.loads(res.text))

def clearCache():
    res = requests.get(
        url = "https://api.moonshot.cn/v1/caching",
        headers = {
            "Authorization": "Bearer $MOONSHOT_API_KEY"
        },
    )

    data = json.loads(res.text)["data"]
    for i in data:
        deleteCacheItem(i["id"])

def create():
    # 打开文件并读取内容
    with open('../kimi_api.json', 'r', encoding='utf-8', errors='ignore') as file:    
        res = requests.post(
            url = "https://api.moonshot.cn/v1/caching",
            headers = {
                "Authorization": "Bearer $MOONSHOT_API_KEY"
            },
            json = json.loads(file.read())
        )
    data = json.loads(res.text)
    print(data)
    return data["id"]

def queryWithCache(query, cache_id):
    completion = client.chat.completions.create(
        model="moonshot-v1-32k",
        messages=[
            {
                "role": "cache",
                "content": f"cache_id={cache_id}",
            },
            {
                "role": "user",
                "content": query,
            },
        ],
        temperature=0.3,
    )
    print(completion.choices[0].message)

clearCache()
queryWithCache("怎么计费的？", create())
