package imageserver

import (
  "net/http"
  "imageserver/handler"
)

func init() {
    http.HandleFunc("/init", handler.Init)
    http.HandleFunc("/upload", handler.Upload)
    http.HandleFunc("/image", handler.Image)
}
