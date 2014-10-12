package imageserver

import (
  "io"
  "net/http"
  "appengine"
)

func serveError(c appengine.Context, w http.ResponseWriter, err error) {
  w.WriteHeader(http.StatusInternalServerError)
  w.Header().Set("Content-Type", "text/plain")
  io.WriteString(w, "Internal Server Error")
  c.Errorf("%v", err)
}

func init() {
  http.HandleFunc("/init", HandleInit)
  http.HandleFunc("/upload", HandleUpload)
  http.HandleFunc("/blobstore", HandleBlobstore)
}

