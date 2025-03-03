import { appendFileSync } from "node:fs";
import {getConfigMap} from "./src/config.ts";

const dotEnvPath = ".env"
const dotEnvExamplePath = ".env.example"

const dotEnv = Bun.file(dotEnvPath)
const dotEnvExample = Bun.file(dotEnvExamplePath)

const dotEnvExists = await dotEnv.exists()

const createEnvFromExample = async () => {
  console.log("Creating .env from .env.example")
  await Bun.write(dotEnvPath, await dotEnvExample.text())
}

if (dotEnvExists) {
  const dotEnvStat = await dotEnv.stat()
  if(dotEnvStat.size === 0) {
    console.log(".env file is empty")
    await createEnvFromExample()
  } else {
    const configMap = await getConfigMap()
    if(configMap.size === 0) {
      console.log("No config found")
      console.log("Append .env.example to .env")
      const dotEnvExampleContent = await dotEnvExample.text()
      appendFileSync(dotEnvPath, `\n${dotEnvExampleContent}`, "utf-8")
    }
  }
} else {
  console.log("No .env file found")
  await createEnvFromExample()
}