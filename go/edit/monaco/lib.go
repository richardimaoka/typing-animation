package monaco

// https://microsoft.github.io/monaco-editor/docs.html#interfaces/IRange.html
type Range struct {
	StartColumn     int `json:"startColumn"`
	StartLineNumber int `json:"startLineNumber"`
	EndColumn       int `json:"endColumn"`
	EndLineNumber   int `json:"endLineNumber"`
}

// https://microsoft.github.io/monaco-editor/docs.html#interfaces/editor.ISingleEditOperation.html#range
type SingleEditOperation struct {
	Text      string `json:"text"`
	Range     Range  `json:"range"`
	Operation string `json:"operation"`
}
