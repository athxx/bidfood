package hdl

import (
	"encoding/json"
	"errors"
	"net/http"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Ok(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(Response{
		Code: 0,
		Msg:  "",
		Data: data,
	})
}

func Err(w http.ResponseWriter, status int, message string, err error) {
	if err == nil {
		err = errors.New("Unknown Error")
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Response{
		Code: status,
		Msg:  message,
		Data: err.Error(),
	})
}
