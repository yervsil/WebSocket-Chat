package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/go-playground/validator.v9"
)

type Response struct {
	Message interface{} `json:"message"`
}

func SendJson(w http.ResponseWriter, message interface{}, status int) {
	w.Header().Add("Content-Type", "application/json")
	resp, err := json.Marshal(Response{Message: message})
	if err != nil{
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}

	w.WriteHeader(status)
	w.Write(resp)
}

func InvalidFields(v validator.ValidationErrors) []string{
	var errorMessages []string
	for _, fieldError := range v {
		errorMessages = append(errorMessages, fmt.Sprintf("Field '%s' failed validation with tag '%s'", fieldError.Field(), fieldError.Tag()))
	}

	return errorMessages
}