import os, time
from openai import OpenAI

def call_api(messages,model_id,max_tokens=1024,temperature=0.3):
    print("query api...")
    print(model_id)
    print(messages)

    client = OpenAI(api_key=os.getenv('MOONSHOT_API_KEY'),
                    base_url=f"https://api.moonshot.cn/v1")   
        
    response = client.chat.completions.create(
        model=model_id,
        messages=messages,
        max_tokens=max_tokens,
        temperature=temperature,
    )
    return response.choices[0].message.content

def query_llm(messages,model_id="moonshot-v1-8k",max_tokens=1024,temperature=0.3):
    st_time = time.time()
    print(f"Query Model ID: {model_id}")
    max_attempts = 100
    for i in range(max_attempts):
        print(f"Try Attempts: {i+1}")
        try:
            response = call_api(messages,model_id,max_tokens,temperature)
            ed_time = time.time()
            print(f"Succuess!")
            print(f"Query Time: {ed_time-st_time}")
            return response
        except Exception as e:
            print(e)
            time.sleep(1)
            continue

    print(f"Query Failed.")
    return
