package vscode

//https://microsoft.github.io/monaco-editor/docs.html#interfaces/IRange.html
type Range struct {
	StarLineNumber int
	StartColumn    int
	EndLineNumber  int
	EndColumn      int
}

// represents Monaco's IIdentifiedSingleEditOperation
//   https://microsoft.github.io/monaco-editor/docs.html#interfaces/editor.IIdentifiedSingleEditOperation.html
type EditOperation struct {
	Id    string
	Range Range
	Text  string
}

//https://microsoft.github.io/monaco-editor/docs.html#interfaces/editor.IStandaloneCodeEditor.html#executeEdits
