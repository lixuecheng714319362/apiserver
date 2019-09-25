package tool

import (
	"bytes"
	"crypto/cipher"
	sha2562 "crypto/sha256"
	"github.com/astaxie/beego"
	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/sm3"
	"github.com/tjfoc/gmsm/sm4"
	"net/http"
)

const (
	Ok         = 200 // 请求成功
	BadRequest = 400 // 请求失败
	Forbidden  = 403 // 参数错误
)

type Result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func ChangeArgsWithSm3(Args []string) (arg [][]byte) {
	for i, v := range Args {
		var hash []byte
		if i == 1 {
			hash = sm3EnCrypt(v)
		} else {
			hash = []byte(v)
		}
		arg = append(arg, hash)
	}
	return arg
}
func ChangeArgsWithSm4(Args []string) (arg [][]byte) {
	for i, v := range Args {
		var d []byte
		var err error
		if i == 1 {
			d, err = sm4Encrypt([]byte("1234567890abcdef"), []byte("1234567890abcdef"), []byte(v))
			if err != nil {
				panic(err)
			}

		} else {
			d = []byte(v)
		}
		arg = append(arg, d)
	}
	return arg
}
func ChangeArgsWithSm2(Args []string) (arg [][]byte) {
	for i, v := range Args {
		var d []byte
		var err error
		if i == 1 {
			d, err = sm2Encrypt([]byte(v))
			if err != nil {
				panic(err)
			}

		} else {
			d = []byte(v)
		}
		arg = append(arg, d)
	}
	return arg
}

func ChangeArgsWithSHA256(Args []string) (arg [][]byte) {
	for i, v := range Args {
		var d []byte
		if i == 1 {
			h := sha2562.Sum256([]byte(v))
			d = h[:]
		} else {
			d = []byte(v)
		}
		arg = append(arg, d)
	}
	return arg
}

func ChangeArgs(Args []string) (arg [][]byte) {
	for _, v := range Args {
		arg = append(arg, []byte(v))
	}
	return arg
}
var LogLevel string
func HanddlerError(c beego.Controller) {
	if LogLevel=="debug"{
		return
	}
	e := recover()
	if e != nil {
		beego.Error("panic err", e)
		BackResError(c, http.StatusBadRequest, "panic err")
	}
}

// 返回错误信息
func BackResError(this beego.Controller, code int, error string) {
	result := Result{Code: code, Msg: error}
	this.Data["json"] = result
	this.ServeJSON()
}

type ResData struct {
	Result
	Data interface{} `json:"data"`
}

// 返回查询信息
func BackResData(this beego.Controller, data interface{}) {
	result := ResData{
		Result: Result{
			Code: Ok,
			Msg:  "success",
		},
		Data: data,
	}
	this.Data["json"] = result
	this.ServeJSON()
}

// 返回请求状态
func BackResSuccess(this beego.Controller) {
	result := Result{
		Code: Ok,
		Msg:  "success",
	}
	this.Data["json"] = result
	this.ServeJSON()
}

// 返回请求状态
func BackResTimeOut(this beego.Controller) {
	result := Result{
		Code: Ok,
		Msg:  "api server time out",
	}
	this.Data["json"] = result
	this.ServeJSON()
}

type ResHash struct {
	Result
	Hash string `json:"hash"`
}

// 返回Hash信息
func BackResHash(this beego.Controller, hash string) {
	result := new(ResHash)
	result.Code = Ok
	result.Hash = hash
	this.Data["json"] = result
	this.ServeJSON()
}

// 组装args参数格式
func GetArgs(keys []string) [][]byte {
	args := make([][]byte, len(keys))
	for k, v := range keys {
		args[k] = []byte(v)
	}
	return args
}

func sm2Encrypt(data []byte) ([]byte, error) {
	pub, err := sm2.ReadPublicKeyFromPem("./controllers/tool/publicKey.pem", []byte(""))
	if err != nil {
		return nil, err
	}
	return sm2.Encrypt(pub, data)
}

func sm3EnCrypt(text string) []byte {
	h := sm3.New()
	return h.Sum([]byte(text))
}

func sm4Encrypt(key, iv, plainText []byte) ([]byte, error) {
	block, err := sm4.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData := pkcs5Padding(plainText, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, iv)
	cryted := make([]byte, len(origData))
	blockMode.CryptBlocks(cryted, origData)
	return cryted, nil
}

func sm4Decrypt(key, iv, cipherText []byte) ([]byte, error) {
	block, err := sm4.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(cipherText))
	blockMode.CryptBlocks(origData, cipherText)
	origData = pkcs5UnPadding(origData)
	return origData, nil
}

// pkcs5填充
func pkcs5Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func pkcs5UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}
