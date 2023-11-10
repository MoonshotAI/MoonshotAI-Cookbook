using MoonshotDotnet;
using Microsoft.Extensions.Logging;
using Microsoft.Extensions.DependencyInjection;
using Newtonsoft.Json;

internal class Program
{
    public static async Task Main(string[] args)
    {
        MoonshotClient.Host = "https://api.moonshot.cn";
        MoonshotClient.ApiKey = "sk-";
        var services = new ServiceCollection();
        services.AddHttpClient();
        var serviceProvider = services.BuildServiceProvider();
        var logger = serviceProvider.GetService<ILogger<MoonshotClient>>();
        var httpClientFactory = serviceProvider.GetService<IHttpClientFactory>();
        var moonshotClient = new MoonshotClient(logger, httpClientFactory);
        Console.WriteLine("moonshotClient created");

        // ListModels example
        var modelsResp = await moonshotClient.ListModels();
        Console.WriteLine($"ListModels successfully, total: {modelsResp.data.Count}");
        foreach (var model in modelsResp.data)
        {
            Console.WriteLine($"{model.id} {model.created}");
        }
        Console.WriteLine("ListModels done");

        // Chat example
        var modelId = modelsResp.data.FirstOrDefault().id;
        var chatRep = new ChatReq
        {
            max_tokens = 2048,
            temperature = 0.5,
            top_p = 1,
            frequency_penalty = 0,
            presence_penalty = 0,
            model = modelId,
            messages = new List<MessagesItem>(){
                new MessagesItem{
                    role = "system",
                    content = "The following is a conversation with an AI assistant. The assistant is helpful, creative, clever, and very friendly."
                },
                new MessagesItem{
                    role = "user",
                    content = "What is human life expectancy in the United States?"
                }
            }
        };
        Console.WriteLine($"Chat created,request: {JsonConvert.SerializeObject(chatRep)}");
        var chatResp = await moonshotClient.Chat(chatRep);
        var chatRespBody = await chatResp.Content.ReadAsStringAsync();
        Console.WriteLine($"Chat successfully, response: {chatRespBody}");
    }
}