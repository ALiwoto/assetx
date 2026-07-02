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
