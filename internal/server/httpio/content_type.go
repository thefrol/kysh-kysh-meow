package httpio

import "net/http"

const (
	TypeTextPlain       = "text/plain"
	TypeApplicationJSON = "application/json"
	TypeTextHTML        = "text/html"
)

func SetContentType(w http.ResponseWriter, contentType string) {
	w.Header().Set("Content-Type", contentType)
}
