package paths

import "path/filepath"

const (
	web       = "web"
	templates = "templates"
	Header    = "header.html"

	Home          = "home.html"
	Note          = "note.html"
	NoteWorkSpace = "note-workspace.html"
)

func TemplatesPath(path string) string {
	return filepath.Join(web, templates, path)
}
