package handler

import (
    "io"
    "net/http"
    "appengine"
)

func ServeError(c appengine.Context, w http.ResponseWriter, err error) {
    w.WriteHeader(http.StatusInternalServerError)
    w.Header().Set("Content-Type", "text/plain")
    io.WriteString(w, "Internal Server Error")
    c.Errorf("%v", err)
}
