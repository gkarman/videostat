package response

import (
	"encoding/json"
	"net/http"
	"log/slog"
)

type ErrorBody struct {
	Error string `json:"error"`
}

func JSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		slog.Error(err.Error())
	}
}

func ErrorJSON(w http.ResponseWriter, code int, msg string) {
	JSON(w, code, ErrorBody{Error: msg})
}
