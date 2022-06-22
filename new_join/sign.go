package signature

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// Sign 加签函数
type Sign func(content, privateKey string) (sign string, err error)

// Verify 验签函数
type Verify func(content, sign, pubKey string) (err error)

// NewSigner 初始化Signer，默认RSA2签名
func NewSigner(s Sign) *Signer {
	if s == nil {
		s = rsa2Sign
	}
	return &Signer{
		S: s,
	}
}

var ErrPemDecode = errors.New("pem.Decode failed") // pem解析失败

func rsa2Sign(content, privateKey string) (sign string, err error) {
	// 1、将密钥解析成密钥实例
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		err = ErrPemDecode
		return
	}
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return
	}

	// 2、生成签名
	hash := sha256.New()
	_, err = hash.Write([]byte(content))
	if err != nil {
		return
	}
	signature, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, hash.Sum(nil))
	if err != nil {
		return
	}

	// 3、签名base64编码
	sign = base64.StdEncoding.EncodeToString(signature)
	return
}

// Signer
type Signer struct {
	S Sign
}

// Sign 签名
func (s *Signer) Sign(content, privateKey string) (sign string, err error) {
	return s.S(content, privateKey)
}

// NewVerifier 初始化Verifier，默认RSA2验签
func NewVerifier(v Verify) *Verifier {
	if v == nil {
		v = rsa2Verify
	}
	return &Verifier{
		V: v,
	}
}

func rsa2Verify(content, sign, pubKey string) (err error) {
	// 1、签名base64解码
	signature, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return
	}

	// 2、密钥解析成公钥实例
	block, _ := pem.Decode([]byte(pubKey))
	if block == nil {
		err = ErrPemDecode
		return
	}
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return
	}
	hash := sha256.New()
	_, err = hash.Write([]byte(content))
	if err != nil {
		return
	}

	// 3、验证签名
	pub := key.(*rsa.PublicKey)
	err = rsa.VerifyPKCS1v15(pub, crypto.SHA256, hash.Sum(nil), signature)
	return
}

// Verifier
type Verifier struct {
	V Verify
}

// Verify 验签
func (s *Verifier) Verify(content, sign, pubKey string) (err error) {
	return s.V(content, sign, pubKey)
}

const InvalidType = "invalid type=%v"

// InterfaceToSortedJSONStr 结构体、Map 转 待加签的排序的json字符串
// json按照字典序排序，值为空或者为0的忽略，不序列化为json的忽略(tag中`json:"-"`)，不参与加签的字段忽略(tag中`sign:"-"`)
func InterfaceToSortedJSONStr(i interface{}) (str string, err error) {
	// 1、数据提取，基础类型提取值，结构体、Map等转换为有序Map
	if i == nil {
		err = fmt.Errorf(InvalidType, i)
		return
	}
	v, err := interfaceValExtract(i)
	if err != nil {
		return
	}

	// 2、字符串类型直接返回
	if vStr, ok := v.(string); ok {
		str = vStr
		return
	}

	// 3、ToSignMap类型 检查len
	if vMap, ok := v.(ToSignMap); ok && len(vMap) == 0 {
		str = ""
		return
	}

	// 4、返回json字符串
	return MarshalToStr(v)
}

// MarshalToStr 转换为json not escape html 解决不转义< >的问题
func MarshalToStr(i interface{}) (str string, err error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err = encoder.Encode(i)
	if err != nil {
		return
	}
	str = strings.TrimSuffix(buffer.String(), "\n")
	return
}

// interfaceValExtract 提取i的值，i为0值或空值时返回""，结构体、Map 转 key排序的Map[string]interface{}
func interfaceValExtract(i interface{}) (v interface{}, err error) {
	// 1、构建默认返回值，反射获取i的类型与值
	v = ""
	typ := reflect.TypeOf(i)
	val := reflect.ValueOf(i)

	// 2、指针类型取出元素类型与值
	if typ.Kind() == reflect.Ptr {
		if val.IsNil() {
			return
		}
		typ = typ.Elem()
		val = val.Elem()
	}

	// 3、分类型处理
	k := typ.Kind()
	switch k {
	case reflect.Bool:
		v = val.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// 忽略0值
		if val.Int() == 0 {
			return
		}
		v = val.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		// 忽略0值
		if val.Uint() == 0 {
			return
		}
		v = val.Uint()
	case reflect.Float32, reflect.Float64:
		if val.IsZero() {
			return
		}
		v = val.Float()
	case reflect.String:
		v = val.String()
	case reflect.Slice, reflect.Array:
		if val.Len() == 0 {
			return
		}
		v, err = sliceValExtract(val)
	case reflect.Struct:
		if val.IsZero() {
			return
		}
		v, err = structValToSortedMap(typ, val)
	case reflect.Map:
		if val.Len() == 0 {
			return
		}
		v, err = mapValToSortedMap(val)
	// 其他类型不支持签名
	default:
		err = fmt.Errorf(InvalidType, k)
	}
	return
}

// structValToSortedMap 结构体转排序的json string，忽略空值和0值
func structValToSortedMap(typs reflect.Type, vals reflect.Value) (sc ToSignMap, err error) {
	// 1、构建map
	sc = make(ToSignMap)

	// 2、反射遍历属性
	num := vals.NumField()
	for i := 0; i < num; i++ {
		val := vals.Field(i)
		typ := typs.Field(i)
		// 判断是否为需要忽略的加签字段
		if isSkippedSignField(typ.Tag) {
			continue
		}
		// 判断属性是否可导出(私有属性不能导出)
		if !val.CanInterface() {
			continue
		}
		// 转换成排序类型
		var v interface{}
		v, err = interfaceValExtract(val.Interface())
		if err != nil {
			return
		}
		// 名称以结构体中的json标签名称为准
		name := typ.Name
		if jsonName := getJSONNameInTag(typ.Tag); jsonName != "" {
			name = jsonName
		}
		sc[name] = v
	}

	// 3、元素排序、去掉空值
	sc = sc.ToSortedNoZeroValue()
	return
}

func isSkippedSignField(tag reflect.StructTag) bool {
	// 1、忽略不序列化的字段
	v, ok := tag.Lookup("json")
	if ok && v == "-" {
		return true
	}

	// 2、忽略不加签的字段
	v, ok = tag.Lookup("sign")
	return ok && v == "-"
}

func getJSONNameInTag(tag reflect.StructTag) string {
	v, ok := tag.Lookup("json")
	if ok {
		return strings.Split(v, ",")[0]
	}
	return ""
}

// mapValToSortedMap map转排序的json string，忽略0值和空值
func mapValToSortedMap(vals reflect.Value) (sc ToSignMap, err error) {
	// 1、构建map
	sc = make(ToSignMap)
	// 2、反射遍历属性
	iter := vals.MapRange()
	for iter.Next() {
		// 处理key
		key, er := interfaceValExtract(iter.Key().Interface())
		if er != nil {
			err = er
			return
		}
		k := fmt.Sprintf("%v", key)

		// 处理value
		var val interface{}
		val, err = interfaceValExtract(iter.Value().Interface())
		if err != nil {
			return
		}

		// 赋值
		sc[k] = val
	}

	// 3、元素排序、去掉空值
	sc = sc.ToSortedNoZeroValue()
	return
}

// sliceValExtract 切片转忽略空值 或 配置了忽略签名 的切片
func sliceValExtract(vals reflect.Value) (s []interface{}, err error) {
	// 1、反射遍历属性
	num := vals.Len()
	for i := 0; i < num; i++ {
		// 类型判断
		val := vals.Index(i)
		k := val.Kind()
		if isNotValidType(k) {
			err = fmt.Errorf(InvalidType, k)
			return
		}

		// 判断属性是否可导出(私有属性不能导出)
		if !val.CanInterface() {
			continue
		}
		// 取出值
		v := val.Interface()

		// 结构体/Map/切片类型进行值的提取
		if k == reflect.Struct || k == reflect.Map || k == reflect.Slice || k == reflect.Array {
			// 提取切片的元素
			v, err = interfaceValExtract(val.Interface())
			if err != nil {
				return
			}
		}
		s = append(s, v)
	}

	// 2、返回处理后的切片
	return
}

func isNotValidType(k reflect.Kind) bool {
	return k == reflect.Invalid || k == reflect.Complex64 || k == reflect.Complex128 ||
		k == reflect.Chan || k == reflect.Func || k == reflect.UnsafePointer
}

// 待签名Map，提供转有序Map方法
type ToSignMap map[string]interface{}

// ToSortedNoZeroValue 转换为去除空值的按照字典序排序key的Map
func (sc ToSignMap) ToSortedNoZeroValue() ToSignMap {
	if len(sc) == 0 {
		return sc
	}

	// 1、取出sc的值不为空的key
	var keys []string
	for k, v := range sc {
		// 忽略空值
		if k == "" || v == "" {
			continue
		}
		keys = append(keys, k)
	}

	// 2、排序
	sort.Strings(keys)

	// 3、重组为排序的map
	sorted := make(ToSignMap)
	for _, v := range keys {
		sorted[v] = sc[v]
	}
	return sorted
}
