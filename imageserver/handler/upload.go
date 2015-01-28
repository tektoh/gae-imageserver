package handler

import (
    "imageserver/info"
    "net/http"
    "strings"
    "appengine"
    "appengine/blobstore"
    "appengine/image"
    "github.com/famz/RFC2047"
)

type ResultUpload struct {
    OriginUrl string `json:"origin_url"`
    OriginSize int `json:"origin_size"`
    Filename string `json:"filename"`
    ContentType string `json:"content_type"`
    ThumbUrl string `json:"thumb_url"`
}

func Upload(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    var err error

    if r.Method != "POST" {
        err = WriteJsonResponse(w, http.StatusBadRequest, "Bad method", nil)
        if err != nil {
            ServeError(c, w, err)
        }
        return
    }

    blobs, _, err := blobstore.ParseUpload(r)

    if err != nil {
        ServeError(c, w, err)
        return
    }

    file := blobs["file"]

    if len(file) == 0 {
        err := WriteJsonResponse(w, http.StatusBadRequest, "No file uploaded", nil)
        if err != nil {
            ServeError(c, w, err)
        }
        return
    }

    if file[0].ContentType != "image/jpeg" && file[0].ContentType != "image/png" {
        err = blobstore.Delete(c, file[0].BlobKey)
        err = WriteJsonResponse(w, http.StatusBadRequest, "Bad mimetype", nil)
        if err != nil {
            ServeError(c, w, err)
        }
        return
    }

    blobKey       := string(file[0].BlobKey)
    originUrl     := "https://"+r.Host+"/image?blobKey="+blobKey
    originSize    := int(file[0].Size)
    filename      := RFC2047.Decode(file[0].Filename)
    contentType   := file[0].ContentType
    thumbUrl, err := image.ServingURL(c, file[0].BlobKey, nil)

    if err != nil {
        ServeError(c, w, err)
        return
    }

    strThumbUrl := strings.Replace(thumbUrl.String(), "http://", "https://", 1)

    imageInfo := info.NewImageInfo(c, blobKey)

    if err = imageInfo.Save(originUrl, originSize, filename, contentType, strThumbUrl); err != nil {

        err = blobstore.Delete(c, file[0].BlobKey)

        err = WriteJsonResponse(w, http.StatusBadRequest, "Save failed", nil)

        if err != nil {
            ServeError(c, w, err)
        }

        return
    }

    err = WriteJsonResponse(w, http.StatusOK, "OK", ResultUpload{
        OriginUrl: originUrl,
        OriginSize: originSize,
        Filename: filename,
        ContentType: contentType,
        ThumbUrl: strThumbUrl,
    })

    if err != nil {
        ServeError(c, w, err)
        return
    }
}
