package utilities

import (
	"database/sql"
	"fmt"
	"html"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"text/template"
)

var Template *template.Template

var TemplateFuncs = template.FuncMap{
    "add": func(a, b int) int { return a + b },
    "toUpper": strings.ToUpper,
    "hasPrefix": strings.HasPrefix,
    "contains": func(arr []string, value string) bool {
        return slices.Contains(arr, value)
    },
    "containsByField": func(list any, fieldName string, value any) bool {
        v := reflect.ValueOf(list)

        // If it's a pointer, dereference it
        if v.Kind() == reflect.Ptr {
            v = v.Elem()
        }

        // Check list is a slice
        if v.Kind() != reflect.Slice {
            return false
        }

        for i := range v.Len() {
            elem := v.Index(i)
    
            // If element is a pointer, dereference
            if elem.Kind() == reflect.Ptr {
                elem = elem.Elem()
            }
    
            if elem.Kind() != reflect.Struct {
                continue
            }
    
            field := elem.FieldByName(fieldName)
            if !field.IsValid() {
                continue
            }
    
            if fmt.Sprint(field.Interface()) == fmt.Sprint(value) {
                return true
            }
        }

        return false
    },
    "default": func(value, fallback string) string {
        if value == "" {
            return fallback
        }
        return value
    },
    "formatDate": func(v any, layout, fallback string) string {
        if t, ok := v.(sql.NullTime); ok && t.Valid {
            return t.Time.Format(layout)
        }
        return fallback
    },
}

func LoadTemplates() *template.Template {
    root := template.New("").Option("missingkey=error").Funcs(TemplateFuncs)

    filepath.Walk("views", func(path string, info os.FileInfo, err error) error {
        if err != nil || info.IsDir() {
            return nil
        }

        if strings.HasSuffix(path, ".html") {
            rel := strings.TrimPrefix(filepath.ToSlash(path), "views/")
			
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

    // 
    tmpl := template.Must(Template.Clone())
    bytes, _ := os.ReadFile("views/" + name)
    tmpl = template.Must(tmpl.Parse(string(bytes)))

	return tmpl.ExecuteTemplate(w, name, payload)
}

func EscapeHTML(s string) string {
    return html.EscapeString(s)
}

func UnescapeHTML(s string) string {
    return html.UnescapeString(s)
}

func Compact(pairs ...any) map[string]any {
    m := make(map[string]any)

    for i := 0; i < len(pairs); i += 2 {
        key := pairs[i].(string)
        var val any
        if i+1 < len(pairs) {
            val = pairs[i+1]
        } else {
            val = nil
        }
        m[key] = val
    }
    return m
}