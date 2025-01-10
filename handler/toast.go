package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	INFO    = "info"
	SUCCESS = "success"
	WARNING = "warning"
	DANGER  = "danger"
)

type Toast struct {
	Level   string `json:"level"`
	Message string `json:"message"`
}

func New(level string, message string) Toast {
	return Toast{level, message}
}

func Info(message string) Toast {
	return New(INFO, message)
}

func Success(w http.ResponseWriter, message string) {
	New(SUCCESS, message).SetHXTriggerHeader(w)
}

func Warning(message string) Toast {
	return New(WARNING, message)
}

func Danger(message string) Toast {
	return New(DANGER, message)
}

func (t Toast) Error() string {
	return fmt.Sprintf("%s: %s", t.Level, t.Message)
}

func (t Toast) jsonify() (string, error) {
	t.Message = t.Error()
	eventMap := map[string]Toast{}
	eventMap["makeToast"] = t
	jsonData, err := json.Marshal(eventMap)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

func (t Toast) SetHXTriggerHeader(w http.ResponseWriter) {
	jsonData, _ := t.jsonify()
	w.Header().Set("HX-Trigger", jsonData)
}
