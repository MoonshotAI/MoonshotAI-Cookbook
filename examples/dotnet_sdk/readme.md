# Moonshot SDK for .NET Core 8.0

Moonshot SDK for .NET is a library that allows you to interact with the Moonshot API from your .NET application.

- [api-reference](https://platform.moonshot.cn/docs/api-reference)
- dotnet core 8.0 sdk

## package

```xml
  <ItemGroup>
    <PackageReference Include="Microsoft.Extensions.Http" Version="8.0.0" />
    <PackageReference Include="Newtonsoft.Json" Version="13.0.3" />
    <PackageReference Include="System.Configuration.ConfigurationManager" Version="8.0.0" />
  </ItemGroup>
```

## src

- [MoonshotClient.cs](MoonshotClient/MoonshotClient.cs): Moonshot API client.
- [MoonshotModel.cs](MoonshotClient/MoonshotModel.cs): Moonshot API model.

## usage

Example usage of the Moonshot SDK for .NET:

```csharp

MoonshotClient.ApiKey = "your-api"; // set your api key

// var logger = LoggerFactory.Create(builder => builder.AddConsole()).CreateLogger<MoonshotClient>();
// httpClientFactory = IHttpClientFactory
var client = new MoonshotClient(logger, httpClientFactory);

```

if your use AutoFac, you can use the following code:

```csharp

// Startup.cs ConfigureServices()

services.AddHttpClient("ProxyClient", c =>
{
    // Configure other settings of HttpClient
})
// var builder = new ContainerBuilder();

builder.RegisterType<MoonshotClient>().As<MoonshotClient>().InstancePerLifetimeScope();


```

```csharp

// HealthController.cs

    [ApiController]
    [Route("[controller]")]
    public class HealthController : ControllerBase
    {

        private readonly MoonshotClient _moonshotClient;

        public HealthController(MoonshotClient moonshotClient)
        {
            _moonshotClient = moonshotClient;
        }


        [HttpGet("/headers")]
        public IActionResult GetHeaders()
        {
            return Ok(new
            {
                data = HttpContext.Request.Headers,
                code = 0
            });
        }

        [HttpGet("/moonshot/test")]
        public async Task<IActionResult> GetMoonshot()
        {
            var modelList = await _moonshotClient.ListModels();
            var fileList = await _moonshotClient.ListFiles();
            return Ok(new
            {
                models = modelList,
                files = fileList,
                code = 0
            });
        }

    }


```

## demo

```csharp

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
```

## build

```bash

dotnet build MoonshotClient
```
