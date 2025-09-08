import { type ZodRawShape, z } from "zod";
import type { Config, KihraToolNames } from "./config.js";

const KIRHA_SEARCH_API_URL = "https://api.kirha.ai/chat/v1/search";

function getApiOptionsForConfig(config: Config) {
  let options = {};
  if (config.api.summarization.enable) {
    options = {
      summarization: { enable: true, model: config.api.summarization.model },
    };
  }
  return options;
}

/** Search Kirha Tool */
const searchKirhaToolInputSchema = { query: z.string() };

async function searchKirhaToolHandler({ query }: { query: string }, config: Config) {
  try {
    const options = getApiOptionsForConfig(config);
    const response = await fetch(KIRHA_SEARCH_API_URL, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${config.apiKey}`,
      },
      body: JSON.stringify({
        ...options,
        query,
        vertical_id: config.verticalId,
      }),
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const result = await response.json();
    return { content: [{ type: "text" as const, text: JSON.stringify(result) }] };
  } catch (error) {
    return {
      content: [
        { type: "text" as const, text: `Error: ${error instanceof Error ? error.message : String(error)}` },
      ],
    };
  }
}

/** Create Kirha Search Plan Tool */
const createKirhaSearchPlanInputSchema = { query: z.string() };

async function createKirhaSearchPlanHandler({ query }: { query: string }, config: Config) {
  try {
    const response = await fetch(`${KIRHA_SEARCH_API_URL}/plan`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${config.apiKey}`,
      },
      body: JSON.stringify({
        query,
        vertical_id: config.verticalId,
      }),
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const result = await response.json();
    return { content: [{ type: "text" as const, text: JSON.stringify(result) }] };
  } catch (error) {
    return {
      content: [
        { type: "text" as const, text: `Error: ${error instanceof Error ? error.message : String(error)}` },
      ],
    };
  }
}

/** Run Kirha Search Plan Tool */
const runKirhaSearchPlanInputSchema = { planId: z.string() };

async function runKirhaSearchPlanHandler({ planId }: { planId: string }, config: Config) {
  try {
    const options = getApiOptionsForConfig(config);
    const response = await fetch(`${KIRHA_SEARCH_API_URL}/plan/run`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${config.apiKey}`,
      },
      body: JSON.stringify({
        ...options,
        plan_id: planId,
        summarization: { enable: true, model: "kirha-flash" },
      }),
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const result = await response.json();
    return { content: [{ type: "text" as const, text: JSON.stringify(result) }] };
  } catch (error) {
    return {
      content: [
        { type: "text" as const, text: `Error: ${error instanceof Error ? error.message : String(error)}` },
      ],
    };
  }
}

/** Run Kirha Search Plan Tool */

export type ToolDefinition = {
  inputSchema: ZodRawShape;
  // biome-ignore lint/suspicious/noExplicitAny: <_explanation>
  handler: (input: any, config: Config) => Promise<{ content: { type: "text"; text: string }[] }>;
};

export const kirhaToolDefinitions: Record<KihraToolNames, ToolDefinition> = {
  searchKirha: {
    inputSchema: searchKirhaToolInputSchema,
    handler: searchKirhaToolHandler,
  },
  createKirhaSearchPlan: {
    inputSchema: createKirhaSearchPlanInputSchema,
    handler: createKirhaSearchPlanHandler,
  },
  runKirhaSearchPlan: {
    inputSchema: runKirhaSearchPlanInputSchema,
    handler: runKirhaSearchPlanHandler,
  },
};
