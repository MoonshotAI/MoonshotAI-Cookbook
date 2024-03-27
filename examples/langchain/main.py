from langchain_openai import ChatOpenAI

llm = ChatOpenAI(
    openai_api_base="https://api.moonshot.cn/v1/", 
    openai_api_key="MOONSHOT_API_KEY",
    model_name="moonshot-v1-8k",
)

print(llm.invoke("how can langsmith help with testing?"))