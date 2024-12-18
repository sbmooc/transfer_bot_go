package openai

type OpenAIClient struct {
	api_key string
	url     string
}

func NewOpenAIClient() *OpenAIClient {
	return &OpenAIClient{
		url:     ,
		api_key: token,
	}
}
