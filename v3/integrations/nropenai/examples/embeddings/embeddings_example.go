package main

import (
	"fmt"
	"os"
	"time"

	"github.com/newrelic/go-agent/v3/integrations/nropenai"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/pkoukk/tiktoken-go"
	openai "github.com/sashabaranov/go-openai"
)

func main() {
	// Start New Relic Application
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName("Basic OpenAI App"),
		newrelic.ConfigLicense(os.Getenv("NEW_RELIC_LICENSE_KEY")),
		newrelic.ConfigDebugLogger(os.Stdout),
		// Enable AI Monitoring
		newrelic.ConfigAIMonitoringEnabled(true),
	)
	if nil != err {
		panic(err)
	}
	app.WaitForConnection(10 * time.Second)
	app.SetLLMTokenCountCallback(func(modelName string, content string) int {
		var tokensPerMessage, tokensPerName int
		switch modelName {
		case "gpt-3.5-turbo-0613",
			"gpt-3.5-turbo-16k-0613",
			"gpt-4-0314",
			"gpt-4-32k-0314",
			"gpt-4-0613",
			"gpt-4-32k-0613":
			tokensPerMessage = 3
			tokensPerName = 1
		case "gpt-3.5-turbo-0301":
			tokensPerMessage = 4
			tokensPerName = -1
		}

		tkm, err := tiktoken.EncodingForModel(modelName)
		if err != nil {
			fmt.Println("error getting tokens", err)
			return 0
		}
		token := tkm.Encode(content, nil, nil)
		totalTokens := len(token) + tokensPerMessage + tokensPerName
		return totalTokens
	})
	// OpenAI Config - Additionally, NRDefaultAzureConfig(apiKey, baseURL string) can be used for Azure
	cfg := nropenai.NRDefaultConfig(os.Getenv("OPEN_AI_API_KEY"))

	// Create OpenAI Client - Additionally, NRNewClient(authToken string) can be used
	client := nropenai.NRNewClientWithConfig(cfg)

	// Add any custom attributes
	// NOTE: Attributes must start with "llm.", otherwise they will be ignored
	client.CustomAttributes = map[string]interface{}{
		"llm.foo": "bar",
		"llm.pi":  3.14,
	}

	fmt.Println("Creating Embedding Request...")
	// Create Embeddings
	embeddingReq := openai.EmbeddingRequest{
		Input: []string{
			"The food was delicious and the waiter",
			"Other examples of embedding request",
		},
		Model:          openai.AdaEmbeddingV2,
		EncodingFormat: openai.EmbeddingEncodingFormatFloat,
	}
	resp, err := nropenai.NRCreateEmbedding(client, embeddingReq, app)
	if err != nil {
		panic(err)
	}

	fmt.Println("Embedding Created!")
	fmt.Println(resp.Usage.PromptTokens)
	// Shutdown Application
	app.Shutdown(5 * time.Second)
}
