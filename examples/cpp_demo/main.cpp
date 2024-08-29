#include <iostream>
#include "restclient-cpp/restclient.h"
using namespace std;

int main() {
    RestClient::Connection* conn = new RestClient::Connection("https://api.moonshot.cn/v1/chat/completions");
    conn->AppendHeader("Content-Type", "application/json");
    conn->AppendHeader("Authorization", "Bearer $MOONSHOT_API_KEY");
    RestClient::Response resp = conn->post("/post", "{ \"model\": \"moonshot-v1-8k\", \"messages\": [ {\"role\": \"system\", \"content\": \"你是 Kimi，由 Moonshot AI 提供的人工智能助手，你更擅长中文和英文的对话。你会为用户提供安全，有帮助，准确的回答。同时，你会拒绝一切涉及恐怖主义，种族歧视，黄色暴力等问题的回答。Moonshot AI 为专有名词，不可翻译成其他语言。\"}, {\"role\": \"user\", \"content\": \"你好，我叫李雷，1+1等于多少？\"}],");
    cout << resp.code << endl;
    cout << resp.body << endl;
    cout << resp.headers << endl;
}