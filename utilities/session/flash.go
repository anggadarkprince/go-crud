package session

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/anggadarkprince/crud-employee-go/configs"
)

type FlashData map[string]any

// SetFlash sets multiple flash messages as key-value pairs
func SetFlash(w http.ResponseWriter, data FlashData) {
    jsonData, err := json.Marshal(data)
    if err != nil {
        return
    }
	// Base64 encode to avoid special characters in cookie
    encoded := base64.StdEncoding.EncodeToString(jsonData)
    
    cookie := &http.Cookie{
        Name: configs.Get().Session.StoreName,
        Value: encoded,
        Path: configs.Get().Session.Path,
        HttpOnly: true,
        MaxAge: 60, // 60 seconds - short lived
        SameSite: http.SameSiteLaxMode,
    }
    http.SetCookie(w, cookie)
}

// SetFlashMessage is a convenience method to set a single typed message
func Flash(w http.ResponseWriter, messageType, message string) {
	flashMessage := FlashData{
		"alert": map[string]string{
			"type": messageType,
			"message": message,
		},
	}
	SetFlash(w, flashMessage)
}

// GetFlash retrieves and deletes all flash messages
func GetFlash(w http.ResponseWriter, r *http.Request) FlashData {
    cookie, err := r.Cookie(configs.Get().Session.StoreName)
    if err != nil {
        return nil
    }
    
    // Delete the cookie immediately
    http.SetCookie(w, &http.Cookie{
        Name: configs.Get().Session.StoreName,
        Value: "",
        Path: "/",
        MaxAge: -1,
    })

	// Base64 decode first
    decoded, err := base64.StdEncoding.DecodeString(cookie.Value)
    if err != nil {
        return nil
    }
    
    var data FlashData
    err = json.Unmarshal([]byte(decoded), &data)
    if err != nil {
        return nil
    }
    
    return data
}

// Flash helper methods for common patterns
func FlashError(w http.ResponseWriter, message string) {
    Flash(w, "error", message)
}

func FlashSuccess(w http.ResponseWriter, message string) {
    Flash(w, "success", message)
}

func FlashWarning(w http.ResponseWriter, message string) {
    Flash(w, "warning", message)
}

func FlashInfo(w http.ResponseWriter, message string) {
    Flash(w, "info", message)
}

func FlashDanger(w http.ResponseWriter, message string) {
    Flash(w, "danger", message)
}

// FlashWithInput stores form input data along with a message
func FlashWithInput(w http.ResponseWriter, messageType, message string, input map[string]any) {
    data := FlashData{
		"alert": map[string]string{
			"type": messageType,
			"message": message,
		},
        "old": input,
    }
    SetFlash(w, data)
}

func ParseFormInput(r *http.Request) map[string]any {
    if err := r.ParseForm(); err != nil {
        return nil
    }
    
    input := make(map[string]any)
    
    for key, values := range r.Form {
        if len(values) == 1 {
            // Single value
            input[key] = values[0]
        } else if len(values) > 1 {
            // Multiple values (array)
            input[key] = values
        }
    }
    
    return input
}

func ParseMultipartFormInput(r *http.Request, maxMemory int64) map[string]interface{} {
    if maxMemory == 0 {
        maxMemory = 32 << 20 // 32 MB default
    }
    
    if err := r.ParseMultipartForm(maxMemory); err != nil {
        return nil
    }
    
    input := make(map[string]interface{})
    
    // Get regular form values
    for key, values := range r.MultipartForm.Value {
        if len(values) == 1 {
            input[key] = values[0]
        } else if len(values) > 1 {
            input[key] = values
        }
    }
    
    // Get file names (not the actual files, just names for repopulation)
    for key, files := range r.MultipartForm.File {
        if len(files) == 1 {
            input[key+"_filename"] = files[0].Filename
        } else if len(files) > 1 {
            filenames := make([]string, len(files))
            for i, file := range files {
                filenames[i] = file.Filename
            }
            input[key+"_filenames"] = filenames
        }
    }
    
    return input
}

// GetOldInput retrieves old input from flash data
func GetOldInput(flash FlashData, key string, defaultValue ...string) string {
    if flash == nil {
        if len(defaultValue) > 0 {
            return defaultValue[0]
        }
        return ""
    }
    
    old, ok := flash["old"]
    if !ok {
        if len(defaultValue) > 0 {
            return defaultValue[0]
        }
        return ""
    }
    
    oldMap, ok := old.(map[string]interface{})
    if !ok {
        if len(defaultValue) > 0 {
            return defaultValue[0]
        }
        return ""
    }
    
    value, ok := oldMap[key]
    if !ok {
        if len(defaultValue) > 0 {
            return defaultValue[0]
        }
        return ""
    }
    
    // Handle string value
    if str, ok := value.(string); ok {
        return str
    }
    
    // Handle array value (return first item)
    if arr, ok := value.([]interface{}); ok && len(arr) > 0 {
        if str, ok := arr[0].(string); ok {
            return str
        }
    }
    
    if len(defaultValue) > 0 {
        return defaultValue[0]
    }
    return ""
}

// GetOldInputArray retrieves old input array from flash data
func GetOldInputArray(flash FlashData, key string) []string {
    if flash == nil {
        return []string{}
    }
    
    old, ok := flash["old"]
    if !ok {
        return []string{}
    }
    
    oldMap, ok := old.(map[string]interface{})
    if !ok {
        return []string{}
    }
    
    value, ok := oldMap[key]
    if !ok {
        return []string{}
    }
    
    // Handle string value (convert to single-item array)
    if str, ok := value.(string); ok {
        return []string{str}
    }
    
    // Handle array value
    if arr, ok := value.([]interface{}); ok {
        result := make([]string, 0, len(arr))
        for _, item := range arr {
            if str, ok := item.(string); ok {
                result = append(result, str)
            }
        }
        return result
    }
    
    return []string{}
}