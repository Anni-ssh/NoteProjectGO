package pages

import (
	"NoteProject/internal/transport/http-server/paths"
	"fmt"
	"html/template"
	"net/http"
)

func Home(w http.ResponseWriter) error {
	tmpl, err := template.ParseFiles(paths.TemplatesPath(paths.Home))
	if err != nil {
		return fmt.Errorf("failed to parse tmpl files: %s", err)
	}

	if err = tmpl.Execute(w, nil); err != nil {
		return fmt.Errorf("failed to execute tmpl: %s", err)
	}
	return nil

}
