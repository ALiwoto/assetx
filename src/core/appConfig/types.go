package appConfig

type AssetxConfig struct {
	ProxyBaseURL         string `json:"proxy_base_url"`
	APIKey               string `json:"api_key"`
	DisableSearchHistory bool   `json:"disable_search_history"`
}
