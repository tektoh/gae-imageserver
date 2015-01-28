package info

import (
	"appengine"
	"appengine/datastore"
)

type ImageInfo struct {
	OriginUrl string `datastore:"originUrl"`
	OriginSize int `datastore:"originSize"`
	Filename string `datastore:"filename"`
	ContentType string `datastore:"contentType"`
	ThumbUrl string `datastore:"thumbUrl"`

	context appengine.Context `datastore:"-"`
	key *datastore.Key `datastore:"-"`
}

func NewImageInfo(c appengine.Context, blobKey string) *ImageInfo {
	info := new(ImageInfo)
	info.context = c
	info.key = datastore.NewKey(c, "ImageInfo", blobKey, 0, nil)
	return info
}

func (info *ImageInfo) Save(originUrl string, originSize int, filename string,
	contentType string, thumbUrl string) error {

	info.OriginUrl   = originUrl
	info.OriginSize  = originSize
	info.Filename    = filename
	info.ContentType = contentType
	info.ThumbUrl    = thumbUrl

	if _, err := datastore.Put(info.context, info.key, info); err != nil {
		info.context.Errorf("Error: datastore.Put(%v, %v, %v) failed", info.context, info.key, info)
		return err
	}

	return nil
}

func (info *ImageInfo) Delete() error {
	if err := datastore.Delete(info.context, info.key); err != nil {
		info.context.Errorf("Error: datastore.Delete(%v, %v) failed", info.context, info.key)
		return err
	}

	return nil
}
