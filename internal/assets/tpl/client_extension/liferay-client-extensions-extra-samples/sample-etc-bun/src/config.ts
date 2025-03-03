import {Glob} from "bun"
import {log} from "./log.ts";

const configTreePaths = [
  Bun.env.LIFERAY_ROUTES_CLIENT_EXTENSION,
  Bun.env.LIFERAY_ROUTES_DXP
]

export const getConfigMap = async () => {
  let configMap = new Map<string, string>([])
  const glob = new Glob("*")
  for (const configTreePath of configTreePaths) {
    if(configTreePath) {
      log.info(configTreePath)
      for (const fileName of glob.scanSync(configTreePath)) {
        const file = Bun.file(`${configTreePath}/${fileName}`)
        configMap.set(fileName, await file.text())
      }
    }
  }
  return configMap
}