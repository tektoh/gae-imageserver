package handler

import (
    "imageserver/auth"
    "net/http"
    "appengine"
    "appengine/blobstore"
)

type ResultInit struct {
  UploadURL string `json:"upload_url"`
}

func Init(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)

    res, message := auth.Authorize(c, w, r);

    if !res {
        err := WriteJsonResponse(w, http.StatusUnauthorized, message, nil)
        if err != nil {
            ServeError(c, w, err)
        }
        return
    }

    uploadURL, err := blobstore.UploadURL(c, "/upload", nil)

    if err != nil {
        ServeError(c, w, err)
        return
    }

    if err := WriteJsonResponse(w, http.StatusOK, "OK", ResultInit{
        UploadURL: uploadURL.String(),
    }); err != nil {
        ServeError(c, w, err)
        return
    }
}
