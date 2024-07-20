package monaco

//https://microsoft.github.io/monaco-editor/docs.html#interfaces/editor.IIdentifiedSingleEditOperation.html

// https://microsoft.github.io/monaco-editor/docs.html#interfaces/IRange.html
type Range struct {
	StartColumn     int `json:"startColumn"`
	StartLineNumber int `json:"startLineNumber"`
	EndColumn       int `json:"endColumn"`
	EndLineNumber   int `json:"endLineNumber"`
}

// https://microsoft.github.io/monaco-editor/docs.html#interfaces/editor.ISingleEditOperation.html#range
type SingleEditOperation struct {
	Range Range  `json:"range"`
	Text  string `json:"text"`
}
