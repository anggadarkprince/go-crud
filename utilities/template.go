package utilities

import (
	"database/sql"
	"fmt"
	"html"
	"maps"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"text/template"

	"github.com/anggadarkprince/crud-employee-go/middlewares"
	"github.com/anggadarkprince/crud-employee-go/utilities/session"
)

var Template *template.Template

var TemplateFuncs = template.FuncMap{
    "add": func(a, b int) int { return a + b },
    "toUpper": strings.ToUpper,
    "hasPrefix": strings.HasPrefix,
    "contains": func(arr any, value string) bool {
        if arr == nil {
            return false
        }
        
        // Handle []string
        if strArr, ok := arr.([]string); ok {
            return slices.Contains(strArr, value)
        }
        
        // Handle []any
        if ifaceArr, ok := arr.([]any); ok {
            for _, item := range ifaceArr {
                if str, ok := item.(string); ok && str == value {
                    return true
                }
            }
            return false
        }
        
        // Handle using reflection for other slice types
        v := reflect.ValueOf(arr)
        if v.Kind() == reflect.Slice {
            for i := range v.Len() {
                item := v.Index(i).Interface()
                if str, ok := item.(string); ok && str == value {
                    return true
                }
            }
        }
        
        return false
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
    "pluck": func(arr any, fieldName string) []string {
        if arr == nil {
            return []string{}
        }
        
        result := []string{}
        v := reflect.ValueOf(arr)
        
        // Dereference pointer if the input is a pointer
        if v.Kind() == reflect.Ptr {
            if v.IsNil() {
                return result
            }
            v = v.Elem()
        }
        
        // Handle if not a slice
        if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
            return result
        }
        
        for i := 0; i < v.Len(); i++ {
            item := v.Index(i)
            
            // Dereference pointer if needed
            if item.Kind() == reflect.Ptr {
                if item.IsNil() {
                    continue
                }
                item = item.Elem()
            }
            
            var fieldValue string
            
            // Handle struct
            if item.Kind() == reflect.Struct {
                field := item.FieldByName(fieldName)
                if field.IsValid() {
                    fieldValue = fmt.Sprintf("%v", field.Interface())
                }
            }
            
            // Handle map
            if item.Kind() == reflect.Map {
                field := item.MapIndex(reflect.ValueOf(fieldName))
                if field.IsValid() {
                    fieldValue = fmt.Sprintf("%v", field.Interface())
                }
            }
            
            if fieldValue != "" {
                result = append(result, fieldValue)
            }
        }
        
        return result
    },
    "default": func(value any, fallback any) any {
        if value == nil {
            return fallback
        }
        
        // Check if it's a slice/array and if it's empty
        v := reflect.ValueOf(value)
        if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
            if v.Len() == 0 {
                return fallback
            }
            return value
        }
        
        str := fmt.Sprintf("%v", value)
        if str == "" || str == "<no value>" || str == "[]" {
            return fallback
        }
        return value
    },
    "emptySlice": func() []string {
        return []string{}
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

func Render(w http.ResponseWriter, r *http.Request, name string, data map[string]any) error {
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
    tmpl, err := Template.Clone()
    if err != nil {
        return err
    }

    tmpl = tmpl.Funcs(funcs)

    //tmpl := template.Must(Template.Clone())
    bytes, err := os.ReadFile("views/" + name)
    if err != nil {
        return err
    }

    tmpl = template.Must(tmpl.Parse(string(bytes)))


    // Query parameters
    queryMap := make(map[string]string)
    for key, values := range r.URL.Query() {
        if len(values) > 0 {
            queryMap[key] = values[0]
        }
    }
    queryAll := make(map[string][]string)
    maps.Copy(queryAll, r.URL.Query())

    // Flash data
    flashData := session.GetFlash(w, r)
    oldData := flashData["old"]
    if oldData == nil {
        oldData = make(map[string]any)
    }
    
    // Authenticated user
    authData := make(map[string]any)
    authData["user"] = middlewares.GetUser(r)

	payload := map[string]any{
        "currentPath": r.URL.Path,
        "query": queryMap, 
        "queryAll": queryAll,
        "flash": flashData,
        "old": oldData,
        "auth": authData,
    }
    maps.Copy(payload, data)

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