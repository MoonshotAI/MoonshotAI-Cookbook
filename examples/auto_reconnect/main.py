from openai import OpenAI
import time
 
client = OpenAI(
    api_key = "$MOONSHOT_API_KEY",
    base_url = "https://api.moonshot.cn/v1",
)

def chat_once(msgs):
    response = client.chat.completions.create(
        model = "moonshot-v1-auto",
        messages = msgs,
        temperature = 0.3,
    )
    return response.choices[0].message.content

def chat(input: str, max_attempts: int = 100) -> str:
    messages = [
	    {"role": "system", "content": "你是 Kimi，由 Moonshot AI 提供的人工智能助手，你更擅长中文和英文的对话。你会为用户提供安全，有帮助，准确的回答。同时，你会拒绝一切涉及恐怖主义，种族歧视，黄色暴力等问题的回答。Moonshot AI 为专有名词，不可翻译成其他语言。"},
    ]
 
	# 我们将用户最新的问题构造成一个 message（role=user），并添加到 messages 的尾部
    messages.append({
		"role": "user",
		"content": input,	
	})
    st_time = time.time()  
    for i in range(max_attempts):
        print(f"Attempts: {i+1}/{max_attempts}")
        try:
            response = chat_once(messages)
            ed_time = time.time()
            print("Query Succuess!")
            print(f"Query Time: {ed_time-st_time}")
            return response
        except Exception as e:
            print(e)
            time.sleep(1)
            continue

    print("Query Failed.")
    return
 
print(chat("你好，请给我讲一个童话故事。"))