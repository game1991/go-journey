package sha

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"io"
	"os"
	"sync"
)

const (
	minSize     = 4 * 1024
	defaultSize = 4 * 1024 * 1024
	maxSize     = 100 * 1024 * 1024
)

var b64 *base64.Encoding

func init() {
	b64 = base64.StdEncoding.WithPadding(base64.NoPadding)
}

//SHA256 计算字符数组的sha256 信息摘要
func SHA256(bts []byte) string {
	sha := sha256.New()
	sha.Write(bts)
	return hex.EncodeToString(sha.Sum(nil))
}

//SHA256BASE64 返回base64编码
func SHA256BASE64(bts []byte) string {
	sha := sha256.New()
	sha.Write(bts)
	return b64.EncodeToString(sha.Sum(nil))
}

//SHA256Stream 计算文件流信息摘要
func SHA256Stream(reader io.Reader, size ...int) ([]byte, error) {
	var length int = 0
	if len(size) == 0 {
		length = defaultSize
	} else {
		length = size[0]
		if length < minSize {
			length = minSize
		} else if length > maxSize {
			length = maxSize
		}
	}

	buff := make([]byte, length)
	return sha256Stream(reader, buff)
}

var buffPool *sync.Pool

func init() {
	buffPool = &sync.Pool{
		New: func() interface{} {
			return make([]byte, 1024*1024) //1m缓存
		},
	}
}

//SHA256StreamPool 高并发文件流计算时，使用byte pool池化技术
func SHA256StreamPool(reader io.Reader) ([]byte, error) {
	buff := buffPool.Get().([]byte)
	defer buffPool.Put(buff)
	return sha256Stream(reader, buff)
}

func sha256Stream(reader io.Reader, buff []byte) ([]byte, error) {
	hash := sha256.New()
	_, err := io.CopyBuffer(hash, reader, buff)
	if err != nil {
		return nil, err
	}
	return hash.Sum(nil), nil
}

//SHA256File 计算文件sha256
func SHA256File(path string) ([]byte, error) {
	fs, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return SHA256StreamPool(fs)
}

//SHA256FileBASE64 计算文件sha256 返回bash64
func SHA256FileBASE64(path string) (string, error) {
	hs, err := SHA256File(path)
	if err != nil {
		return "", err
	}
	return b64.EncodeToString(hs), nil
}

//SHA256FileHEX 计算文件sha256 返回bash64
func SHA256FileHEX(path string) (string, error) {
	hs, err := SHA256File(path)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hs), nil
}
