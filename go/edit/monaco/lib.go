package vscode

//https://microsoft.github.io/monaco-editor/docs.html#interfaces/IRange.html
type Range struct {
	StarLineNumber int
	StartColumn    int
	EndLineNumber  int
	EndColumn      int
}

// represents Monaco's ISingleEditOperation
//   https://microsoft.github.io/monaco-editor/docs.html#interfaces/editor.IIdentifiedSingleEditOperation.html
type EditOperation struct {
	Range Range
	Text  string
}

//https://microsoft.github.io/monaco-editor/docs.html#interfaces/editor.IStandaloneCodeEditor.html#executeEdits
