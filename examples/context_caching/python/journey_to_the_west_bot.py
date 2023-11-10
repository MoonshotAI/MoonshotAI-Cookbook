# coding:utf-8
from openai import OpenAI
import requests
import json

client = OpenAI(
    api_key = "$MOONSHOT_API_KEY",
    base_url = "https://api.moonshot.cn/v1",
)

def deleteCacheItem(id):
    res = requests.delete(
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
    with open('../西游记.txt', 'r', encoding='utf-8', errors='ignore') as file:
        content = file.read()
        res = requests.post(
            url = "https://api.moonshot.cn/v1/caching",
            headers = {
                "Authorization": "Bearer $MOONSHOT_API_KEY"
            },
            json = {
                "model": "moonshot-v1",
                "messages":
                [
                    {
                        "role": "user",
                        "content": "{\"type\":\"file\",\"document_number\":1,\"file_type\":\"text\",\"file_name\":\"西游记.txt\",\"content\":" + content + "}",
                        "name": "_resource"
                    },
                ],
                "name": "journey_to_west_bot",
                "ttl": 3600
            }
        )
    data = json.loads(res.text)
    print(data)
    return data["id"]

def queryWithCache(query, cache_id):
    completion = client.chat.completions.create(
        model="moonshot-v1-128k",
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
queryWithCache("大闹天宫篇中孙悟空打败了哪些天兵天将？", create())
