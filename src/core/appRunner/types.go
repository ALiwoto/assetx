package appRunner

type ImageRequest struct {
	Avoids       []string
	Background   string
	ConfigPath   string
	ExampleNotes []string
	Examples     []string
	Model        string
	OutputPath   string
	Prompt       string
	Quality      string
	Size         string
}

type ConvertTGSRequest struct {
	FFMPEGPath string
	InputPath  string
	OutputPath string
}

type ConvertWEBPRequest struct {
	InputPath  string
	OutputPath string
}

type RemoveBackgroundRequest struct {
	Color      string
	InputPath  string
	OutputPath string
	Tolerance  int
}

type SearchRequest struct {
	AllowedDomains    []string
	ConfigPath        string
	Model             string
	Query             string
	SearchContextSize string
}
