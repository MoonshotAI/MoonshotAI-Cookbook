using MoonshotDotnet;
using Microsoft.Extensions.Logging;
using Microsoft.Extensions.DependencyInjection;

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
        var modelsResp = await moonshotClient.ListModels();
        Console.WriteLine($"ListModels successfully, total: {modelsResp.data.Count}");
        foreach (var model in modelsResp.data)
        {
            Console.WriteLine($"{model.id} {model.created}");
        }
        Console.WriteLine("ListModels done");
    }
}