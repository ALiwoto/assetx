package openaiClient

import "net/http"

type APIErrorResponse struct {
	Error APIErrorBody `json:"error"`
}

type APIErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Param   string `json:"param"`
	Type    string `json:"type"`
}

type Client struct {
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client
}

type ImageData struct {
	B64JSON       string `json:"b64_json"`
	RevisedPrompt string `json:"revised_prompt"`
	URL           string `json:"url"`
}

type ImageRequest struct {
	Background   string
	Examples     []string
	Model        string
	OutputFormat string
	Prompt       string
	Quality      string
	Size         string
}

type ImageResponse struct {
	Data []ImageData `json:"data"`
}
