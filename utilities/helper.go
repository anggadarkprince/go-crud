package utilities

import (
	"slices"
	"fmt"
	"reflect"
	"strings"
	"text/template"
)

var TemplateFuncs = template.FuncMap{
    "add": func(a, b int) int { return a + b },
    "toUpper": strings.ToUpper,
    "hasPrefix": strings.HasPrefix,
    "contains": func(arr []string, value string) bool {
        return slices.Contains(arr, value)
    },
    "containsByField": func(list any, fieldName string, value any) bool {
        v := reflect.ValueOf(list)

        // Check list is a slice
        if v.Kind() != reflect.Slice {
            return false
        }

        for i := 0; i < v.Len(); i++ {
            elem := v.Index(i)

            // Check elem adalah struct
            if elem.Kind() == reflect.Ptr {
                elem = elem.Elem()
            }
            if elem.Kind() != reflect.Struct {
                continue
            }

            // Get field by name
            field := elem.FieldByName(fieldName)
            if !field.IsValid() {
                continue
            }

            // Check field with value then convert to string
            if fmt.Sprint(field.Interface()) == fmt.Sprint(value) {
                return true
            }
        }

        return false
    },
    // "formatDate": FormatDate,
}