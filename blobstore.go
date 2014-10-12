package imageserver

import (
    "net/http"
    "appengine"
    "appengine/blobstore"
)

func HandleBlobstore(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    strBlobKey := r.FormValue("blobKey")

    if strBlobKey == "" {
        http.NotFound(w, r)
        return
    }

    blobKey := appengine.BlobKey(strBlobKey)

    if r.Method == "delete" {
        err := blobstore.Delete(c, blobKey)

        if err != nil {
            serveError(c, w, err)
            return
        }
    } else {
        blobstore.Send(w, blobKey)
    }
}
