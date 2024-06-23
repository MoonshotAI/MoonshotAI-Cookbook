const axios = require('axios')
const openai = require('openai')
const fs = require('fs')

const client = new openai.OpenAI({
  apiKey: '$MOONSHOT_API_KEY',
  baseURL: 'https://api.moonshot.cn/v1'
})

// 删除缓存项的函数
async function deleteCacheItem (id) {
  try {
    const response = await fetch(`https://api.moonshot.cn/v1/caching/${id}`, {
      method: 'DELETE',
      headers: {
        Authorization: 'Bearer $MOONSHOT_API_KEY'
      }
    })

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }

    const result = await response.json()
    console.log(result)
  } catch (error) {
    console.error('Error deleting cache item:', error)
  }
}

// 清除所有缓存项的函数
async function clearCache () {
  try {
    const response = await fetch('https://api.moonshot.cn/v1/caching', {
      method: 'GET',
      headers: {
        Authorization: 'Bearer $MOONSHOT_API_KEY'
      }
    })

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }

    const data = await response.json()
    const cacheItems = data.data || [] // 确保data对象中有data属性

    for (const item of cacheItems) {
      await deleteCacheItem(item.id)
    }
  } catch (error) {
    console.error('Error clearing cache:', error)
  }
}


function readFile(path) {
  try {
    return fs.readFileSync(path, 'utf8');
  } catch (err) {
    console.error(err);
    throw err;
  }
}

async function create () {
  try {
    const content = readFile('../西游记.txt')
    const response = await axios.post('https://api.moonshot.cn/v1/caching',      
      {
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
      },
      {
        headers: {
          Authorization: `Bearer $MOONSHOT_API_KEY`
        }
      }
    )
    console.log(response.data)
    return response.data.id
  } catch (error) {
    console.error('Error:', error)
  }
}

async function queryWithCache (query, cacheId) {
  try {
    const response = await client.chat.completions.create({
      model: 'moonshot-v1-128k',
      messages: [
        {
          role: 'cache',
          content: `cache_id=${cacheId}`
        },
        {
          role: 'user',
          content: query
        }
      ],
      temperature: 0.3
    })
    console.log(response.choices[0].message)
  } catch (error) {
    console.error('An error occurred:', error)
  }
}

clearCache().then(() => {
  create().then(cache_id => {
    queryWithCache('大闹天宫篇中孙悟空打败了哪些天兵天将？', cache_id)
  })
})
