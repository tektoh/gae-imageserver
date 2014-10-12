package imageserver

import (
    "net/http"
    "appengine"
    "appengine/blobstore"
)

type ResultInit struct {
  UploadURL string `json:"upload_url"`
}

func HandleInit(w http.ResponseWriter, r *http.Request) {
    var err error
    c := appengine.NewContext(r)

    uploadURL, err := blobstore.UploadURL(c, "/upload", nil)

    if err != nil {
        serveError(c, w, err)
        return
    }

    err = WriteJsonResponse(w, http.StatusOK, "OK", ResultInit{
        UploadURL: uploadURL.String(),
    })

    if err != nil {
        serveError(c, w, err)
        return
    }
}
