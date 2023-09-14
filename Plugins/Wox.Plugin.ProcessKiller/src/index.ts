import { PluginInitContext, PublicAPI, Query, Result, Plugin } from "@wox-launcher/wox-plugin"

let api: PublicAPI

export const plugin: Plugin = {
  init: async (context: PluginInitContext) => {
    api = context.API
    await api.Log("process killer initialized")
    await api.ShowApp()
  },

  query: async (query: Query) => {
    await api.Log("process killer got query: " + query.Search)
    return [
      {
        Title: `Kill process ${query.RawQuery}`,
        IcoPath: "Images/app.png",
        Action: async () => {
          const translationResult = await api.GetTranslation("processKillerKilling")
          await api.Log(translationResult)
          return false
        }
      }
    ] as Result[]
  }
}