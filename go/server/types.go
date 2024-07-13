package server

type Position struct {
	Line      int `json:"line"`      //The zero-based line value.
	Character int `json:"character"` //The zero-based character value.
}

type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

type EditInsert struct {
	NewText  string   `json:"newText"`
	Position Position `json:"position"`
}

type EditDelete struct {
	DeleteRange Range `json:"deleteRange"`
}

type Edit struct {
	EditType string `json:"editType"`
	// Either of Insert or Delete must be non-nil, dependent on EditType
	Insert *EditInsert `json:"insert"`
	Delete *EditDelete `json:"delete"`
}

type NextTransition struct {
	Edits []Edit `json:"edits"`
}

type FileData struct {
	CommitHash string         `json:"commitHash"`
	FilePath   string         `json:"filePath"`
	Contents   string         `json:"contents"`
	Next       NextTransition `json:"nextTransaction"`
}
