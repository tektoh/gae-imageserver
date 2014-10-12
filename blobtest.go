package blobtest

import (
        "html/template"
        "io"
        "net/http"

        "appengine"
        "appengine/blobstore"
)

func serveError(c appengine.Context, w http.ResponseWriter, err error) {
        w.WriteHeader(http.StatusInternalServerError)
        w.Header().Set("Content-Type", "text/plain")
        io.WriteString(w, "Internal Server Error")
        c.Errorf("%v", err)
}

var rootTemplate = template.Must(template.New("root").Parse(rootTemplateHTML))

const rootTemplateHTML = `
<html><body>
<form action="{{.}}" method="POST" enctype="multipart/form-data">
Upload File: <input type="file" name="file"><br>
<input type="submit" name="submit" value="Submit">
</form></body></html>
`

func handleRoot(w http.ResponseWriter, r *http.Request) {
        c := appengine.NewContext(r)
        uploadURL, err := blobstore.UploadURL(c, "/upload", nil)
        if err != nil {
                serveError(c, w, err)
                return
        }
        w.Header().Set("Content-Type", "text/html")
        err = rootTemplate.Execute(w, uploadURL)
        if err != nil {
                c.Errorf("%v", err)
        }
}

func handleServe(w http.ResponseWriter, r *http.Request) {
        blobstore.Send(w, appengine.BlobKey(r.FormValue("blobKey")))
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
        c := appengine.NewContext(r)
        blobs, _, err := blobstore.ParseUpload(r)
        if err != nil {
                serveError(c, w, err)
                return
        }
        file := blobs["file"]
        if len(file) == 0 {
                c.Errorf("no file uploaded")
                http.Redirect(w, r, "/", http.StatusFound)
                return
        }
        http.Redirect(w, r, "/serve/?blobKey="+string(file[0].BlobKey), http.StatusFound)
}

func init() {
        http.HandleFunc("/", handleRoot)
        http.HandleFunc("/serve/", handleServe)
        http.HandleFunc("/upload", handleUpload)
}
