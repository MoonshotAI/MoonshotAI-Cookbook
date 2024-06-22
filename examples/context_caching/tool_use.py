from openai import OpenAI
import requests
import json

client = OpenAI(
    api_key = "$MOONSHOT_API_KEY",
    base_url = "https://api.moonshot.cn/v1",
)

def delete(id):
    res = requests.delete(
        url = f"https://api.moonshot.cn/v1/caching/{id}",
        headers = {
            "Authorization": "Bearer $MOONSHOT_API_KEY"
        },
    )
    print(json.loads(res.text))

def clear():
    res = requests.get(
        url = "https://api.moonshot.cn/v1/caching",
        headers = {
            "Authorization": "Bearer $MOONSHOT_API_KEY"
        },
    )

    data = json.loads(res.text)["data"]
    for i in data:
        delete(i["id"])

def create():
    res = requests.post(
        url = "https://api.moonshot.cn/v1/caching",
        headers = {
            "Authorization": "Bearer $MOONSHOT_API_KEY"
        },
        json = {
            "model": "moonshot-v1",
            "messages": [
                {
                    "role": "system",
                    "content": "你是 Kimi，由 Moonshot AI 提供的人工智能助手，你更擅长中文和英文的对话。你会为用户提供安全，有帮助，准确的回答。同时，你会拒绝一切涉及恐怖主义，种族歧视，黄色暴力等问题的回答。Moonshot AI 为专有名词，不可翻译成其他语言。"
                },
            ],
            "tools": [{
                "type": "function",
                "function": {
                    "name": "CodeRunner",
                    "description": "代码执行器，支持运行 python 和 javascript 代码",
                    "parameters": {
                        "properties": {
                            "language": {
                                "type": "string",
                                "enum": ["python", "javascript"]
                            },
                            "code": {
                                "type": "string",
                                "description": "代码写在这里"
                            }
                        },
                        "type": "object"
                    }
                }
            }],
            "name": "CodeRunner",
            "ttl": 3600
        }
    )
    print(json.loads(res.text))

def query_with_cache(query, cache_id):
    completion = client.chat.completions.create(
        model="moonshot-v1-8k",
        messages=[
            {
                "role": "cache",
                "content": f"cache_id={cache_id};dry_run=0",
            },
            {
                "role": "user",
                "content": query,
            },
        ],
        temperature=0.3,
    )
    #print(completion)
    print(completion.choices[0].message)

query_with_cache("编程判断 3214567 是否是素数。", "cache-essqmysd6h1111dauub1")
