using System.Text.Json;
using Wox.Core.Utils;
using Wox.Plugin;

namespace Wox.Core.Plugin.Host;

public class NonDotnetPlugin : IPlugin
{
    public required PluginMetadata Metadata { get; init; }
    public required PluginHostBase PluginHost { get; init; }

    public void Init(PluginInitContext context)
    {
        PluginHost.InvokeMethod(Metadata, "init").Wait();
    }

    public async Task<List<Result>> Query(Query query)
    {
        var rawResults = await PluginHost.InvokeMethod(Metadata, "query", new Dictionary<string, string?>
        {
            { "RawQuery", query.RawQuery },
            { "TriggerKeyword", query.TriggerKeyword },
            { "Command", query.Command },
            { "Search", query.Search }
        });

        if (!rawResults.HasValue)
            return new List<Result>();

        var results = rawResults.Value.Deserialize<List<Result>>();
        if (results == null)
        {
            Logger.Error($"[{Metadata.Name}] Fail to deserialize query result");
            return new List<Result>();
        }

        foreach (var result in results)
            result.Action = () =>
            {
                var actionRawResult = PluginHost.InvokeMethod(Metadata, "action", new Dictionary<string, string?>
                {
                    { "ActionId", result.Id }
                }).Result;
                if (!actionRawResult.HasValue) return true;

                return actionRawResult.Value.Deserialize<bool>();
            };

        return results;
    }
}