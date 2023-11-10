using System;
using System.Text;
using System.Threading.Tasks;
using System.Net.Http.Headers;
using System.Configuration;
using System.Net.Http;
using System.IO;
using Microsoft.Extensions.Logging;
using Newtonsoft.Json.Linq;

namespace MoonshotDotnet
{

    /// <summary>
    /// https://platform.moonshot.cn/docs/api-reference
    /// </summary>
    public class MoonshotClient
    {
        private readonly ILogger<MoonshotClient> _logger;

        private readonly IHttpClientFactory _httpClientFactory;

        public MoonshotClient(ILogger<MoonshotClient> logger, IHttpClientFactory httpClientFactory)
        {
            _logger = logger;
            _httpClientFactory = httpClientFactory;
        }

        /// <summary>
        ///  list models
        /// </summary>
        /// <returns></returns>
        public async Task<ModelListResp> ListModels()
        {
            var response = await GetAsync("/v1/models");
            return await ParseResp<ModelListResp>(response);
        }

        /// <summary>
        /// Chat
        /// </summary>
        /// <param name="requestBody"></param>
        /// <returns>Return HttpResponseMessage for SSE</returns>
        public async Task<HttpResponseMessage> Chat(string requestBody)
        {
            return await PostJsonStreamAsync("/v1/chat/completions", requestBody);
        }

        /// <summary>
        /// Chat
        /// </summary>
        /// <param name="chatReq"></param>
        /// <returns>Return HttpResponseMessage for SSE</returns>
        public async Task<HttpResponseMessage> Chat(ChatReq chatReq)
        {
            var requestBody = Newtonsoft.Json.JsonConvert.SerializeObject(chatReq);
            return await PostJsonStreamAsync("/v1/chat/completions", requestBody);
        }

        /// <summary>
        ///  Get as timate token count
        /// </summary>
        public async Task<int?> GetAsTiMateTokenCount(string chatReqText)
        {
            var response = await PostJsonAsync("/v1/tokenizers/estimate-token-count", chatReqText);
            var responseText = await response.Content.ReadAsStringAsync();
            if (response.IsSuccessStatusCode)
            {
                var responseObj = JToken.Parse(responseText);
                return responseObj?["data"]?["total_tokens"]?.ToObject<int>();
            }
            var error = Newtonsoft.Json.JsonConvert.DeserializeObject<ErrorResponse>(responseText);
            _logger.LogError($"{error?.error?.type}: {error?.error?.message}");
            throw new Exception($"{error?.error.type}: {error?.error.message}");
        }


        /// <summary>
        /// Get as timate token count
        /// </summary>
        /// <param name="chatReq"></param>
        /// <returns></returns>
        public async Task<int?> GetAsTiMateTokenCount(ChatReq chatReq)
        {
            var chatReqText = Newtonsoft.Json.JsonConvert.SerializeObject(chatReq);
            return await GetAsTiMateTokenCount(chatReqText);
        }


        /// <summary>
        ///  List files
        /// </summary>
        public virtual async Task<FileListResp> ListFiles()
        {
            var response = await GetAsync("/v1/files");
            return await ParseResp<FileListResp>(response);
        }





        /// <summary>
        ///  Upload file
        /// </summary>
        public virtual async Task<FileItem> UploadFile(string filePath)
        {
            if (!File.Exists(filePath))
            {
                throw new FileNotFoundException($"{filePath} not found");
            }
            var client = _httpClientFactory.CreateClient();
            client.DefaultRequestHeaders.Authorization = new AuthenticationHeaderValue("Bearer", ApiKey);
            var request = new HttpRequestMessage(HttpMethod.Post, $"{Host}/v1/files");
            var content = new MultipartFormDataContent
            {
                { new StreamContent(File.OpenRead(filePath)), "file", filePath }
            };
            request.Content = content;
            var response = await client.SendAsync(request);
            return await ParseResp<FileItem>(response);
        }



        /// <summary>
        ///  Upload file stream
        /// </summary>
        public virtual async Task<FileItem> UploadFileStream(Stream stream, string fileName)
        {
            var client = _httpClientFactory.CreateClient();
            client.DefaultRequestHeaders.Authorization = new AuthenticationHeaderValue("Bearer", ApiKey);
            var request = new HttpRequestMessage(HttpMethod.Post, $"{Host}/v1/files");
            var content = new MultipartFormDataContent
            {
                { new StreamContent(stream), "file", fileName }
            };
            request.Content = content;
            var response = await client.SendAsync(request);
            return await ParseResp<FileItem>(response);
        }


        /// <summary>
        ///  Get file content
        /// </summary>

        public virtual async Task<FileContent> GetFileContent(string fileId)
        {
            var response = await GetAsync($"/v1/files/{fileId}/content");
            return await ParseResp<FileContent>(response);
        }


        private async Task<HttpResponseMessage> GetAsync(string path)
        {
            var client = _httpClientFactory.CreateClient();
            client.DefaultRequestHeaders.Authorization = new AuthenticationHeaderValue("Bearer", ApiKey);
            return await client.GetAsync(Host + path);
        }

        private async Task<HttpResponseMessage> PostJsonAsync(string path, string json)
        {
            var client = _httpClientFactory.CreateClient();
            client.DefaultRequestHeaders.Authorization = new AuthenticationHeaderValue("Bearer", ApiKey);
            return await client.PostAsync(Host + path, new StringContent(json, Encoding.UTF8, "application/json"));
        }

        private async Task<HttpResponseMessage> PostJsonStreamAsync(string path, string json)
        {
            var client = _httpClientFactory.CreateClient();
            client.DefaultRequestHeaders.Authorization = new AuthenticationHeaderValue("Bearer", ApiKey);
            var request = ToHttpRequest(path);
            request.Content = new StringContent(json, Encoding.UTF8, "application/json");
            return await client.SendAsync(request, HttpCompletionOption.ResponseHeadersRead);
        }

        private HttpRequestMessage ToHttpRequest(string path)
        {
            var request = new HttpRequestMessage();
            var uriBuilder = new UriBuilder(Host + path);
            request.RequestUri = uriBuilder.Uri;
            request.Method = new HttpMethod("POST");
            request.Headers.Host = (new Uri(Host)).Host;
            return request;
        }



        /// <summary>
        /// Parse response
        /// </summary>
        private async Task<T> ParseResp<T>(HttpResponseMessage response)
        {
            var responseText = await response.Content.ReadAsStringAsync();
            if (response.IsSuccessStatusCode)
            {
                return Newtonsoft.Json.JsonConvert.DeserializeObject<T>(responseText) ?? default;
            }
            var error = Newtonsoft.Json.JsonConvert.DeserializeObject<ErrorResponse>(responseText);
            _logger.LogError($"{error?.error.type}: {error?.error.message}");
            throw new Exception($"{error?.error.type}: {error?.error.message}");
        }



        private static string _host = "https://api.moonshot.cn";

        public static string Host
        {
            get
            {
                if (string.IsNullOrEmpty(_host) && !string.IsNullOrEmpty(ConfigurationManager.AppSettings["MoonshotHost"]))
                {
                    _host = ConfigurationManager.AppSettings?["MoonshotHost"] ?? "";
                }

                return _host;
            }
            set
            {

                _host = value;
            }
        }


        private static string _apiKey = "sk_";

        public static string ApiKey
        {
            get
            {
                if (string.IsNullOrEmpty(_apiKey) && !string.IsNullOrEmpty(ConfigurationManager.AppSettings["MoonshotApiKey"]))
                {
                    _apiKey = ConfigurationManager.AppSettings["MoonshotApiKey"] ?? "";
                }

                return _apiKey;
            }
            set
            {
                _apiKey = value;
            }
        }

    }
}