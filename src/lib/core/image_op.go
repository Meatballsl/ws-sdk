package core

import (
	"errors"
	"github.com/Meatballsl/ws-sdk/src/lib/utility"
	"net/http"
	"net/url"
	"strconv"
)

type ImageOp struct {
	auth        *utility.Auth
	config      *Config
	httpManager *utility.HttpManager
}

func NewImageOp(auth *utility.Auth, config *Config, client *http.Client) (bm *ImageOp) {
	if nil == auth {
		panic("auth is nil")
	}
	if nil == config {
		config = NewDefaultConfig()
	}
	return &ImageOp{auth, config, utility.NewHttpManager(client)}
}

// 图片鉴定
// https://wcs.chinanetcenter.com/document/API/Image-op/imageDetect
// type: porn-鉴黄,terror-暴恐,political-政治人物识别
func (this *ImageOp) ImageDetect(image string, _type string, bucket string) (response *http.Response, err error) {
	if 0 == len(image) {
		err = errors.New("image is empty")
		return
	}
	if 0 == len(_type) {
		err = errors.New("type is empty")
		return
	}
	if "porn" != _type &&
		"terror" != _type &&
		"political" != _type {
		err = errors.New("type is invalid")
		return
	}
	if 0 == len(bucket) {
		err = errors.New("bucket is empty")
		return
	}

	query := make(url.Values)
	query.Add("image", utility.UrlSafeEncodeString(image))
	query.Add("type", _type)
	query.Add("bucket", bucket)
	return InnerFOps(this.auth, this.httpManager.GetClient(), this.config.GetManageUrlPrefix()+"/imageDetect", query.Encode())
}

// 图片持久化处理
// https://wcs.chinanetcenter.com/document/API/Image-op/imagePersistentOp
func (this *ImageOp) ImagePersistentOp(bucket string, key string, fops string, notify_url string, force int, separate int) (response *http.Response, err error) {
	if 0 == len(bucket) {
		err = errors.New("bucket is empty")
		return
	}
	if 0 == len(key) {
		err = errors.New("key is empty")
		return
	}

	query := make(url.Values)
	query.Add("bucket", utility.UrlSafeEncodeString(bucket))
	query.Add("key", utility.UrlSafeEncodeString(key))
	query.Add("fops", utility.UrlSafeEncodeString(fops))
	if len(notify_url) > 0 {
		query.Add("notifyURL", utility.UrlSafeEncodeString(notify_url))
	}
	if 0 == force || 1 == force {
		query.Add("force", strconv.Itoa(force))
	}
	if 0 == separate || 1 == separate {
		query.Add("separate", strconv.Itoa(separate))
	}
	return FOps(this.auth, this.config, this.httpManager.GetClient(), query.Encode())
}
