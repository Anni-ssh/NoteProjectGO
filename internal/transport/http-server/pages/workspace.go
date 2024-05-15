package pages

import (
	"NoteProject/internal/entities"
	"NoteProject/internal/transport/http-server/paths"
	"fmt"
	"html/template"
	"net/http"
)

func NoteWorkspace(w http.ResponseWriter, notes []entities.Note) error {
	tmpl, err := template.ParseFiles(paths.TemplatesPath(paths.NoteWorkSpace), paths.TemplatesPath(paths.Note))
	if err != nil {

		return fmt.Errorf("failed to parse tmpl files: %s", err)
	}

	if err = tmpl.Execute(w, notes); err != nil {
		return fmt.Errorf("failed to execute tmpl: %s", err)
	}
	return nil
}
