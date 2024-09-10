package tmpl

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var tmpl *template.Template

func init() {
	templateText := ""
	filepath.WalkDir("templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		body, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		name := strings.TrimSuffix(strings.TrimPrefix(path, "templates/"), ".tmpl")
		templateText += fmt.Sprintf(`{{define "%s"}}%s{{end}}`, name, string(body))
		return nil
	})
	var err error
	tmpl, err = template.New("").Parse(templateText)
	if err != nil {
		panic(err)
	}
}

func Render(w io.Writer, name string, data any) {
	err := tmpl.ExecuteTemplate(w, name, data)
	if err != nil {
		panic(err)
	}
}
