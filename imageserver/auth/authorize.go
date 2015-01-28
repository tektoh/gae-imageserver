package auth

import (
	"fmt"
	"time"
	"strconv"
	"net/http"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"appengine"
	"appengine/datastore"
)

type Applications struct {
	AccessKey string `datastore:"accessKey"`
	SecretKey string `datastore:"secretKey,noindex"`
}

func Authorize(c appengine.Context, w http.ResponseWriter, r *http.Request) (bool, string) {
	values := r.URL.Query()
	header := r.Header

	accessKey := values["accessKey"]

	if len(accessKey) == 0 {
		accessKey = header["X-Imageserver-Access-Key"]
	}
	if len(accessKey) == 0 {
		return false, "Bad parameter: accessKey"
	}

	expires := values["expires"]

	if len(expires) == 0 {
		expires = header["X-Imageserver-Expires"]
	}
	if len(expires) == 0 {
		return false, "Bad parameter: expires"
	}

	expiresUnix, err := strconv.ParseInt(expires[0], 10, 64);

	if err != nil {
		c.Errorf("%v", err)
		return false, "Bad parameter: expires"
	}

	nowTime := time.Now().Unix()

	if expiresUnix < nowTime {
		return false, "Has expired"
	}

	signature := values["signature"]
	if len(signature) == 0 {
		signature = header["X-Imageserver-Signature"]
	}
	if len(signature) == 0 {
		return false, "Bad parameter: signature"
	}

	secretKey, err := getSecretKey(c, accessKey[0])
	if err != nil {
		c.Errorf("%v", err)
		return false, "Bad parameter: signature"
	}

	if !checkSignature(c, accessKey[0], expires[0], secretKey, signature[0]) {
		return false, "Bad parameter: signature"
	}

	return true, "OK"
}

func getSecretKey(c appengine.Context, accessKey string) (string, error) {

	query := datastore.NewQuery("Applications").Filter("accessKey =", accessKey).Limit(1)

	var apps []Applications

	if key, err := query.GetAll(c, &apps); len(key) == 0 || err != nil {
		return "", fmt.Errorf("AccessKey %s is not found", accessKey)
	}

	return apps[0].SecretKey, nil
}

func checkSignature(c appengine.Context, accessKey string, expires string, secretKey string, signature string) bool {
	message := accessKey + "&" + expires;

	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(message))
	expectedSignature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	if signature != expectedSignature {
		c.Errorf("Bad signature: client=%s server=%s", signature, expectedSignature)
		return false;
	}

	return true
}
