package imageserver

import (
    "net/http"
    "crypto/hmac"
    "crypto/sha1"
    "appengine"
    "appengine/blobstore"
    "appengine/image"
    "appengine/datastore"
)

type Applications struct {
    AccessKey string
    SecretKey string
}

type ResultGetUpload struct {
  UploadURL string `json:"upload_url"`
}

type ResultPostUpload struct {
    OriginUrl string `json:"origin_url"`
    OriginSize int `json:"origin_size"`
    ContentType string `json:"content_type"`
    ThumbUrl string `json:"thumb_url"`
}

func Authorize(c appengine.Context, w http.ResponseWriter, r *http.Request) bool {
    values := r.URL.Query()

    accessKey := values["accessKey"]
    if len(accessKey) == 0 {
        err := WriteJsonResponse(w, http.StatusUnauthorized, "Bad accessKey", nil)
        if err != nil {
            serveError(c, w, err)
        }
        return false
    }

    time := values["time"]
    if len(time) == 0 {
        err := WriteJsonResponse(w, http.StatusUnauthorized, "Bad time", nil)
        if err != nil {
            serveError(c, w, err)
        }
        return false
    }

    signature := values["signature"]
    if len(signature) == 0 {
        err := WriteJsonResponse(w, http.StatusUnauthorized, "Bad signature", nil)
        if err != nil {
            serveError(c, w, err)
        }
        return false
    }

    query := datastore.NewQuery("Applications").Filter("AccessKey =", accessKey[0])

    var apps []Applications

    if key, err := query.GetAll(c, &apps); len(key) == 0 || err != nil {
        err := WriteJsonResponse(w, http.StatusInternalServerError, "Application is not registed", nil)
        if err != nil {
            serveError(c, w, err)
        }
        return false
    }

    message := accessKey[0] + "&" + time[0];
    mac := hmac.New(sha1.New, []byte(apps[0].SecretKey))
    mac.Write([]byte(message))
    expectedSignature := string(mac.Sum(nil))

    if signature[0] != expectedSignature {
        err := WriteJsonResponse(w, http.StatusUnauthorized, "Bad signature: " + signature[0], nil)
        if err != nil {
            serveError(c, w, err)
        }
        return false
    }

    return true
}

func HandleUpload(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)

    switch r.Method {
    case "POST","PUT":
        PostUpload(c, w, r)
    default:
        GetUpload(c, w, r)
    }
}

func GetUpload(c appengine.Context, w http.ResponseWriter, r *http.Request) {
    var err error

    if !Authorize(c, w, r) {
        return
    }

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

func PostUpload(c appengine.Context, w http.ResponseWriter, r *http.Request) {
    var err error

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
