package utilities

import "net/http"

func Redirect(w http.ResponseWriter, path string) {
	http.Redirect(w, &http.Request{}, path, http.StatusFound)
}

func QueryParam(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}

func FormValue(r *http.Request, key string) string {
	_ = r.ParseForm()
	return r.FormValue(key)
}