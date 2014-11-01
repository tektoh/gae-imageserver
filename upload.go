package imageserver

import (
    "net/http"
    "appengine"
    "appengine/blobstore"
    "appengine/image"
)

type ResultGetUpload struct {
  UploadURL string `json:"upload_url"`
}

type ResultPostUpload struct {
    OriginUrl string `json:"origin_url"`
    OriginSize int `json:"origin_size"`
    ContentType string `json:"content_type"`
    ThumbUrl string `json:"thumb_url"`
}

func HandleUpload(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case "POST","PUT":
        PostUpload(w, r)
    default:
        GetUpload(w, r)
    }
}

func GetUpload(w http.ResponseWriter, r *http.Request) {
    var err error
    c := appengine.NewContext(r)

    uploadURL, err := blobstore.UploadURL(c, "/upload", nil)

    if err != nil {
        serveError(c, w, err)
        return
    }

    err = WriteJsonResponse(w, http.StatusOK, "OK", ResultGetUpload{
        UploadURL: uploadURL.String(),
    })

    if err != nil {
        serveError(c, w, err)
        return
    }
}

func PostUpload(w http.ResponseWriter, r *http.Request) {
    var err error
    c := appengine.NewContext(r)

    blobs, _, err := blobstore.ParseUpload(r)

    if err != nil {
        serveError(c, w, err)
        return
    }

    file := blobs["file"]

    if len(file) == 0 {
        err := WriteJsonResponse(w, http.StatusBadRequest, "No file uploaded", nil)
        if err != nil {
            serveError(c, w, err)
        }
        return
    }

    if file[0].ContentType != "image/jpeg" && file[0].ContentType != "image/png" {
        err = blobstore.Delete(c, file[0].BlobKey)
        err = WriteJsonResponse(w, http.StatusBadRequest, "Bad mimetype", nil)
        if err != nil {
            serveError(c, w, err)
        }
        return
    }

    originUrl := "http://"+r.Host+"/blobstore?blobKey="+string(file[0].BlobKey)
    thumbUrl, err := image.ServingURL(c, file[0].BlobKey, nil)

    if err != nil {
        serveError(c, w, err)
        return
    }

    err = WriteJsonResponse(w, http.StatusOK, "OK", ResultPostUpload{
        OriginUrl: originUrl,
        OriginSize: int(file[0].Size),
        ContentType: file[0].ContentType,
        ThumbUrl: thumbUrl.String(),
    })

    if err != nil {
        serveError(c, w, err)
        return
    }
}
