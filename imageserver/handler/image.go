package handler

import (
    "imageserver/auth"
    "imageserver/info"
    "net/http"
    "appengine"
    "appengine/blobstore"
)

func Image(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    strBlobKey := r.FormValue("blobKey")

    if strBlobKey == "" {
        http.NotFound(w, r)
        return
    }

    blobKey := appengine.BlobKey(strBlobKey)

    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET, DELETE, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "X-Imageserver-Access-Key, X-Imageserver-Expires, X-Imageserver-Signature")

    if r.Method == "DELETE" {

        if res, message := auth.Authorize(c, w, r); !res {
            err := WriteJsonResponse(w, http.StatusUnauthorized, message, nil)
            if err != nil {
                ServeError(c, w, err)
            }
            return
        }

        if err := blobstore.Delete(c, blobKey); err != nil {
            ServeError(c, w, err)
            return
        }

        imageInfo := info.NewImageInfo(c, strBlobKey)

        if err := imageInfo.Delete(); err != nil {
            c.Errorf("%v", err)
        }

        if err := WriteJsonResponse(w, http.StatusOK, "OK", nil); err != nil {
            ServeError(c, w, err)
            return
        }

    } else if r.Method == "OPTIONS" {

        if err := WriteJsonResponse(w, http.StatusOK, "OK", nil); err != nil {
            ServeError(c, w, err)
            return
        }

    } else {
        blobstore.Send(w, blobKey)
    }
}
