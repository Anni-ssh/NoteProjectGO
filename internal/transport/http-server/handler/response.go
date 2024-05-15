package handler

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"strings"
)

type Response struct {
	ErrMsg string `json:"Error,omitempty"`
}

func NewErrResponse(w http.ResponseWriter, statusCode int, errMessage string) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(Response{ErrMsg: errMessage}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("handler.NewResponse failed Encode JSON to response", slog.Any("error", err))
	}
}

func ErrValidator(errs validator.ValidationErrors) Response {
	var errList []string
	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errList = append(errList, fmt.Sprintf("field %s is required field", err.Field()))
		case "username":
			errList = append(errList, fmt.Sprintf("field %s is not a valid username", err.Field()))
		case "password":
			errList = append(errList, fmt.Sprintf("ield %s is not a valid username", err.Field()))
		default:
			errList = append(errList, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}
	return Response{ErrMsg: strings.Join(errList, ", ")}
}
