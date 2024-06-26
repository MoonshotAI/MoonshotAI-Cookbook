require 'net/http'
require 'json'
require 'uri'

# 设置 API 密钥和请求的参数
API_KEY = 'sk-YOUR_API_KEY'
ENDPOINT = 'https://api.moonshot.cn/v1/chat/completions'

model = 'moonshot-v1-32k'
prompt = "帮我写一段Ruby代码"

# 构建请求体
request_body = {
  model: model,
  max_tokens: 1024,
  temperature: 0.7,
  messages: [
    { role: "user", content: prompt }
  ]
}.to_json

# 设置请求的 URL
url = URI.parse(ENDPOINT)

# 创建一个 Net::HTTP::Post 对象，设置请求的参数
request = Net::HTTP::Post.new(url.path)
request['Content-Type'] = 'application/json'
request['Authorization'] = "Bearer #{API_KEY}"
request.body = request_body

# 发送 POST 请求并获取响应
response = Net::HTTP.start(url.host, url.port, use_ssl: url.scheme == 'https') do |http|
  http.request(request)
end

# 打印响应结果
if response.is_a?(Net::HTTPSuccess)
  puts response.body
else
  puts "请求失败: #{response.code} #{response.message}"
end
