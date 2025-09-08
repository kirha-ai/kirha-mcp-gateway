import { readFileSync } from "node:fs";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";
import { z } from "zod";

const EnvSchema = z.object({
  KIRHA_API_KEY: z.string().optional(),
  PLAN_MODE_ENABLED: z.string().optional(),
  VERTICAL_ID: z.string(),
  PORT: z.number().optional(),
});

const McpServerConfigSchema = z.object({
  name: z.string(),
  version: z.string(),
  mode: z.enum(["stdio", "http"]),
});

const ApiConfigSchema = z.object({
  summarization: z.object({
    enable: z.boolean(),
    model: z.string(),
  }),
});

export enum KihraToolNames {
  SearchKirha = "searchKirha",
  CreateKirhaSearchPlan = "createKirhaSearchPlan",
  RunKirhaSearchPlan = "runKirhaSearchPlan",
}

const ToolConfigSchema = z.object({
  name: z.nativeEnum(KihraToolNames),
  title: z.string(),
  description: z.string(),
});

export const configFileSchema = z.object({
  mcp: McpServerConfigSchema,
  api: ApiConfigSchema,
  verticals: z.array(
    z.object({
      id: z.string(),
      tools: z.array(ToolConfigSchema),
    }),
  ),
});

function getCurrentDirname() {
  try {
    if (typeof import.meta !== "undefined" && import.meta.url) {
      return dirname(fileURLToPath(import.meta.url));
    }
  } catch (e) {}

  return process.cwd();
}

const __dirname = getCurrentDirname();

function loadConfig(): ConfigFile {
  const path =
    __dirname === process.cwd() ? join(__dirname, "config.json") : join(dirname(__dirname), "config.json");

  try {
    const configData = readFileSync(path, "utf-8");
    const config = JSON.parse(configData);
    return configFileSchema.parse(config);
  } catch (error) {
    console.error(`Error loading configuration from ${path}:`, error);

    throw new Error(
      `Failed to load configuration: ${error instanceof Error ? error.message : String(error)}`,
    );
  }
}

type ConfigFile = z.infer<typeof configFileSchema>;

export type Config = {
  apiKey: string | undefined;
  port: number;
  planModeEnabled: boolean;
  verticalId: string;
  api: z.infer<typeof ApiConfigSchema>;
  mcpServer: z.infer<typeof McpServerConfigSchema>;
  tools: z.infer<typeof ToolConfigSchema>[];
};

const searchModeTools = ["searchKirha"] as const;
const planModeTools = ["createKirhaSearchPlan", "runKirhaSearchPlan"] as const;

export const config: Config = (() => {
  const parsedEnv = EnvSchema.safeParse(process.env);

  if (!parsedEnv.success) {
    const logErrorMessage = parsedEnv.error.issues
      .map((issue: z.ZodIssue) => {
        return `env variable '${issue.path.join(".")}': ${issue.message}`;
      })
      .join("\n");

    console.error(logErrorMessage);

    throw new Error("invalid environment configuration");
  }

  const configFile = loadConfig();

  const tools = configFile.verticals.find((v) => v.id === parsedEnv.data.VERTICAL_ID)?.tools;

  if (!tools) {
    throw new Error(
      `invalid configuration: No tools configuration found for vertical ID: ${parsedEnv.data.VERTICAL_ID}`,
    );
  }

  const planModeEnabled = parsedEnv.data.PLAN_MODE_ENABLED === "true";

  const expectedTools = tools.filter((tool) =>
    planModeEnabled
      ? planModeTools.includes(tool.name as (typeof planModeTools)[number])
      : searchModeTools.includes(tool.name as (typeof searchModeTools)[number]),
  );

  return {
    port: parsedEnv.data.PORT ?? 3400,
    apiKey: parsedEnv.data.KIRHA_API_KEY,
    planModeEnabled: parsedEnv.data.PLAN_MODE_ENABLED === "true",
    verticalId: parsedEnv.data.VERTICAL_ID,
    api: configFile.api,
    mcpServer: configFile.mcp,
    tools: expectedTools,
  };
})();
