#include <iostream>
#include "openai/openai.hpp"

int main() {
    openai::start(
        "$YOUR_API_KEY",
        "", // organization id
        true, // throw exception
        "https://api.moonshot.cn/v1/"  // base url
    );

    auto chat = openai::chat().create(R"(
        {
            "model": "moonshot-v1-8k",
            "messages":[{"role":"user", "content":"Say hello world."}] 
        }
    )"_json);
    std::cout << "Response is:\n" << chat.dump(2) << '\n'; 
}