# gae-imageserver

GoogleAppEngine上で動作する画像アップロードサーバーです。
オリジナル画像の保存、サムネイルの生成を行います。

## API

### GET /upload
BlobStoreに画像アップロードするためのURLを取得します。

#### クエリパラメータ
* accessKey: アクセスキー
* signature: 署名
* expires: 署名の有効期限

```
/upload?accessKey=eKWuZo6wDXJyuciVQkHMu6Bd5eZLLukh&expires=1414857216&signature=PtemPdruhb3d%2Bq5AoiAQRIBDR9oO5BzukmQ4D5GcINo%3D
```

#### レスポンスデータ
* status: ステータスコード（HTTPレスポンスと同じ)
* message: エラーメッセージ
* result
    * upload_url: アップロード用URL

```
{
  "status": 200,
  "message": "OK",
  "result": {
    "upload_url": "http://example.appspot.com/_ah/upload/?accessKey=eKWuZo6wDXJyuciVQkHMu6Bd5eZLLukh&expires=1414857216&signature=PtemPdruhb3d%2Bq5AoiAQRIBDR9oO5BzukmQ4D5GcINo%3D/AMmfu6bjD2nDqtT1nX8nws4-ImEPQluohIFjqSAnsEQOrglUu-Ma63R_gH4Hs0CZ5DqsGe6typbyc5ZmHhE2bzcHBcukhbBSAWTrMtz0PFH0vEPd7To955VeZWYB-FBt50umy6yQthFg/ALBNUaYAAAAAVFT8OVElO0ebF3GPOVWRDbhAcf7Fc85o/"
  }
}
```

#### 署名の作り方

1. accessKeyとexpiresの間に"&"を入れた文字列(message)を作ります。
2. secretKeyをキーにし、messageのHMAC-SHA256のハッシュ値を求めます。
3. ハッシュ値をさらにbase64エンコードします。
4. さらにURLエンコードしてクエリパラメータに設定します。

```
message = accessKey + "&" + expires
urlencode(base64(hmac-sha256(message, secretKey)))
```

#### 署名の有効期限

サーバー側の時刻と比較し、期限が切れていた場合はエラーを返します。

#### アクセスキーとシークレットキー

* index.yaml を参考にしてDataStoreに`accessKey`と`secretKey`をDeveloper Consoleから直接入力してエンティティを作成してください。
* `accessKey`にはインデックスを設定してください。`secretKey`には不要です。

### POST /_ah/upload
* `GET /upload` で取得できたURLに画像ファイルをPOSTします。

#### レスポンスデータ
* status: ステータスコード（HTTPレスポンスと同じ)
* message: エラーメッセージ
* result
    * origin_url: オリジナル画像のURL
    * origin_size: アップロードした画像のファイルサイズ
    * content_type: アップロードした画像のMIMEタイプ
    * thumb_url: サムネイルのURL

```
{
  "status": 200,
  "message": "OK",
  "result": {
    "origin_url": "http://example.appspot.com/blobstore?blobKey=AMIfv97ZUCKTJf1-AdenPtrCbXCJkfyxzw0LVjJxY-4KLEWGHu67aZRZKwkW3Itkda9esI3Wt1jKvJ2Usr0E5h4NfFCYV7-J5VJJC_deaJFqLlfPAQYWqatZaWtcM_JLOq6drJ6__8CTTQAb5gRTyUZJYA0ZeSa2XDGyR98UfswpNWhnVX_m4bo",
    "origin_size": 491699,
    "content_type": "image/png",
    "thumb_url": "http://lh3.ggpht.com/yGorNPKra-Jmcc-x5i9-ncs83dUlh8YV8j0EWjYb8NM3Vt-vjEzwvCL1Uh6mpo-FqPDHTOSaEoa_OBtEZO30-j69xxrh"
  }
}
```

### GET /blobstore

BlobStoreからオリジナル画像を取り出します。

### DELETE /blobstore

BlobStoreからオリジナル画像を削除します。

### サムネイルのURL
0〜1600ピクセルの間で画像サイズを指定できます。GoogleAppEngineのドキュメントを参照してください。    
https://cloud.google.com/appengine/docs/go/images/#Go_Serving_and_re-sizing_images_from_the_Blobstore
