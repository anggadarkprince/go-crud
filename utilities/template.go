package utilities

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var Template *template.Template

func LoadTemplates() *template.Template {
    root := template.New("").Option("missingkey=error").Funcs(TemplateFuncs)

    filepath.Walk("views", func(path string, info os.FileInfo, err error) error {
        if err != nil || info.IsDir() {
            return nil
        }

        if strings.HasSuffix(path, ".html") {
            rel := strings.TrimPrefix(path, "views/")
			
            bytes, _ := os.ReadFile(path)

            _, err := root.New(rel).Parse(string(bytes))
            if err != nil {
                panic(err)
            }
        }
        return nil
    })

    return root
}

func InitTemplates() {
	Template = LoadTemplates()

	for _, t := range Template.Templates() {
		fmt.Println("Loaded template:", t.Name())
	}
}

func Render(w http.ResponseWriter, r *http.Request, name string, data any) error {
	funcs := template.FuncMap{
        "query": func(key string) string {
            return r.URL.Query().Get(key)
        },
        "header": func(key string) string {
            return r.Header.Get(key)
        },
		"currentPath": func() string {
			return r.URL.Path
		},
    }

    // Clone template so funcs are local to this request
    Template, err := Template.Clone()
    if err != nil {
        return err
    }

    Template = Template.Funcs(funcs)

	payload := map[string]any{
        "currentPath": r.URL.Path,
    }

	// If original data is map, merge it
    if m, ok := data.(map[string]any); ok {
        for k, v := range m {
            payload[k] = v
        }
    } else if data != nil {
        // Otherwise put struct under "Data"
        payload["data"] = data
    }

	return Template.ExecuteTemplate(w, name, payload)
}
