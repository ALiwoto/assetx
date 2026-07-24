package openaiClient

type WebSearchRequest struct {
	AllowedDomains    []string
	Model             string
	Query             string
	SearchContextSize string
}

type WebSearchResult struct {
	Sources []WebSearchSource
	Text    string
}

type WebSearchSource struct {
	Title string `json:"title"`
	Type  string `json:"type"`
	URL   string `json:"url"`
}

type responsesCreateRequest struct {
	Include    []string        `json:"include"`
	Input      string          `json:"input"`
	Model      string          `json:"model"`
	ToolChoice string          `json:"tool_choice"`
	Tools      []webSearchTool `json:"tools"`
}

type webSearchTool struct {
	Filters           *webSearchFilters `json:"filters,omitempty"`
	SearchContextSize string            `json:"search_context_size"`
	Type              string            `json:"type"`
}

type webSearchFilters struct {
	AllowedDomains []string `json:"allowed_domains"`
}

type responsesAPIResponse struct {
	Error  *responsesAPIError    `json:"error"`
	Output []responsesOutputItem `json:"output"`
	Status string                `json:"status"`
}

type responsesAPIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type responsesOutputItem struct {
	Action  *webSearchAction  `json:"action"`
	Content []responseContent `json:"content"`
	Type    string            `json:"type"`
}

type webSearchAction struct {
	Sources []WebSearchSource `json:"sources"`
}

type responseContent struct {
	Annotations []responseAnnotation `json:"annotations"`
	Text        string               `json:"text"`
	Type        string               `json:"type"`
}

type responseAnnotation struct {
	Title string `json:"title"`
	Type  string `json:"type"`
	URL   string `json:"url"`
}
