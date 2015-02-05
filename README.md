# gae-imageserver

* GoogleAppEngine上で動作する画像アップロードサーバーです。
* オリジナル画像の保存、サムネイルの生成を行います。

## GoogleAppEngineについて

* Google製のPaaSです。
    * https://cloud.google.com/appengine/
* この画像サーバーはAppEngineの`BlobStore`、`Images Service`を利用しています。
    * BlobStore: https://cloud.google.com/appengine/docs/go/blobstore/
        * CloudStorageへのファイルアップロード機能を簡単に実装できます
        * ファイルへのアクセスにはアップロード時に発行される`BlobKey`を使用します。
    * Images: https://cloud.google.com/appengine/docs/go/images/
        * `BlobStore`にアップロードされた画像ファイルのサムネイルを作成します。
        * 0～1600pxの間で自由なサイズのサムネイルを動的に取得できます。

## API

### GET /init
* BlobStoreに画像アップロードするためのURLを取得します。
* 署名によるクライアント認証を行います。

#### クエリパラメータ
* accessKey: アクセスキー
* signature: 署名
* expires: 署名の有効期限

```
/upload?accessKey=eKWuZo6wDXJyuciVQkHMu6Bd5eZLLukh&expires=1414857216&signature=PtemPdruhb3d%2Bq5AoiAQRIBDR9oO5BzukmQ4D5GcINo%3D
```

#### リクエストヘッダー
署名の情報をリクエストヘッダーに含めることもできます。

* X-Imageserver-Access-Key: アクセスキー
* X-Imageserver-Signature: 署名
* X-Imageserver-Expires: 署名の有効期限


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
* AppEngineにより長いURLが発行されるため、認証は行っていません。

#### リクエストデータ

multipart/form-data 形式。

* file: 画像ファイル

#### レスポンスデータ
* status: ステータスコード（HTTPレスポンスと同じ)
* message: エラーメッセージ
* result
    * origin_url: オリジナル画像のURL
    * origin_size: アップロードした画像のファイルサイズ
    * filename: アップロードした画像のファイル名
    * content_type: アップロードした画像のMIMEタイプ
    * thumb_url: サムネイルのURL

```
{
  "status": 200,
  "message": "OK",
  "result": {
    "origin_url": "https://example.appspot.com/image?blobKey=AMIfv96V2rPci4U9WoK-HmfAaFbcdSTPwcQot-yfq0Op9Auz3QSkNp4mzvQXGrl8-kgRkKTDtPYVChvq9ALIrhebrC0vAjyiym7LSXIQsk2ipcT00fur7aa1Y-bUALqI78FDqG-KkqskKVxzkjmXkw4iwyE_Kje15ygtf9BhfJXWb3jxe7XKy6A",
    "origin_size": 583122,
    "filename": "G0017228.JPG",
    "content_type": "image/jpeg",
    "thumb_url": "https://lh5.ggpht.com/hmCTZkntu1SO929RgUZJbrJd5B1QgSe57Cw3fTD4u-MA10SexhV4CAlz1r6bMROkiAGhl_nqu940fwpomEArDeG5HB0SWQ"
  }
}
```

### GET /image

* BlobStoreからオリジナル画像を取り出します。
* AppEngineにより長いオブジェクトキーが発行されるため、認証は行っていません。

#### クエリパラメータ
* blobKey: BlobStoreのKey

```
/image?blobKey=AMIfv97ZUCKTJf1-AdenPtrCbXCJkfyxzw0LVjJxY-4KLEWGHu67aZRZKwkW3Itkda9esI3Wt1jKvJ2Usr0E5h4NfFCYV7-J5VJJC_deaJFqLlfPAQYWqatZaWtcM_JLOq6drJ6__8CTTQAb5gRTyUZJYA0ZeSa2XDGyR98UfswpNWhnVX_m4bo
```

### DELETE /image

* BlobStoreからオリジナル画像を削除します。
* 署名によるクライアント認証を行います。

#### クエリパラメータ
* blobKey: BlobStoreのKey

```
/image?blobKey=AMIfv97ZUCKTJf1-AdenPtrCbXCJkfyxzw0LVjJxY-4KLEWGHu67aZRZKwkW3Itkda9esI3Wt1jKvJ2Usr0E5h4NfFCYV7-J5VJJC_deaJFqLlfPAQYWqatZaWtcM_JLOq6drJ6__8CTTQAb5gRTyUZJYA0ZeSa2XDGyR98UfswpNWhnVX_m4bo
```

### サムネイルのURL
0〜1600ピクセルの間で画像サイズを指定できます。GoogleAppEngineのドキュメントを参照してください。
https://cloud.google.com/appengine/docs/go/images/#Go_Serving_and_re-sizing_images_from_the_Blobstore

URLは以下の様になります。オリジナル画像をアップロードしたAppEngineとは別インスタンス（おそらくPicasa）になります。
```
http://your_app_id.appspot.com/randomStringImageId
```

URLの末尾に`=sxxx`もしくは`=sxxx-c`を加えることでサムネイル画像のサイズを自由に指定できます。
```
// Resize the image to 32 pixels (aspect-ratio preserved)
http://your_app_id.appspot.com/randomStringImageId=s32

// Crop the image to 32 pixels
http://your_app_id.appspot.com/randomStringImageId=s32-c
```
