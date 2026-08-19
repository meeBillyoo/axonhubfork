package responses

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/looplj/axonhub/llm"
)

func TestMergeRawOnlyInputItems_NormalizesToolSearchCallStringArguments(t *testing.T) {
	structuredRaw := json.RawMessage(`[
		{"type":"message","role":"user","content":[{"type":"input_text","text":"hello"}]}
	]`)
	requestExt := &llm.OpenAIResponsesRequestExtensions{
		RawInputItems: []llm.OpenAIResponsesRawFragment{{
			Type:          "tool_search_call",
			CallID:        "call_search",
			OriginalIndex: 1,
			Raw:           json.RawMessage(`{"type":"tool_search_call","call_id":"call_search","status":"completed","arguments":"{\"query\":\"agent_kb\",\"limit\":10}"}`),
		}},
	}

	items, ok := mergeRawOnlyInputItems(structuredRaw, requestExt)
	require.True(t, ok)
	require.Len(t, items, 2)

	var rawItem map[string]any
	require.NoError(t, json.Unmarshal(items[1], &rawItem))
	require.Equal(t, "tool_search_call", rawItem["type"])

	arguments, ok := rawItem["arguments"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "agent_kb", arguments["query"])
	require.Equal(t, float64(10), arguments["limit"])
}

func TestNormalizeRawInputItem_DoesNotRewriteFunctionCallArguments(t *testing.T) {
	raw := json.RawMessage(`{"type":"function_call","call_id":"call_fn","name":"lookup","arguments":"{\"query\":\"agent_kb\",\"limit\":10}"}`)

	normalized := normalizeRawInputItem(raw)

	var item map[string]any
	require.NoError(t, json.Unmarshal(normalized, &item))
	require.Equal(t, "function_call", item["type"])

	arguments, ok := item["arguments"].(string)
	require.True(t, ok)
	require.JSONEq(t, `{"query":"agent_kb","limit":10}`, arguments)
}
