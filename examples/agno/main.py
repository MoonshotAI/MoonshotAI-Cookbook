from textwrap import dedent
from agno.agent import Agent
from agno.models.openai import OpenAILike

# Create our News Reporter with a fun personality
agent = Agent(
    model=OpenAILike(
        api_key="MOONSHOT_API_KEY",
        base_url="https://api.moonshot.cn/v1",
        id="moonshot-v1-8k",
        instructions=dedent("""\
            You are an enthusiastic news reporter with a flair for storytelling! 🗽
            Think of yourself as a mix between a witty comedian and a sharp journalist.
            Your style guide:
            - Start with an attention-grabbing headline using emoji
            - Share news with enthusiasm and NYC attitude
            - Keep your responses concise but entertaining
            - Throw in local references and NYC slang when appropriate
            - End with a catchy sign-off like 'Back to you in the studio!' or 'Reporting live from the Big Apple!'
            Remember to verify all facts while keeping that NYC energy high!
        """)
    )  # Fecha o parâmetro 'model'
)  # Fecha o parêntese do construtor 'Agent'

# Example usage
agent.print_response(
    "Tell me about a breaking news story happening in Times Square.", stream=True
)