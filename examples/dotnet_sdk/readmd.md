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

## build

```bash

dotnet build MoonshotClient
```
