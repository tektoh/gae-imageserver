package handler

import (
    "net/http"
    "encoding/json"
)

type JsonResponse struct {
  Status  int         `json:"status"`
  Message string      `json:"message"`
  Result  interface{} `json:"result"`
}

func WriteJsonResponse(w http.ResponseWriter, status int, message string, result interface{}) error {
    res := JsonResponse{
        Status: status,
        Message: message,
        Result: result,
    }

    json, err := json.Marshal(res)

    if err != nil {
      return err
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(res.Status)

    w.Write(json)

    return nil
}
